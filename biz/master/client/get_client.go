package client

import (
	"context"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/samber/lo"
)

func GetClientHandler(ctx context.Context, req *pb.GetClientRequest) (*pb.GetClientResponse, error) {
	logger.Logger(ctx).Infof("get client, req: [%+v]", req)

	var (
		userInfo = common.GetUserInfo(ctx)
		clientID = req.GetClientId()
	)

	if !userInfo.Valid() {
		return &pb.GetClientResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid user"},
		}, nil
	}

	if len(clientID) == 0 {
		return &pb.GetClientResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid client id"},
		}, nil
	}

	client, err := dao.GetClientByClientID(userInfo, clientID)
	if err != nil {
		return nil, err
	}

	return &pb.GetClientResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
		Client: &pb.Client{
			Id:       lo.ToPtr(client.ClientID),
			Secret:   lo.ToPtr(client.ConnectSecret),
			Config:   lo.ToPtr(string(client.ConfigContent)),
			ServerId: lo.ToPtr(client.ServerID),
			Stopped:  lo.ToPtr(client.Stopped),
			Comment:  lo.ToPtr(client.Comment),
		},
	}, nil
}
