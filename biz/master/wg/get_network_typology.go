package wg

import (
	"errors"
	"fmt"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/VaalaCat/frp-panel/services/wg"
	"github.com/samber/lo"
)

func GetNetworkTopology(ctx *app.Context, req *pb.GetNetworkTopologyRequest) (*pb.GetNetworkTopologyResponse, error) {
	log := ctx.Logger().WithField("op", "GetNetworkTopology")

	userInfo := common.GetUserInfo(ctx)
	if !userInfo.Valid() {
		return nil, errors.New("invalid user")
	}

	networkID := uint(req.GetId())
	if networkID == 0 {
		return nil, errors.New("invalid id")
	}

	q := dao.NewQuery(ctx)

	peers, err := q.GetWireGuardsByNetworkID(userInfo, networkID)
	if err != nil {
		log.WithError(err).Errorf("failed to get wireguard peers by network id: %d", networkID)
		return nil, err
	}
	links, err := q.ListWireGuardLinksByNetwork(userInfo, networkID)
	if err != nil {
		log.WithError(err).Errorf("failed to get wireguard links by network id: %d", networkID)
		return nil, err
	}

	if len(peers) == 0 {
		log.Errorf("no wireguard peers found")
		return nil, fmt.Errorf("no wireguard peers found")
	}

	policy := wg.DefaultRoutingPolicy(
		wg.NewACL().LoadFromPB(peers[0].Network.ACL.Data),
		ctx.GetApp().GetNetworkTopologyCache(),
		ctx.GetApp().GetClientsManager(),
	)

	var resp map[uint][]wg.Edge

	if req.GetSpf() {
		resp, err = wg.NewDijkstraAllowedIPsPlanner(policy).BuildFinalGraph(peers, links)
	} else {
		resp, err = wg.NewDijkstraAllowedIPsPlanner(policy).BuildGraph(peers, links)
	}

	if err != nil {
		log.WithError(err).Errorf("failed to build graph")
		return nil, err
	}

	adjs := make(map[uint32]*pb.WireGuardLinks)
	for id, peerConfigs := range resp {
		adjs[uint32(id)] = &pb.WireGuardLinks{
			Links: lo.Map(peerConfigs, func(e wg.Edge, _ int) *pb.WireGuardLink {
				v := e.ToPB()
				v.FromWireguardId = uint32(uint(id))
				return v
			}),
		}
	}

	for id, links := range adjs {
		for _, link := range links.GetLinks() {
			toWireguardEdges, ok := adjs[uint32(link.GetToWireguardId())]
			if ok {
				for _, edge := range toWireguardEdges.GetLinks() {
					if edge.GetToWireguardId() == uint32(uint(id)) {
						link.DownBandwidthMbps = edge.GetUpBandwidthMbps()
					}
				}
			}
		}
	}

	return &pb.GetNetworkTopologyResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"},
		Adjs:   adjs,
	}, nil
}
