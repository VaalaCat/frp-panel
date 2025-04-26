package user

import (
	"context"

	"github.com/VaalaCat/frp-panel/biz/master/client"
	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/VaalaCat/frp-panel/utils/logger"
)

func UpdateUserInfoHander(c *app.Context, req *pb.UpdateUserInfoRequest) (*pb.UpdateUserInfoResponse, error) {
	var (
		userInfo = common.GetUserInfo(c)
	)

	if !userInfo.Valid() {
		return &pb.UpdateUserInfoResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid user"},
		}, nil
	}
	newUserEntity := userInfo.(*models.UserEntity)
	newUserInfo := req.GetUserInfo()

	if newUserInfo.GetEmail() != "" {
		newUserEntity.Email = newUserInfo.GetEmail()
	}

	if newUserInfo.GetRawPassword() != "" {
		hashedPassword, err := utils.HashPassword(newUserInfo.GetRawPassword())
		if err != nil {
			logger.Logger(context.Background()).WithError(err).Errorf("cannot hash password")
			return nil, err
		}
		newUserEntity.Password = hashedPassword
	}

	if newUserInfo.GetUserName() != "" {
		newUserEntity.UserName = newUserInfo.GetUserName()
	}

	if newUserInfo.GetToken() != "" {
		newUserEntity.Token = newUserInfo.GetToken()
	}

	if err := dao.NewQuery(c).UpdateUser(userInfo, newUserEntity); err != nil {
		return &pb.UpdateUserInfoResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: err.Error()},
		}, err
	}

	go func() {
		newUser, err := dao.NewQuery(app.NewContext(context.Background(), c.GetApp())).GetUserByUserID(userInfo.GetUserID())
		if err != nil {
			logger.Logger(context.Background()).WithError(err).Errorf("cannot get user")
			return
		}

		if err := client.SyncTunnel(c, newUser); err != nil {
			logger.Logger(context.Background()).WithError(err).Errorf("cannot sync tunnel, user need to retry update")
		}
	}()

	return &pb.UpdateUserInfoResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
