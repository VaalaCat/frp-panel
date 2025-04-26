package client

import (
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils/logger"
)

func StopFRPCHandler(ctx *app.Context, req *pb.StopFRPCRequest) (*pb.StopFRPCResponse, error) {
	logger.Logger(ctx).Infof("client get a stop client request, origin is: [%+v]", req)

	ctx.GetApp().GetClientController().StopAll()
	ctx.GetApp().GetClientController().DeleteAll()

	return &pb.StopFRPCResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
