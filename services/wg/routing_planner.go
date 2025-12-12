package wg

import (
	"errors"
	"math"
	"sort"
	"time"

	"github.com/samber/lo"

	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
)

// RoutingPolicy 决定边权重的计算方式。
// cost = LatencyWeight*latency_ms + InverseBandwidthWeight*(1/max(up_mbps,1e-6)) + HopWeight
type RoutingPolicy struct {
	LatencyWeight            float64
	InverseBandwidthWeight   float64
	HopWeight                float64
	MinUpMbps                uint32
	DefaultEndpointUpMbps    uint32
	DefaultEndpointLatencyMs uint32
	OfflineThreshold         time.Duration

	ACL                  *ACL
	NetworkTopologyCache app.NetworkTopologyCache
	CliMgr               app.ClientsManager
}

func (p *RoutingPolicy) LoadACL(acl *ACL) *RoutingPolicy {
	p.ACL = acl
	return p
}

func DefaultRoutingPolicy(acl *ACL, networkTopologyCache app.NetworkTopologyCache, cliMgr app.ClientsManager) RoutingPolicy {
	return RoutingPolicy{
		LatencyWeight:            1.0,
		InverseBandwidthWeight:   50.0, // 对低带宽路径给予更高惩罚
		HopWeight:                1.0,
		DefaultEndpointUpMbps:    50,
		DefaultEndpointLatencyMs: 30,
		OfflineThreshold:         2 * time.Minute,
		ACL:                      acl,
		NetworkTopologyCache:     networkTopologyCache,
		CliMgr:                   cliMgr,
	}
}

type AllowedIPsPlanner interface {
	// Compute 基于拓扑与链路指标，计算每个节点应配置到直连邻居的 AllowedIPs。
	// 输入的 peers 应包含同一 Network 下的所有 WireGuard 节点，links 为其有向链路。
	// 返回节点ID->PeerConfig 列表，节点所有 ID->Edge 列表。
	Compute(peers []*models.WireGuard, links []*models.WireGuardLink) (map[uint][]*pb.WireGuardPeerConfig, map[uint][]Edge, error)
	// BuildGraph 基于拓扑与链路指标，计算每个节点应配置到直连邻居的 AllowedIPs，并返回节点ID->Edge 列表。
	BuildGraph(peers []*models.WireGuard, links []*models.WireGuardLink) (map[uint][]Edge, error)
	// BuildFinalGraph 最短路径算法，返回节点ID->Edge 列表。
	BuildFinalGraph(peers []*models.WireGuard, links []*models.WireGuardLink) (map[uint][]Edge, error)
}

type dijkstraAllowedIPsPlanner struct {
	policy RoutingPolicy
}

func NewDijkstraAllowedIPsPlanner(policy RoutingPolicy) AllowedIPsPlanner {
	return &dijkstraAllowedIPsPlanner{policy: policy}
}

func PlanAllowedIPs(peers []*models.WireGuard, links []*models.WireGuardLink, policy RoutingPolicy) (map[uint][]*pb.WireGuardPeerConfig, map[uint][]Edge, error) {
	return NewDijkstraAllowedIPsPlanner(policy).Compute(peers, links)
}

func (p *dijkstraAllowedIPsPlanner) Compute(peers []*models.WireGuard, links []*models.WireGuardLink) (map[uint][]*pb.WireGuardPeerConfig, map[uint][]Edge, error) {
	if len(peers) == 0 {
		return map[uint][]*pb.WireGuardPeerConfig{}, map[uint][]Edge{}, nil
	}

	idToPeer, order := buildNodeIndex(peers)
	adj := buildAdjacency(order, idToPeer, links, p.policy)
	aggByNode, edgeInfoMap := runAllPairsDijkstra(order, adj, idToPeer, p.policy)
	result, err := assemblePeerConfigs(order, aggByNode, edgeInfoMap, idToPeer)
	if err != nil {
		return nil, nil, err
	}
	fillIsolates(order, result)

	// 填充没有链路的节点
	for _, id := range order {
		if _, ok := adj[id]; !ok {
			adj[id] = []Edge{}
		}
	}

	return result, adj, nil
}

func (p *dijkstraAllowedIPsPlanner) BuildGraph(peers []*models.WireGuard, links []*models.WireGuardLink) (map[uint][]Edge, error) {
	idToPeer, order := buildNodeIndex(peers)
	adj := buildAdjacency(order, idToPeer, links, p.policy)
	// 填充没有链路的节点
	for _, id := range order {
		if _, ok := adj[id]; !ok {
			adj[id] = []Edge{}
		}
	}
	return adj, nil
}

func (p *dijkstraAllowedIPsPlanner) BuildFinalGraph(peers []*models.WireGuard, links []*models.WireGuardLink) (map[uint][]Edge, error) {
	idToPeer, order := buildNodeIndex(peers)
	adj := buildAdjacency(order, idToPeer, links, p.policy)
	routesInfoMap, edgeInfoMap := runAllPairsDijkstra(order, adj, idToPeer, p.policy)

	ret := map[uint][]Edge{}
	for src, edgeInfo := range edgeInfoMap {
		for next := range edgeInfo {
			if _, ok := adj[src]; !ok {
				continue
			}
			originEdge := Edge{}
			finded := false
			for _, e := range adj[src] {
				if e.to == next {
					originEdge = e
					finded = true
					break
				}
			}
			if !finded {
				continue
			}

			routesInfo := routesInfoMap[src][next]

			ret[src] = append(ret[src], Edge{
				to:         next,
				latency:    originEdge.latency,
				upMbps:     originEdge.upMbps,
				toEndpoint: originEdge.toEndpoint,
				routes:     lo.Keys(routesInfo),
			})
		}
	}
	for _, id := range order {
		if _, ok := ret[id]; !ok {
			ret[id] = []Edge{}
		}
	}
	return ret, nil
}

type Edge struct {
	to         uint
	latency    uint32
	upMbps     uint32
	toEndpoint *models.Endpoint // 指定的目标端点，可能为 nil
	routes     []string         // 路由信息
}

func (e *Edge) ToPB() *pb.WireGuardLink {
	link := &pb.WireGuardLink{
		ToWireguardId:   uint32(e.to),
		LatencyMs:       e.latency,
		UpBandwidthMbps: e.upMbps,
		Active:          true,
		Routes:          e.routes,
	}
	if e.toEndpoint != nil {
		link.ToEndpoint = e.toEndpoint.ToPB()
	}
	return link
}

func buildNodeIndex(peers []*models.WireGuard) (map[uint]*models.WireGuard, []uint) {
	idToPeer := make(map[uint]*models.WireGuard, len(peers))
	order := make([]uint, 0, len(peers))
	for _, p := range peers {
		idToPeer[uint(p.ID)] = p
		order = append(order, uint(p.ID))
	}
	return idToPeer, order
}

func buildAdjacency(order []uint, idToPeer map[uint]*models.WireGuard, links []*models.WireGuardLink, policy RoutingPolicy) map[uint][]Edge {
	adj := make(map[uint][]Edge, len(order))
	// 1) 显式链路
	for _, l := range links {
		if !l.Active {
			continue
		}
		from := l.FromWireGuardID
		to := l.ToWireGuardID

		if _, ok := idToPeer[from]; !ok {
			continue
		}

		if _, ok := idToPeer[to]; !ok {
			continue
		}

		if lastSeenAt, ok := policy.CliMgr.GetLastSeenAt(idToPeer[from].ClientID); !ok || time.Since(lastSeenAt) > policy.OfflineThreshold {
			continue
		}

		if lastSeenAt, ok := policy.CliMgr.GetLastSeenAt(idToPeer[to].ClientID); !ok || time.Since(lastSeenAt) > policy.OfflineThreshold {
			continue
		}

		// 如果两个peer都没有endpoint，则不建立链路
		if len(idToPeer[from].AdvertisedEndpoints) == 0 && len(idToPeer[to].AdvertisedEndpoints) == 0 {
			continue
		}

		latency := l.LatencyMs
		if latency == 0 { // 如果指定latency为0，则使用真实值
			if latencyMs, ok := policy.NetworkTopologyCache.GetLatencyMs(from, to); ok {
				latency = latencyMs
			} else {
				latency = policy.DefaultEndpointLatencyMs
			}
		}

		adj[from] = append(adj[from], Edge{to: to, latency: latency, upMbps: l.UpBandwidthMbps, toEndpoint: l.ToEndpoint})
	}

	// 2) 若某节点具备 endpoint，则所有其他节点可直连它
	edgeSet := make(map[[2]uint]struct{}, 16)
	for from, edges := range adj {
		for _, e := range edges { // 先拿到所有直连的节点
			edgeSet[[2]uint{from, e.to}] = struct{}{}
			edgeSet[[2]uint{e.to, from}] = struct{}{}
		}
	}

	for _, to := range order {
		peer := idToPeer[to]
		if peer == nil || len(peer.AdvertisedEndpoints) == 0 {
			continue
		}
		for _, from := range order {
			if from == to {
				continue
			}
			if _, ok := idToPeer[from]; !ok {
				continue
			}

			latency := policy.DefaultEndpointLatencyMs
			if latencyMs, ok := policy.NetworkTopologyCache.GetLatencyMs(from, to); ok {
				latency = latencyMs
			}
			if latencyMs, ok := policy.NetworkTopologyCache.GetLatencyMs(to, from); ok {
				latency = latencyMs
			}

			if lastSeenAt, ok := policy.CliMgr.GetLastSeenAt(idToPeer[from].ClientID); !ok || time.Since(lastSeenAt) > policy.OfflineThreshold {
				continue
			}

			if lastSeenAt, ok := policy.CliMgr.GetLastSeenAt(idToPeer[to].ClientID); !ok || time.Since(lastSeenAt) > policy.OfflineThreshold {
				continue
			}

			// 有 acl 限制
			if policy.ACL.CanConnect(idToPeer[from], idToPeer[to]) {
				key1 := [2]uint{from, to}
				if _, exists := edgeSet[key1]; exists {
					continue
				}

				adj[from] = append(adj[from], Edge{to: to, latency: latency, upMbps: policy.DefaultEndpointUpMbps})
				edgeSet[key1] = struct{}{}
			}

			if policy.ACL.CanConnect(idToPeer[to], idToPeer[from]) {
				key2 := [2]uint{to, from}
				if _, exists := edgeSet[key2]; exists {
					continue
				}
				adj[to] = append(adj[to], Edge{to: from, latency: latency, upMbps: policy.DefaultEndpointUpMbps})
				edgeSet[key2] = struct{}{}
			}
		}
	}
	return adj
}

// EdgeInfo 保存边的端点信息，用于后续组装 PeerConfig
type EdgeInfo struct {
	toEndpoint *models.Endpoint
}

// runAllPairsDijkstra returns: map[src]map[nextHop]map[CIDR], map[src]map[nextHop]*EdgeInfo
func runAllPairsDijkstra(order []uint, adj map[uint][]Edge, idToPeer map[uint]*models.WireGuard, policy RoutingPolicy) (map[uint]map[uint]map[string]struct{}, map[uint]map[uint]*EdgeInfo) {
	aggByNode := make(map[uint]map[uint]map[string]struct{}, len(order))
	edgeInfoMap := make(map[uint]map[uint]*EdgeInfo, len(order)) // 保存 src -> nextHop 的边信息

	for _, src := range order {
		dist, prev, visited := initSSSP(order)
		dist[src] = 0

		for {
			u, ok := pickNext(order, dist, visited)
			if !ok {
				break
			}
			visited[u] = true
			for _, e := range adj[u] {
				invBw := 1.0 / math.Max(float64(e.upMbps), 1e-6)
				w := policy.LatencyWeight*float64(e.latency) + policy.InverseBandwidthWeight*invBw + policy.HopWeight
				alt := dist[u] + w
				if alt < dist[e.to] {
					dist[e.to] = alt
					prev[e.to] = u
				}
			}
		}

		// 累计 nextHop -> CIDR，并保存边信息
		for _, dst := range order {
			if dst == src {
				continue
			}
			if _, ok := prev[dst]; !ok {
				continue
			}
			next := findNextHop(src, dst, prev)
			if next == 0 {
				continue
			}
			dstPeer := idToPeer[dst]
			allowed, err := dstPeer.AsBasePeerConfig(nil) // 这里只获取 CIDR，不需要指定 endpoint
			if err != nil || len(allowed.GetAllowedIps()) == 0 {
				continue
			}
			cidr := allowed.GetAllowedIps()[0]
			if _, ok := aggByNode[src]; !ok {
				aggByNode[src] = make(map[uint]map[string]struct{})
			}
			if _, ok := aggByNode[src][next]; !ok {
				aggByNode[src][next] = map[string]struct{}{}
			}
			aggByNode[src][next][cidr] = struct{}{}

			// 保存从 src 到 next 的边信息（查找直接边）
			if _, ok := edgeInfoMap[src]; !ok {
				edgeInfoMap[src] = make(map[uint]*EdgeInfo)
			}
			if _, ok := edgeInfoMap[src][next]; !ok {
				// 查找从 src 到 next 的边
				for _, e := range adj[src] {
					if e.to == next {
						edgeInfoMap[src][next] = &EdgeInfo{toEndpoint: e.toEndpoint}
						break
					}
				}
			}
		}
	}
	return aggByNode, edgeInfoMap
}

func initSSSP(order []uint) (map[uint]float64, map[uint]uint, map[uint]bool) {
	dist := make(map[uint]float64, len(order))
	prev := make(map[uint]uint, len(order))
	visited := make(map[uint]bool, len(order))
	for _, vid := range order {
		dist[vid] = math.Inf(1)
	}
	return dist, prev, visited
}

func pickNext(order []uint, dist map[uint]float64, visited map[uint]bool) (uint, bool) {
	best := uint(0)
	bestVal := math.Inf(1)
	found := false
	for _, vid := range order {
		if visited[vid] {
			continue
		}
		if dist[vid] < bestVal {
			bestVal = dist[vid]
			best = vid
			found = true
		}
	}
	return best, found
}

func findNextHop(src, dst uint, prev map[uint]uint) uint {
	next := dst
	for {
		p, ok := prev[next]
		if !ok {
			return 0
		}
		if p == src {
			return next
		}
		next = p
	}
}

func assemblePeerConfigs(order []uint, aggByNode map[uint]map[uint]map[string]struct{}, edgeInfoMap map[uint]map[uint]*EdgeInfo, idToPeer map[uint]*models.WireGuard) (map[uint][]*pb.WireGuardPeerConfig, error) {
	result := make(map[uint][]*pb.WireGuardPeerConfig, len(order))
	for src, nextMap := range aggByNode {
		peersForSrc := make([]*pb.WireGuardPeerConfig, 0, len(nextMap))
		for nextHop, cidrSet := range nextMap {
			remote := idToPeer[nextHop]

			// 获取从 src 到 nextHop 的边信息，确定使用哪个 endpoint
			var specifiedEndpoint *models.Endpoint
			if edgeInfo, ok := edgeInfoMap[src][nextHop]; ok && edgeInfo != nil && edgeInfo.toEndpoint != nil {
				specifiedEndpoint = edgeInfo.toEndpoint
			}

			base, err := remote.AsBasePeerConfig(specifiedEndpoint)
			if err != nil {
				return nil, errors.Join(errors.New("build peer base config failed"), err)
			}
			cidrs := make([]string, 0, len(cidrSet))
			for c := range cidrSet {
				cidrs = append(cidrs, c)
			}
			sort.Strings(cidrs)
			base.AllowedIps = lo.Uniq(cidrs)
			peersForSrc = append(peersForSrc, base)
		}
		sort.SliceStable(peersForSrc, func(i, j int) bool {
			return peersForSrc[i].GetClientId() < peersForSrc[j].GetClientId()
		})
		result[src] = peersForSrc
	}
	return result, nil
}

func fillIsolates(order []uint, result map[uint][]*pb.WireGuardPeerConfig) {
	for _, id := range order {
		if _, ok := result[id]; !ok {
			result[id] = []*pb.WireGuardPeerConfig{}
		}
	}
}
