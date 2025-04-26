package auth

import (
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
)

func GetClientCert(ctx *app.Context, req *pb.GetClientCertRequest) (*pb.GetClientCertResponse, error) {
	var err error
	if req.ClientType == pb.ClientType_CLIENT_TYPE_FRPC {
		_, err = dao.NewQuery(ctx).ValidateClientSecret(req.GetClientId(), req.GetClientSecret())
	}
	if req.ClientType == pb.ClientType_CLIENT_TYPE_FRPS {
		_, err = dao.NewQuery(ctx).ValidateServerSecret(req.GetClientId(), req.GetClientSecret())
	}
	if err != nil {
		return &pb.GetClientCertResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: err.Error()},
		}, err
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
