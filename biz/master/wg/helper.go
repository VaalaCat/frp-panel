package wg

import (
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/wg"
	"github.com/samber/lo"
)

func adjsToPB(resp map[uint][]wg.Edge) map[uint32]*pb.WireGuardLinks {
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

	return adjs
}

// peerConfigsToPBAdjs 将“真实下发给节点的路由表（PeerConfig.AllowedIps）”转换为拓扑展示所需的 Adjs。
//
// - routes: 直接使用 peerCfg.AllowedIps（这才是 WireGuard 实际使用的路由表）
// - latency/up/down: 尽量从 allEdges（buildAdjacency 的直连边指标）中补齐，仅用于展示
// - endpoint: 优先使用 peerCfg.Endpoint（与实际下发一致）
func peerConfigsToPBAdjs(peerCfgs map[uint][]*pb.WireGuardPeerConfig, allEdges map[uint][]wg.Edge) map[uint32]*pb.WireGuardLinks {
	adjs := make(map[uint32]*pb.WireGuardLinks, len(peerCfgs))

	for src, pcs := range peerCfgs {
		srcID := uint32(src)
		links := make([]*pb.WireGuardLink, 0, len(pcs))

		// 构建 toID -> edge 指标索引（仅用于展示 latency/up）
		edgePBByTo := make(map[uint32]*pb.WireGuardLink, 16)
		for _, e := range allEdges[src] {
			epb := e.ToPB()
			edgePBByTo[epb.GetToWireguardId()] = epb
		}

		for _, pc := range pcs {
			if pc == nil || pc.GetId() == 0 {
				continue
			}
			toID := pc.GetId()
			var latency uint32
			var up uint32
			if epb, ok := edgePBByTo[toID]; ok && epb != nil {
				latency = epb.GetLatencyMs()
				up = epb.GetUpBandwidthMbps()
			}

			links = append(links, &pb.WireGuardLink{
				FromWireguardId:   srcID,
				ToWireguardId:     toID,
				LatencyMs:         latency,
				UpBandwidthMbps:   up,
				DownBandwidthMbps: 0, // 下面统一填充
				Active:            true,
				ToEndpoint:        pc.GetEndpoint(),
				Routes:            pc.GetAllowedIps(),
			})
		}

		adjs[srcID] = &pb.WireGuardLinks{Links: links}
	}

	// 填充 down bandwidth（参考 adjsToPB 的做法：取反向边的 up）
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

	return adjs
}
