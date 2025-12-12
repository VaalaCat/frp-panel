package wg

import (
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	wgsvc "github.com/VaalaCat/frp-panel/services/wg"
	"github.com/samber/lo"
)

func ListClientWireGuards(ctx *app.Context, req *pb.ListClientWireGuardsRequest) (*pb.ListClientWireGuardsResponse, error) {
	clientId := req.GetBase().GetClientId()
	log := ctx.Logger().WithField("op", "ListClientWireGuards")

	wgCfgs, err := dao.NewQuery(ctx).AdminListWireGuardsWithClientID(clientId)
	if err != nil {
		log.WithError(err).Errorf("failed to list wireguard configs with client id: %s", clientId)
		return nil, err
	}

	networkPeers := map[uint][]*models.WireGuard{}
	networkIDs := lo.Map(wgCfgs, func(wgCfg *models.WireGuard, _ int) uint {
		return wgCfg.NetworkID
	})

	allRelatedWgCfgs, err := dao.NewQuery(ctx).AdminListWireGuardsWithNetworkIDs(networkIDs)
	if err != nil {
		log.WithError(err).Errorf("failed to list wireguard configs with network ids: %v", networkIDs)
		return nil, err
	}

	allRelatedWgLinks, err := dao.NewQuery(ctx).AdminListWireGuardLinksWithNetworkIDs(networkIDs)
	if err != nil {
		log.WithError(err).Errorf("failed to list wireguard links with network ids: %v", networkIDs)
		return nil, err
	}
	networkLinksMap := make(map[uint][]*models.WireGuardLink)
	for _, link := range allRelatedWgLinks {
		if _, ok := networkLinksMap[link.NetworkID]; !ok {
			networkLinksMap[link.NetworkID] = []*models.WireGuardLink{}
		}
		networkLinksMap[link.NetworkID] = append(networkLinksMap[link.NetworkID], link)
	}

	for _, wgCfg := range allRelatedWgCfgs {
		if _, ok := networkPeers[wgCfg.NetworkID]; !ok {
			networkPeers[wgCfg.NetworkID] = []*models.WireGuard{}
		}
		networkPeers[wgCfg.NetworkID] = append(networkPeers[wgCfg.NetworkID], wgCfg)
	}

	networkPeerConfigsMap := make(map[uint]map[uint][]*pb.WireGuardPeerConfig)
	networkAllEdgesMap := make(map[uint]map[uint][]wgsvc.Edge)

	for _, networkID := range networkIDs {
		peerConfigs, allEdges, err := wgsvc.PlanAllowedIPs(
			networkPeers[networkID], networkLinksMap[networkID],
			wgsvc.DefaultRoutingPolicy(
				wgsvc.NewACL().LoadFromPB(networkPeers[networkID][0].Network.ACL.Data),
				ctx.GetApp().GetNetworkTopologyCache(),
				ctx.GetApp().GetClientsManager(),
			))

		if err != nil {
			log.WithError(err).Errorf("failed to plan allowed ips for wireguard configs: %v", wgCfgs)
			return nil, err
		}

		networkPeerConfigsMap[networkID] = peerConfigs
		networkAllEdgesMap[networkID] = allEdges
	}

	resp := &pb.ListClientWireGuardsResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"},
		WireguardConfigs: lo.Map(wgCfgs, func(wgCfg *models.WireGuard, _ int) *pb.WireGuardConfig {
			if wgCfg == nil || wgCfg.Network == nil {
				log.Warnf("wireguard config or network is nil, wireguard id: %d", wgCfg.ID)
				return nil
			}

			r := wgCfg.ToPB()
			r.Peers = lo.Map(networkPeerConfigsMap[wgCfg.NetworkID][wgCfg.ID],
				func(peerCfg *pb.WireGuardPeerConfig, _ int) *pb.WireGuardPeerConfig {
					return peerCfg
				})

			r.Adjs = adjsToPB(networkAllEdgesMap[wgCfg.NetworkID])

			return r
		}),
	}

	return resp, nil
}
