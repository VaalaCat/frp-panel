package client

import (
	"context"
	"fmt"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/rpc"
	"github.com/sirupsen/logrus"
)

func RemoveFrpcHandler(c context.Context, req *pb.RemoveFRPCRequest) (*pb.RemoveFRPCResponse, error) {
	logrus.Infof("remove frpc, req: [%+v]", req)

	var (
		clientID = req.GetClientId()
		userInfo = common.GetUserInfo(c)
	)

	if len(clientID) == 0 {
		logrus.Errorf("invalid client id")
		return nil, fmt.Errorf("invalid client id")
	}

	_, err := dao.GetClientByClientID(userInfo, clientID)
	if err != nil {
		logrus.WithError(err).Errorf("cannot get client, id: [%s]", clientID)
		return nil, err
	}

	err = dao.DeleteClient(userInfo, clientID)
	if err != nil {
		logrus.WithError(err).Errorf("cannot delete client, id: [%s]", clientID)
		return nil, err
	}

	go func() {
		resp, err := rpc.CallClient(c, req.GetClientId(), pb.Event_EVENT_REMOVE_FRPC, req)
		if err != nil {
			logrus.WithError(err).Errorf("remove event send to client error, client id: [%s]", req.GetClientId())
		}

		if resp == nil {
			logrus.Errorf("cannot get response, client id: [%s]", req.GetClientId())
		}
	}()

	logrus.Infof("remove frpc success, client id: [%s]", req.GetClientId())
	return &pb.RemoveFRPCResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
