package wg

import (
	"errors"
	"fmt"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/VaalaCat/frp-panel/services/wg"
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

	if req.GetSpf() {
		// SPF 模式：展示“真实下发的路由表”（即 PeerConfig.AllowedIps），确保与实际一致。
		peerCfgs, allEdges, err := wg.PlanAllowedIPs(peers, links, policy)
		if err != nil {
			log.WithError(err).Errorf("failed to plan allowed ips")
			return nil, err
		}
		adjs := peerConfigsToPBAdjs(peerCfgs, allEdges)

		return &pb.GetNetworkTopologyResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"},
			Adjs:   adjs,
		}, nil
	}

	resp, err := wg.NewDijkstraAllowedIPsPlanner(policy).BuildGraph(peers, links)
	if err != nil {
		log.WithError(err).Errorf("failed to build graph")
		return nil, err
	}
	adjs := adjsToPB(resp)

	return &pb.GetNetworkTopologyResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"},
		Adjs:   adjs,
	}, nil
}
