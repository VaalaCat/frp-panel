package auth

import (
	"fmt"

	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/middleware"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
)

func LoginHandler(ctx *app.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	username := req.GetUsername()
	password := req.GetPassword()
	ok, user, err := dao.NewQuery(ctx).CheckUserPassword(username, password)
	if err != nil {
		return nil, err
	}

	if !ok {
		return &pb.LoginResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid username or password"},
		}, nil
	}

	tokenStr := conf.GetCommonJWT(ctx.GetApp().GetConfig(), fmt.Sprint(user.GetUserID()))

	ginCtx := ctx.GetGinCtx()
	middleware.PushTokenStr(ginCtx, ctx.GetApp(), tokenStr)

	return &pb.LoginResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
		Token:  &tokenStr,
	}, nil
}
