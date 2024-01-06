package server

import (
	"context"

	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/sirupsen/logrus"
)

func FRPAuth(ctx context.Context, req *pb.FRPAuthRequest) (*pb.FRPAuthResponse, error) {
	logrus.Infof("frpc auth, req: [%+v]", req)
	var (
		err error
		cli *models.ServerEntity
	)

	if cli, err = ValidateServerRequest(req.GetBase()); err != nil {
		return &pb.FRPAuthResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: err.Error()},
			Ok:     false,
		}, err
	}

	logrus.Infof("frpc auth success, server: [%+v]", cli.ServerID)

	return &pb.FRPAuthResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
		Ok:     true,
	}, nil
}
