package client

import (
	"context"

	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/pb"
)

func StartFRPCHandler(ctx context.Context, req *pb.StartFRPCRequest) (*pb.StartFRPCResponse, error) {
	logger.Logger(ctx).Infof("client get a start client request, origin is: [%+v]", req)

	if err := PullConfig(req.GetClientId(), conf.Get().Client.Secret); err != nil {
		logger.Logger(ctx).WithError(err).Error("cannot pull client config")
		return nil, err
	}

	return &pb.StartFRPCResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
