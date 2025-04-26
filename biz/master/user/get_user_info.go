package user

import (
	"fmt"

	"github.com/VaalaCat/frp-panel/app"
	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/samber/lo"
)

func GetUserInfoHandler(c *app.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	userInfo := common.GetUserInfo(c)
	return &pb.GetUserInfoResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
		UserInfo: &pb.User{
			UserID:   lo.ToPtr(int64(userInfo.GetUserID())),
			TenantID: lo.ToPtr(int64(userInfo.GetTenantID())),
			UserName: lo.ToPtr(userInfo.GetUserName()),
			Email:    lo.ToPtr(userInfo.GetEmail()),
			Status:   lo.ToPtr(fmt.Sprint(userInfo.GetStatus())),
			Role:     lo.ToPtr(userInfo.GetRole()),
			Token:    lo.ToPtr(userInfo.GetToken()),
		},
	}, nil
}
