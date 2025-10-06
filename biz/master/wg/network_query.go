package wg

import (
	"errors"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/samber/lo"
)

func GetNetwork(ctx *app.Context, req *pb.GetNetworkRequest) (*pb.GetNetworkResponse, error) {
	userInfo := common.GetUserInfo(ctx)
	if !userInfo.Valid() {
		return nil, errors.New("invalid user")
	}
	id := uint(req.GetId())
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	net, err := dao.NewQuery(ctx).GetNetworkByID(userInfo, id)
	if err != nil {
		return nil, err
	}
	logger.Logger(ctx).Infof("get network: %+v", net)
	return &pb.GetNetworkResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"},
		Network: net.ToPB(),
	}, nil
}

func ListNetworks(ctx *app.Context, req *pb.ListNetworksRequest) (*pb.ListNetworksResponse, error) {
	userInfo := common.GetUserInfo(ctx)
	if !userInfo.Valid() {
		return nil, errors.New("invalid user")
	}
	page, pageSize := int(req.GetPage()), int(req.GetPageSize())
	keyword := req.GetKeyword()
	var (
		list  []*models.Network
		err   error
		total int64
	)
	if len(keyword) > 0 {
		list, err = dao.NewQuery(ctx).ListNetworksWithKeyword(userInfo, page, pageSize, keyword)
		if err != nil {
			return nil, err
		}
		total, err = dao.NewQuery(ctx).CountNetworksWithKeyword(userInfo, keyword)
		if err != nil {
			return nil, err
		}
	} else {
		list, err = dao.NewQuery(ctx).ListNetworks(userInfo, page, pageSize)
		if err != nil {
			return nil, err
		}
		total, err = dao.NewQuery(ctx).CountNetworks(userInfo)
		if err != nil {
			return nil, err
		}
	}
	resp := lo.Map(list, func(item *models.Network, _ int) *pb.Network {
		return item.ToPB()
	})

	return &pb.ListNetworksResponse{
		Status:   &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"},
		Networks: resp,
		Total:    lo.ToPtr(int32(total)),
	}, nil
}
