package client

import (
	"context"

	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/tunnel"
)

func StopFRPCHandler(ctx context.Context, req *pb.StopFRPCRequest) (*pb.StopFRPCResponse, error) {
	logger.Logger(ctx).Infof("client get a stop client request, origin is: [%+v]", req)

	tunnel.GetClientController().StopAll()
	tunnel.GetClientController().DeleteAll()

	return &pb.StopFRPCResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
