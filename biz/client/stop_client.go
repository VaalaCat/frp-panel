package client

import (
	"fmt"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils/logger"
)

func StopFRPCHandler(ctx *app.Context, req *pb.StopFRPCRequest) (*pb.StopFRPCResponse, error) {
	logger.Logger(ctx).Infof("client get a stop client request, origin is: [%+v]", req)

	ctx.GetApp().GetClientController().StopAll()
	ctx.GetApp().GetClientController().DeleteAll()

	if ctx.GetApp().GetConfig().Client.Features.EnableFunctions {
		ctx.GetApp().GetWorkersManager().StopAllWorkers(ctx)
	}

	errs := ctx.GetApp().GetWireGuardManager().StopAllServices()
	if len(errs) > 0 {
		logger.Logger(ctx).
			WithError(fmt.Errorf("wireguard manager stop all wireguard error, errs: %v", errs)).
			Errorf("wireguard manager stop all wireguard error")
	}

	return &pb.StopFRPCResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
