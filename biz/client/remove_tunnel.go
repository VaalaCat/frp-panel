package client

import (
	"context"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/tunnel"
	"github.com/sirupsen/logrus"
)

func RemoveFrpcHandler(ctx context.Context, req *pb.RemoveFRPCRequest) (*pb.RemoveFRPCResponse, error) {
	logrus.Infof("remove frpc, req: [%+v]", req)
	cli := tunnel.GetClientController().Get(req.GetClientId())
	if cli == nil {
		logrus.Infof("client not found, no need to remove")
		return &pb.RemoveFRPCResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "client not found"},
		}, nil
	}
	cli.Stop()
	tunnel.GetClientController().Delete(req.GetClientId())

	return &pb.RemoveFRPCResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
