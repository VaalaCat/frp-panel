package server

import (
	"context"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/VaalaCat/frp-panel/services/rpc"
	"github.com/VaalaCat/frp-panel/utils/logger"
)

func RemoveFrpsHandler(c *app.Context, req *pb.RemoveFRPSRequest) (*pb.RemoveFRPSResponse, error) {
	logger.Logger(c).Infof("remove frps, req: [%+v]", req)

	var (
		serverID = req.GetServerId()
		userInfo = common.GetUserInfo(c)
	)

	srv, err := dao.NewQuery(c).GetServerByServerID(userInfo, serverID)
	if srv == nil || err != nil {
		logger.Logger(context.Background()).WithError(err).Errorf("cannot get server, id: [%s]", serverID)
		return nil, err
	}

	if err = dao.NewQuery(c).DeleteServer(userInfo, serverID); err != nil {
		logger.Logger(context.Background()).WithError(err).Errorf("cannot delete server, id: [%s]", serverID)
		return nil, err
	}

	go func() {
		resp, err := rpc.CallClient(app.NewContext(context.Background(), c.GetApp()), req.GetServerId(), pb.Event_EVENT_REMOVE_FRPS, req)
		if err != nil {
			logger.Logger(context.Background()).WithError(err).Errorf("remove event send to server error, server id: [%s]", req.GetServerId())
		}

		if resp == nil {
			logger.Logger(c).Errorf("cannot get response, server id: [%s]", req.GetServerId())
		}
	}()

	return &pb.RemoveFRPSResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
