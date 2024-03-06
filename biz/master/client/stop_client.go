package client

import (
	"context"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/rpc"
	"github.com/sirupsen/logrus"
)

func StopFRPCHandler(ctx context.Context, req *pb.StopFRPCRequest) (*pb.StopFRPCResponse, error) {
	logrus.Infof("master get a stop client request, origin is: [%+v]", req)

	userInfo := common.GetUserInfo(ctx)
	clientID := req.GetClientId()

	if !userInfo.Valid() {
		return &pb.StopFRPCResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid user"},
		}, nil
	}

	if len(clientID) == 0 {
		return &pb.StopFRPCResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid client id"},
		}, nil
	}

	client, err := dao.GetClientByClientID(userInfo, clientID)
	if err != nil {
		return nil, err
	}

	client.Stopped = true

	if err = dao.UpdateClient(userInfo, client); err != nil {
		return nil, err
	}

	go func() {
		resp, err := rpc.CallClient(context.Background(), req.GetClientId(), pb.Event_EVENT_STOP_FRPC, req)
		if err != nil {
			logrus.WithError(err).Errorf("stop client event send to client error, client id: [%s]", req.GetClientId())
		}

		if resp == nil {
			logrus.Errorf("cannot get response, client id: [%s]", req.GetClientId())
		}
	}()

	return &pb.StopFRPCResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
