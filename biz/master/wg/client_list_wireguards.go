package wg

import (
	"sort"

	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	wgsvc "github.com/VaalaCat/frp-panel/services/wg"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
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

			// 构建 network 内 WireGuard 索引，用于补齐“可直连 peer 的基础配置”
			idToWg := make(map[uint32]*models.WireGuard, len(networkPeers[wgCfg.NetworkID]))
			for _, item := range networkPeers[wgCfg.NetworkID] {
				if item == nil {
					continue
				}
				idToWg[uint32(item.ID)] = item
			}

			r := wgCfg.ToPB()
			r.Peers = lo.Map(networkPeerConfigsMap[wgCfg.NetworkID][wgCfg.ID],
				func(peerCfg *pb.WireGuardPeerConfig, _ int) *pb.WireGuardPeerConfig {
					return peerCfg
				})

			r.Adjs = adjsToPB(networkAllEdgesMap[wgCfg.NetworkID])

			fillConnectablePeersAsPreconnect(r, uint32(wgCfg.ID), idToWg, log)
			sortPeersStable(r)

			return r
		}),
	}

	return resp, nil
}

// fillConnectablePeersAsPreconnect 将 adj[localID] 中可直连的 peer 补齐到 r.peers 中，并将 AllowedIPs 置空（只预连接，不承载路由）。
func fillConnectablePeersAsPreconnect(r *pb.WireGuardConfig, localID uint32, idToWg map[uint32]*models.WireGuard, log *logrus.Entry) {
	if r == nil || localID == 0 {
		return
	}
	exists := make(map[uint32]struct{}, len(r.GetPeers()))
	for _, p := range r.GetPeers() {
		if p == nil {
			continue
		}
		if p.GetId() != 0 {
			exists[p.GetId()] = struct{}{}
		}
		if p.GetEndpoint() != nil && p.GetEndpoint().GetWireguardId() != 0 {
			exists[p.GetEndpoint().GetWireguardId()] = struct{}{}
		}
	}

	links := r.GetAdjs()[localID]
	if links == nil {
		return
	}
	for _, l := range links.GetLinks() {
		if l == nil {
			continue
		}
		toID := l.GetToWireguardId()
		if toID == 0 || toID == localID {
			continue
		}
		if _, ok := exists[toID]; ok {
			continue
		}
		remote, ok := idToWg[toID]
		if !ok || remote == nil {
			continue
		}

		// 优先使用链路显式 to_endpoint
		var specifiedEndpoint *models.Endpoint
		if l.GetToEndpoint() != nil {
			m := &models.Endpoint{}
			m.FromPB(l.GetToEndpoint())
			specifiedEndpoint = m
		}

		base, err := remote.AsBasePeerConfig(specifiedEndpoint)
		if err != nil {
			log.WithError(err).Warnf("failed to build base peer config for preconnect: local=%d to=%d", localID, toID)
			continue
		}
		base.AllowedIps = nil
		r.Peers = append(r.Peers, base)
		exists[toID] = struct{}{}
	}
}

func sortPeersStable(r *pb.WireGuardConfig) {
	if r == nil || len(r.Peers) <= 1 {
		return
	}
	sort.SliceStable(r.Peers, func(i, j int) bool {
		pi := r.Peers[i]
		pj := r.Peers[j]
		if pi == nil && pj == nil {
			return false
		}
		if pi == nil {
			return false
		}
		if pj == nil {
			return true
		}
		return pi.GetClientId() < pj.GetClientId()
	})
}
