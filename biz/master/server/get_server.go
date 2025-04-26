package server

import (
	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/samber/lo"
)

func GetServerHandler(c *app.Context, req *pb.GetServerRequest) (*pb.GetServerResponse, error) {
	var (
		userServerID = req.GetServerId()
		userInfo     = common.GetUserInfo(c)
	)

	if !userInfo.Valid() {
		return &pb.GetServerResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid user"},
		}, nil
	}

	if len(userServerID) == 0 {
		return &pb.GetServerResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid client id"},
		}, nil
	}

	serverEntity, err := dao.NewQuery(c).GetServerByServerID(userInfo, userServerID)
	if err != nil {
		return nil, err
	}

	return &pb.GetServerResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
		Server: &pb.Server{
			Id:      lo.ToPtr(serverEntity.ServerID),
			Config:  lo.ToPtr(string(serverEntity.ConfigContent)),
			Secret:  lo.ToPtr(serverEntity.ConnectSecret),
			Comment: lo.ToPtr(serverEntity.Comment),
			Ip:      lo.ToPtr(serverEntity.ServerIP),
		},
	}, nil
}
