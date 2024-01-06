package auth

import (
	"context"
	"fmt"

	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/middleware"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/gin-gonic/gin"
)

func LoginHandler(c context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	username := req.GetUsername()
	password := req.GetPassword()
	ok, user, err := dao.CheckUserPassword(username, password)
	if err != nil {
		return nil, err
	}

	if !ok {
		return &pb.LoginResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid username or password"},
		}, nil
	}

	tokenStr := conf.GetCommonJWT(fmt.Sprint(user.GetUserID()))

	ginCtx := c.(*gin.Context)
	middleware.PushTokenStr(ginCtx, tokenStr)

	return &pb.LoginResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
		Token:  &tokenStr,
	}, nil
}
