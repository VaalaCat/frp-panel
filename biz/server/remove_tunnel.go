package server

import (
	"context"
	"fmt"

	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/tunnel"
)

func RemoveFrpsHandler(ctx context.Context, req *pb.RemoveFRPSRequest) (*pb.RemoveFRPSResponse, error) {
	logger.Logger(ctx).Infof("remove frps, req: [%+v]", req)

	if req.GetServerId() == "" {
		return &pb.RemoveFRPSResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "server id cannot be empty"},
		}, fmt.Errorf("server id cannot be empty")
	}

	srv := tunnel.GetServerController().Get(req.GetServerId())
	if srv == nil {
		logger.Logger(ctx).Infof("server not found, no need to remove")
		return &pb.RemoveFRPSResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "server not found"},
		}, nil
	}

	srv.Stop()
	tunnel.GetServerController().Delete(req.GetServerId())

	return &pb.RemoveFRPSResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
