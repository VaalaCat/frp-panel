package auth

import (
	"context"

	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/pb"
)

func GetClientCert(ctx context.Context, req *pb.GetClientCertRequest) (*pb.GetClientCertResponse, error) {
	var err error
	if req.ClientType == pb.ClientType_CLIENT_TYPE_FRPC {
		_, err = dao.ValidateClientSecret(req.GetClientId(), req.GetClientSecret())
	}
	if req.ClientType == pb.ClientType_CLIENT_TYPE_FRPS {
		_, err = dao.ValidateServerSecret(req.GetClientId(), req.GetClientSecret())
	}
	if err != nil {
		return &pb.GetClientCertResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: err.Error()},
		}, err
	}

	_, cert, err := dao.GetDefaultKeyPair()
	if err != nil {
		return &pb.GetClientCertResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: err.Error()},
		}, err
	}
	return &pb.GetClientCertResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
		Cert:   cert,
	}, nil
}
