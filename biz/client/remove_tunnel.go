package client

import (
	"os"
	"time"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils/logger"
)

func RemoveFrpcHandler(ctx *app.Context, req *pb.RemoveFRPCRequest) (*pb.RemoveFRPCResponse, error) {
	logger.Logger(ctx).Infof("remove frpc, req: [%+v], will exit in 10s", req)

	go func() {
		time.Sleep(10 * time.Second)
		os.Exit(0)
	}()

	return &pb.RemoveFRPCResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
