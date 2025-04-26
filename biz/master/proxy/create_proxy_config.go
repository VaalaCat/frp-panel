package proxy

import (
	"errors"
	"fmt"

	"github.com/VaalaCat/frp-panel/biz/master/client"
	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/VaalaCat/frp-panel/utils/logger"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/samber/lo"

	"gorm.io/gorm"
)

func CreateProxyConfig(c *app.Context, req *pb.CreateProxyConfigRequest) (*pb.CreateProxyConfigResponse, error) {

	if len(req.GetClientId()) == 0 || len(req.GetServerId()) == 0 || len(req.GetConfig()) == 0 {
		return nil, fmt.Errorf("request invalid")
	}

	var (
		userInfo = common.GetUserInfo(c)
		clientID = req.GetClientId()
		serverID = req.GetServerId()
	)

	// 1. 检查是否有已连接该服务端的客户端
	// 2. 检查是否有Shadow客户端
	// 3. 如果没有，则新建Shadow客户端和子客户端
	clientEntity, err := dao.NewQuery(c).GetClientByFilter(userInfo, &models.ClientEntity{OriginClientID: clientID, ServerID: serverID}, lo.ToPtr(false))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		clientEntity, err = dao.NewQuery(c).GetClientByFilter(userInfo, &models.ClientEntity{ClientID: clientID}, nil)
		if err != nil {
			logger.Logger(c).WithError(err).Errorf("cannot get client, id: [%s]", clientID)
			return nil, err
		}
		if (!clientEntity.IsShadow || len(clientEntity.ConfigContent) != 0) && len(clientEntity.OriginClientID) == 0 {
			// 没shadow过，需要shadow
			_, err = client.MakeClientShadowed(c, serverID, clientEntity)
			if err != nil {
				logger.Logger(c).WithError(err).Errorf("cannot make client shadow, id: [%s]", clientID)
				return nil, err
			}
		}
		// shadow过，但没找到子客户端，需要新建
		clientEntity, err = client.ChildClientForServer(c, serverID, clientEntity)
		if err != nil {
			logger.Logger(c).WithError(err).Errorf("cannot create child client, id: [%s]", clientID)
			return nil, err
		}
	}
	// 有任何失败，返回
	if err != nil {
		logger.Logger(c).WithError(err).Errorf("cannot get client, id: [%s]", clientID)
		return nil, err
	}

	_, err = dao.NewQuery(c).GetServerByServerID(userInfo, serverID)
	if err != nil {
		logger.Logger(c).WithError(err).Errorf("cannot get server, id: [%s]", serverID)
		return nil, err
	}

	proxyCfg := &models.ProxyConfigEntity{}

	if err := proxyCfg.FillClientConfig(clientEntity); err != nil {
		logger.Logger(c).WithError(err).Errorf("cannot fill client config, id: [%s]", clientID)
		return nil, err
	}

	typedProxyCfgs, err := utils.LoadProxiesFromContent(req.GetConfig())
	if err != nil {
		logger.Logger(c).WithError(err).Errorf("cannot load proxies from content")
		return nil, err
	}
	if len(typedProxyCfgs) == 0 || len(typedProxyCfgs) > 1 {
		logger.Logger(c).Errorf("invalid config, cfg len: [%d]", len(typedProxyCfgs))
		return nil, fmt.Errorf("invalid config")
	}

	if err := proxyCfg.FillTypedProxyConfig(typedProxyCfgs[0]); err != nil {
		logger.Logger(c).WithError(err).Errorf("cannot fill typed proxy config")
		return nil, err
	}

	var existedProxyCfg *models.ProxyConfig
	existedProxyCfg, err = dao.NewQuery(c).GetProxyConfigByOriginClientIDAndName(userInfo, clientID, proxyCfg.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Logger(c).WithError(err).Errorf("cannot get proxy config, id: [%s]", clientID)
		return nil, err
	}

	if !req.GetOverwrite() && err == nil {
		logger.Logger(c).Errorf("proxy config already exist, cfg: [%+v]", proxyCfg)
		return nil, fmt.Errorf("proxy config already exist")
	}

	// update client config
	if oldCfg, err := clientEntity.GetConfigContent(); err != nil {
		logger.Logger(c).WithError(err).Errorf("cannot get client config, id: [%s]", clientID)
		return nil, err
	} else {
		oldCfg.Proxies = lo.Filter(oldCfg.Proxies, func(proxy v1.TypedProxyConfig, _ int) bool {
			return proxy.GetBaseConfig().Name != typedProxyCfgs[0].GetBaseConfig().Name
		})
		oldCfg.Proxies = append(oldCfg.Proxies, typedProxyCfgs...)

		if err := clientEntity.SetConfigContent(*oldCfg); err != nil {
			logger.Logger(c).WithError(err).Errorf("cannot set client config, id: [%s]", clientID)
			return nil, err
		}
	}

	rawCfg, err := clientEntity.MarshalJSONConfig()
	if err != nil {
		logger.Logger(c).WithError(err).Errorf("cannot marshal client config, id: [%s]", clientID)
		return nil, err
	}

	_, err = client.UpdateFrpcHander(c, &pb.UpdateFRPCRequest{
		ClientId: &clientEntity.ClientID,
		ServerId: &serverID,
		Config:   rawCfg,
		Comment:  &clientEntity.Comment,
	})
	if err != nil {
		logger.Logger(c).WithError(err).Warnf("cannot update frpc failed, id: [%s]", clientID)
	}

	if existedProxyCfg != nil && existedProxyCfg.ServerID != serverID {
		logger.Logger(c).Warnf("client and server not match, delete old proxy, client: [%s], server: [%s], proxy: [%s]", clientID, serverID, proxyCfg.Name)
		if _, err := DeleteProxyConfig(c, &pb.DeleteProxyConfigRequest{
			ClientId: lo.ToPtr(existedProxyCfg.ClientID),
			ServerId: lo.ToPtr(existedProxyCfg.ServerID),
			Name:     &proxyCfg.Name,
		}); err != nil {
			logger.Logger(c).WithError(err).Errorf("cannot delete old proxy, client: [%s], server: [%s], proxy: [%s]", clientID, clientEntity.ServerID, proxyCfg.Name)
			return nil, err
		}
	}

	return &pb.CreateProxyConfigResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
