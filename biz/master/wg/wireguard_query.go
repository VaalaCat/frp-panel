package wg

import (
	"errors"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/samber/lo"
)

func GetWireGuard(ctx *app.Context, req *pb.GetWireGuardRequest) (*pb.GetWireGuardResponse, error) {
	userInfo := common.GetUserInfo(ctx)
	if !userInfo.Valid() {
		return nil, errors.New("invalid user")
	}
	id := uint(req.GetId())
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	wg, err := dao.NewQuery(ctx).GetWireGuardByID(userInfo, id)
	if err != nil {
		return nil, err
	}
	resp := wg.ToPB()
	return &pb.GetWireGuardResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"}, WireguardConfig: resp}, nil
}

func ListWireGuards(ctx *app.Context, req *pb.ListWireGuardsRequest) (*pb.ListWireGuardsResponse, error) {
	userInfo := common.GetUserInfo(ctx)
	if !userInfo.Valid() {
		return nil, errors.New("invalid user")
	}
	page, pageSize := int(req.GetPage()), int(req.GetPageSize())
	keyword := req.GetKeyword()
	filter := &models.WireGuardEntity{}
	if cid := req.GetClientId(); len(cid) > 0 {
		filter.ClientID = cid
	}
	// proto 的 network_id 是 string，但模型为 uint；此处仅在非空时参与过滤，解析失败则忽略
	if nid := req.GetNetworkId(); nid > 0 {
		filter.NetworkID = uint(nid)
	}
	list, err := dao.NewQuery(ctx).ListWireGuardsWithFilters(userInfo, page, pageSize, filter, keyword)
	if err != nil {
		return nil, err
	}
	total, err := dao.NewQuery(ctx).CountWireGuardsWithFilters(userInfo, filter, keyword)
	if err != nil {
		return nil, err
	}
	resp := lo.Map(list, func(item *models.WireGuard, _ int) *pb.WireGuardConfig {
		return item.ToPB()
	})
	return &pb.ListWireGuardsResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"}, WireguardConfigs: resp, Total: lo.ToPtr(int32(total))}, nil
}
