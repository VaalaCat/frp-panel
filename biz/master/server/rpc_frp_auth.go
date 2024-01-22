package server

import (
	"context"
	"fmt"

	"github.com/VaalaCat/frp-panel/cache"
	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/sirupsen/logrus"
)

func FRPAuth(ctx context.Context, req *pb.FRPAuthRequest) (*pb.FRPAuthResponse, error) {
	logrus.Infof("frpc auth, req: [%+v]", req)
	var (
		err error
	)

	userToken, err := cache.Get().Get([]byte(req.User))
	if err != nil {
		u, err := dao.GetUserByUserName(req.User)
		if err != nil || u == nil {
			logrus.WithError(err).Errorf("invalid user: %s", req.User)
			return &pb.FRPAuthResponse{
				Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: err.Error()},
				Ok:     false,
			}, fmt.Errorf("invalid user: %s", req.User)
		}
		cache.Get().Set([]byte(u.GetUserName()), []byte(u.GetToken()), 0)
		userToken = []byte(u.GetToken())
	}

	if string(userToken) != req.GetToken() {
		logrus.Error("invalid token")
		return &pb.FRPAuthResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid token"},
			Ok:     false,
		}, fmt.Errorf("invalid token")
	}

	return &pb.FRPAuthResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
		Ok:     true,
	}, nil
}
