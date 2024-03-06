package client

import (
	"context"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/tunnel"
	"github.com/sirupsen/logrus"
)

func StartFRPCHandler(ctx context.Context, req *pb.StartFRPCRequest) (*pb.StartFRPCResponse, error) {
	logrus.Infof("client get a start client request, origin is: [%+v]", req)

	tunnel.GetClientController().Run(req.GetClientId())

	return &pb.StartFRPCResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
