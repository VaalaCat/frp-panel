package auth

import (
	"errors"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
)

func GetClientCert(ctx *app.Context, req *pb.GetClientCertRequest) (*pb.GetClientCertResponse, error) {
	var err error
	switch req.GetClientType() {
	case pb.ClientType_CLIENT_TYPE_FRPC:
		_, err = dao.NewQuery(ctx).ValidateClientSecret(req.GetClientId(), req.GetClientSecret())
	case pb.ClientType_CLIENT_TYPE_FRPS:
		_, err = dao.NewQuery(ctx).ValidateServerSecret(req.GetClientId(), req.GetClientSecret())
	default:
		err = errors.New("invalid client credentials")
	}
	if err != nil {
		return &pb.GetClientCertResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid client credentials"},
		}, errors.New("invalid client credentials")
	}

	_, cert, err := dao.NewQuery(ctx).GetDefaultKeyPair()
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
