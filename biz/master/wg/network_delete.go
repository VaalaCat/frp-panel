package wg

import (
	"errors"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
)

func DeleteNetwork(ctx *app.Context, req *pb.DeleteNetworkRequest) (*pb.DeleteNetworkResponse, error) {
	userInfo := common.GetUserInfo(ctx)
	if !userInfo.Valid() {
		return nil, errors.New("invalid user")
	}
	id := uint(req.GetId())
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	if err := dao.NewQuery(ctx).DeleteNetwork(userInfo, id); err != nil {
		return nil, err
	}
	return &pb.DeleteNetworkResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"},
	}, nil
}
