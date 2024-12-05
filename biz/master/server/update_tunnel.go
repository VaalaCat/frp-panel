package server

import (
	"context"
	"fmt"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/rpc"
	"github.com/VaalaCat/frp-panel/utils"
	v1 "github.com/fatedier/frp/pkg/config/v1"
)

func UpdateFrpsHander(c context.Context, req *pb.UpdateFRPSRequest) (*pb.UpdateFRPSResponse, error) {
	logger.Logger(c).Infof("update frps, req: [%+v]", req)
	var (
		serverID  = req.GetServerId()
		configStr = req.GetConfig()
		userInfo  = common.GetUserInfo(c)
	)

	if len(configStr) == 0 || len(serverID) == 0 {
		return nil, fmt.Errorf("request invalid")
	}

	srv, err := dao.GetServerByServerID(userInfo, serverID)
	if srv == nil || err != nil {
		logger.Logger(context.Background()).WithError(err).Errorf("cannot get server, id: [%s]", serverID)
		return nil, err
	}

	srvCfg, err := utils.LoadServerConfig(configStr, true)
	if srvCfg == nil || err != nil {
		logger.Logger(context.Background()).WithError(err).Errorf("cannot load server config")
		return nil, err
	}

	srvCfg.HTTPPlugins = []v1.HTTPPluginOptions{conf.FRPsAuthOption(common.DefaultServerID == serverID)}

	if err := srv.SetConfigContent(srvCfg); err != nil {
		logger.Logger(context.Background()).WithError(err).Errorf("cannot set server config")
		return nil, err
	}

	srv.Comment = req.GetComment()
	if len(req.GetServerIp()) > 0 {
		srv.ServerIP = req.GetServerIp()
	}

	if err := dao.UpdateServer(userInfo, srv); err != nil {
		logger.Logger(context.Background()).WithError(err).Errorf("cannot update server, id: [%s]", serverID)
		return nil, err
	}

	go func() {
		resp, err := rpc.CallClient(context.Background(), req.GetServerId(), pb.Event_EVENT_UPDATE_FRPS, req)
		if err != nil {
			logger.Logger(context.Background()).WithError(err).Errorf("update event send to server error, server id: [%s]", req.GetServerId())
		}
		if resp == nil {
			logger.Logger(c).Errorf("cannot get response, server id: [%s]", req.GetServerId())
		}
	}()

	logger.Logger(c).Infof("update frps success, id: [%s]", serverID)
	return &pb.UpdateFRPSResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
