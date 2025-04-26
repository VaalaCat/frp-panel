package client

import (
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils/logger"
)

func StartFRPCHandler(ctx *app.Context, req *pb.StartFRPCRequest) (*pb.StartFRPCResponse, error) {
	logger.Logger(ctx).Infof("client get a start client request, origin is: [%+v]", req)

	if err := PullConfig(ctx.GetApp(), req.GetClientId(), ctx.GetApp().GetConfig().Client.Secret); err != nil {
		logger.Logger(ctx).WithError(err).Error("cannot pull client config")
		return nil, err
	}

	return &pb.StartFRPCResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
