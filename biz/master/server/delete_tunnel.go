package server

import (
	"context"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/rpc"
	"github.com/sirupsen/logrus"
)

func RemoveFrpsHandler(c context.Context, req *pb.RemoveFRPSRequest) (*pb.RemoveFRPSResponse, error) {
	logrus.Infof("remove frps, req: [%+v]", req)

	var (
		serverID = req.GetServerId()
		userInfo = common.GetUserInfo(c)
	)

	srv, err := dao.GetServerByServerID(userInfo, serverID)
	if srv == nil || err != nil {
		logrus.WithError(err).Errorf("cannot get server, id: [%s]", serverID)
		return nil, err
	}

	if err = dao.DeleteServer(userInfo, serverID); err != nil {
		logrus.WithError(err).Errorf("cannot delete server, id: [%s]", serverID)
		return nil, err
	}

	go func() {
		resp, err := rpc.CallClient(context.Background(), req.GetServerId(), pb.Event_EVENT_REMOVE_FRPS, req)
		if err != nil {
			logrus.WithError(err).Errorf("remove event send to server error, server id: [%s]", req.GetServerId())
		}

		if resp == nil {
			logrus.Errorf("cannot get response, server id: [%s]", req.GetServerId())
		}
	}()

	return &pb.RemoveFRPSResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
