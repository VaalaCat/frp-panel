package wg

import (
	"errors"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/samber/lo"
)

func CreateWireGuardLink(ctx *app.Context, req *pb.CreateWireGuardLinkRequest) (*pb.CreateWireGuardLinkResponse, error) {
	userInfo := common.GetUserInfo(ctx)
	if !userInfo.Valid() {
		return nil, errors.New("invalid user")
	}
	l := req.GetWireguardLink()
	if l == nil || l.GetFromWireguardId() == 0 || l.GetToWireguardId() == 0 || l.GetFromWireguardId() == l.GetToWireguardId() {
		return nil, errors.New("invalid link params")
	}
	// 校验两端属于同一 network
	q := dao.NewQuery(ctx)
	mut := dao.NewMutation(ctx)

	from, err := q.GetWireGuardByID(userInfo, uint(l.GetFromWireguardId()))
	if err != nil {
		return nil, err
	}
	to, err := q.GetWireGuardByID(userInfo, uint(l.GetToWireguardId()))
	if err != nil {
		return nil, err
	}
	if from.NetworkID == 0 || from.NetworkID != to.NetworkID {
		return nil, errors.New("wireguard not in same network")
	}
	m := &models.WireGuardLink{}
	m.FromPB(l)
	m.NetworkID = from.NetworkID

	reverse := &models.WireGuardLink{}
	reverse.FromPB((&defs.WireGuardLink{WireGuardLink: l}).GetReverse().WireGuardLink)
	reverse.NetworkID = from.NetworkID
	reverse.ToEndpointID = 0

	if err := mut.CreateWireGuardLinks(userInfo, m, reverse); err != nil {
		return nil, err
	}
	return &pb.CreateWireGuardLinkResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"}, WireguardLink: m.ToPB()}, nil
}

func UpdateWireGuardLink(ctx *app.Context, req *pb.UpdateWireGuardLinkRequest) (*pb.UpdateWireGuardLinkResponse, error) {
	userInfo := common.GetUserInfo(ctx)
	if !userInfo.Valid() {
		return nil, errors.New("invalid user")
	}
	l := req.GetWireguardLink()
	if l == nil || l.GetId() == 0 || l.GetFromWireguardId() == 0 || l.GetToWireguardId() == 0 || l.GetFromWireguardId() == l.GetToWireguardId() {
		return nil, errors.New("invalid link params")
	}
	q := dao.NewQuery(ctx)
	mut := dao.NewMutation(ctx)

	m, err := q.GetWireGuardLinkByID(userInfo, uint(l.GetId()))
	if err != nil {
		return nil, err
	}

	var newEp *models.Endpoint
	if l.GetToEndpoint() != nil && l.GetToEndpoint().GetId() > 0 {
		newEp, err = q.GetEndpointByID(userInfo, uint(l.GetToEndpoint().GetId()))
		if err != nil {
			return nil, err
		}
	}

	// 只能修改这些
	m.Active = l.GetActive()
	m.LatencyMs = l.GetLatencyMs()
	m.UpBandwidthMbps = l.GetUpBandwidthMbps()
	m.DownBandwidthMbps = l.GetDownBandwidthMbps()
	if l.GetToEndpoint() != nil && l.GetToEndpoint().GetId() > 0 {
		m.ToEndpoint = newEp
	}

	if err := mut.UpdateWireGuardLink(userInfo, uint(l.GetId()), m); err != nil {
		return nil, err
	}
	return &pb.UpdateWireGuardLinkResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"}, WireguardLink: m.ToPB()}, nil
}

func DeleteWireGuardLink(ctx *app.Context, req *pb.DeleteWireGuardLinkRequest) (*pb.DeleteWireGuardLinkResponse, error) {
	userInfo := common.GetUserInfo(ctx)
	if !userInfo.Valid() {
		return nil, errors.New("invalid user")
	}
	id := uint(req.GetId())
	if id == 0 {
		return nil, errors.New("invalid id")
	}

	q := dao.NewQuery(ctx)
	m := dao.NewMutation(ctx)

	link, err := q.GetWireGuardLinkByID(userInfo, id)
	if err != nil {
		return nil, err
	}

	rev, err := q.GetWireGuardLinkByClientIDs(userInfo, link.ToWireGuardID, link.FromWireGuardID)
	if err != nil {
		return nil, err
	}

	if err := m.DeleteWireGuardLink(userInfo, uint(link.ID)); err != nil {
		return nil, err
	}

	if err := m.DeleteWireGuardLink(userInfo, uint(rev.ID)); err != nil {
		return nil, err
	}

	return &pb.DeleteWireGuardLinkResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"}}, nil
}

func GetWireGuardLink(ctx *app.Context, req *pb.GetWireGuardLinkRequest) (*pb.GetWireGuardLinkResponse, error) {
	userInfo := common.GetUserInfo(ctx)
	if !userInfo.Valid() {
		return nil, errors.New("invalid user")
	}
	id := uint(req.GetId())
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	m, err := dao.NewQuery(ctx).GetWireGuardLinkByID(userInfo, id)
	if err != nil {
		return nil, err
	}
	return &pb.GetWireGuardLinkResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"}, WireguardLink: m.ToPB()}, nil
}

func ListWireGuardLinks(ctx *app.Context, req *pb.ListWireGuardLinksRequest) (*pb.ListWireGuardLinksResponse, error) {
	userInfo := common.GetUserInfo(ctx)
	if !userInfo.Valid() {
		return nil, errors.New("invalid user")
	}
	page, pageSize := int(req.GetPage()), int(req.GetPageSize())
	keyword := req.GetKeyword()
	networkID := uint(req.GetNetworkId())
	list, err := dao.NewQuery(ctx).ListWireGuardLinksWithFilters(userInfo, page, pageSize, networkID, keyword)
	if err != nil {
		return nil, err
	}
	total, err := dao.NewQuery(ctx).CountWireGuardLinksWithFilters(userInfo, networkID, keyword)
	if err != nil {
		return nil, err
	}
	return &pb.ListWireGuardLinksResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"},
		Total: lo.ToPtr(int32(total)),
		WireguardLinks: lo.Map(list, func(x *models.WireGuardLink, _ int) *pb.WireGuardLink {
			return x.ToPB()
		})}, nil
}
