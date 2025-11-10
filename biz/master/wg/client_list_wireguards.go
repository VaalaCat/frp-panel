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

	resp := &pb.ListClientWireGuardsResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"},
		WireguardConfigs: lo.Map(wgCfgs, func(wgCfg *models.WireGuard, _ int) *pb.WireGuardConfig {
			if wgCfg == nil || wgCfg.Network == nil {
				log.Warnf("wireguard config or network is nil, wireguard id: %d", wgCfg.ID)
				return nil
			}

			peerConfigs, err := wgsvc.PlanAllowedIPs(
				networkPeers[wgCfg.NetworkID], networkLinksMap[wgCfg.NetworkID],
				wgsvc.DefaultRoutingPolicy(
					wgsvc.NewACL().LoadFromPB(wgCfg.Network.ACL.Data),
					ctx.GetApp().GetNetworkTopologyCache(),
					ctx.GetApp().GetClientsManager(),
				))
			if err != nil {
				log.WithError(err).Errorf("failed to plan allowed ips for wireguard configs: %v", wgCfgs)
				return nil
			}

			r := wgCfg.ToPB()
			r.Peers = lo.Map(peerConfigs[wgCfg.ID], func(peerCfg *pb.WireGuardPeerConfig, _ int) *pb.WireGuardPeerConfig {
				return peerCfg
			})

			return r
		}),
	}

	return resp, nil
}
