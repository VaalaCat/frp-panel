package client

import (
	"context"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/VaalaCat/frp-panel/services/rpc"
	"github.com/VaalaCat/frp-panel/utils/logger"
)

func StopFRPCHandler(ctx *app.Context, req *pb.StopFRPCRequest) (*pb.StopFRPCResponse, error) {
	logger.Logger(ctx).Infof("master get a stop client request, origin is: [%+v]", req)

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

	cli, err := dao.NewQuery(ctx).GetClientByClientID(userInfo, clientID)
	if err != nil {
		return nil, err
	}

	client := cli.ClientEntity

	client.Stopped = true

	if err = dao.NewQuery(ctx).UpdateClient(userInfo, client); err != nil {
		return nil, err
	}

	go func() {
		resp, err := rpc.CallClient(app.NewContext(context.Background(), ctx.GetApp()), req.GetClientId(), pb.Event_EVENT_STOP_FRPC, req)
		if err != nil {
			logger.Logger(context.Background()).WithError(err).Errorf("stop client event send to client error, client id: [%s]", req.GetClientId())
		}

		if resp == nil {
			logger.Logger(ctx).Errorf("cannot get response, client id: [%s]", req.GetClientId())
		}
	}()

	return &pb.StopFRPCResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
