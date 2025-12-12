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
