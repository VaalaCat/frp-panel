package proxy

import (
	"github.com/VaalaCat/frp-panel/biz/master/client"
	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/VaalaCat/frp-panel/utils/logger"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/samber/lo"
)

func StopProxy(ctx *app.Context, req *pb.StopProxyRequest) (*pb.StopProxyResponse, error) {
	var (
		userInfo  = common.GetUserInfo(ctx)
		clientID  = req.GetClientId()
		serverID  = req.GetServerId()
		proxyName = req.GetName()
	)

	clientEntity, err := getClientWithMakeShadow(ctx, clientID, serverID)
	if err != nil {
		logger.Logger(ctx).WithError(err).Errorf("cannot get client, id: [%s]", clientID)
		return nil, err
	}

	_, err = dao.NewQuery(ctx).GetServerByServerID(userInfo, serverID)
	if err != nil {
		logger.Logger(ctx).WithError(err).Errorf("cannot get server, id: [%s]", serverID)
		return nil, err
	}

	proxyConfig, err := dao.NewQuery(ctx).GetProxyConfigByFilter(userInfo, &models.ProxyConfigEntity{
		ClientID: clientID,
		ServerID: serverID,
		Name:     proxyName,
	})
	if err != nil {
		logger.Logger(ctx).WithError(err).Errorf("cannot get proxy config, client: [%s], server: [%s], proxy name: [%s]", clientID, serverID, proxyName)
		return nil, err
	}

	// 1. 更新proxy状态
	proxyConfig.Stopped = true
	err = dao.NewQuery(ctx).UpdateProxyConfig(userInfo, proxyConfig)
	if err != nil {
		logger.Logger(ctx).WithError(err).Errorf("cannot update proxy config, client: [%s], server: [%s], proxy name: [%s]", clientID, serverID, proxyName)
		return nil, err
	}

	// 2. 从client移除proxy
	if oldCfg, err := clientEntity.GetConfigContent(); err != nil {
		logger.Logger(ctx).WithError(err).Errorf("cannot get client config, id: [%s]", clientID)
		return nil, err
	} else {
		oldCfg.Proxies = lo.Filter(oldCfg.Proxies, func(proxy v1.TypedProxyConfig, _ int) bool {
			return proxy.GetBaseConfig().Name != proxyName
		})

		if err := clientEntity.SetConfigContent(*oldCfg); err != nil {
			logger.Logger(ctx).WithError(err).Errorf("cannot set client config, id: [%s]", clientID)
			return nil, err
		}
	}

	// 3. 更新client的配置
	rawCfg, err := clientEntity.MarshalJSONConfig()
	if err != nil {
		logger.Logger(ctx).WithError(err).Errorf("cannot marshal client config, id: [%s]", clientID)
		return nil, err
	}

	_, err = client.UpdateFrpcHander(ctx, &pb.UpdateFRPCRequest{
		ClientId: &clientEntity.ClientID,
		ServerId: &serverID,
		Config:   rawCfg,
		Comment:  &clientEntity.Comment,
		FrpsUrl:  &clientEntity.FrpsUrl,
	})
	if err != nil {
		logger.Logger(ctx).WithError(err).Warnf("cannot update frpc, id: [%s]", clientID)
	}

	return &pb.StopProxyResponse{
		Status: &pb.Status{
			Code:    pb.RespCode_RESP_CODE_SUCCESS,
			Message: "stop proxy success",
		},
	}, nil
}
