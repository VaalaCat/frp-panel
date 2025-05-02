package proxy

import (
	"fmt"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/VaalaCat/frp-panel/services/rpc"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/samber/lo"
)

func GetProxyConfig(c *app.Context, req *pb.GetProxyConfigRequest) (*pb.GetProxyConfigResponse, error) {
	var (
		userInfo  = common.GetUserInfo(c)
		clientID  = req.GetClientId()
		serverID  = req.GetServerId()
		proxyName = req.GetName()
	)

	proxyConfig, err := dao.NewQuery(c).GetProxyConfigByFilter(userInfo, &models.ProxyConfigEntity{
		ClientID: clientID,
		ServerID: serverID,
		Name:     proxyName,
	})
	if err != nil {
		logger.Logger(c).WithError(err).Errorf("cannot get proxy config, client: [%s], server: [%s], proxy name: [%s]", clientID, serverID, proxyName)
		return nil, err
	}

	resp := &pb.GetProxyConfigResponse{WorkingStatus: &pb.ProxyWorkingStatus{
		Status: lo.ToPtr("stopped"),
	}}
	if !proxyConfig.Stopped {
		if err := rpc.CallClientWrapper(c, proxyConfig.OriginClientID, pb.Event_EVENT_GET_PROXY_INFO, &pb.GetProxyConfigRequest{
			ClientId: lo.ToPtr(proxyConfig.ClientID),
			ServerId: lo.ToPtr(proxyConfig.ServerID),
			Name:     lo.ToPtr(fmt.Sprintf("%s.%s", userInfo.GetUserName(), proxyName)),
		}, resp); err != nil {
			resp.WorkingStatus = &pb.ProxyWorkingStatus{
				Status: lo.ToPtr("error"),
			}
			logger.Logger(c).WithError(err).Errorf("cannot get proxy config, client: [%s], server: [%s], proxy name: [%s]", proxyConfig.OriginClientID, proxyConfig.ServerID, proxyConfig.Name)
		}

		if len(resp.GetWorkingStatus().GetStatus()) == 0 {
			resp.WorkingStatus = &pb.ProxyWorkingStatus{
				Status: lo.ToPtr("unknown"),
			}
		}
	}

	return &pb.GetProxyConfigResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"},
		ProxyConfig: &pb.ProxyConfig{
			Id:             lo.ToPtr(uint32(proxyConfig.ID)),
			Name:           lo.ToPtr(proxyConfig.Name),
			Type:           lo.ToPtr(proxyConfig.Type),
			ClientId:       lo.ToPtr(proxyConfig.ClientID),
			ServerId:       lo.ToPtr(proxyConfig.ServerID),
			Config:         lo.ToPtr(string(proxyConfig.Content)),
			OriginClientId: lo.ToPtr(proxyConfig.OriginClientID),
			Stopped:        lo.ToPtr(proxyConfig.Stopped),
		},
		WorkingStatus: resp.GetWorkingStatus(),
	}, nil
}
