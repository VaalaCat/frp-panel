package wg

import (
	"errors"
	"fmt"
	"math"
	"net/netip"
	"sort"
	"time"

	"github.com/samber/lo"

	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
)

// WireGuard 的 AllowedIPs 同时承担两件事：
// 1) 出站选路：目的 IP 匹配哪个 peer 的 AllowedIPs，就把包发给哪个 peer
// 2) 入站源地址校验：从某 peer 解密出来的 inner packet，其 source IP 必须落在该 peer 的 AllowedIPs
//
// 因此，多跳转发时，某节点 i 从“上一跳 peer=j”收到的包，其 inner source 仍是“原始源节点 s 的 /32”，
// 所以 i 配置 peer(j) 的 AllowedIPs 必须包含这些会经由 j 转发进来的“源地址集合”，否则会直接丢包。
// 思路：
// - 在一个“对称权重”的图上做最短路（保证路径可逆，避免重复/冲突）
// - 同时产出：
//   - Out(i->nextHop): i 出站时，哪些目的 /32 应走 nextHop（目的集合）
//   - In(i<-prevHop): i 入站时，从 prevHop 过来的包允许哪些源 /32（源集合）
// - 最终对每个 i 的每个直连 peer(j)，AllowedIPs = Out(i->j) ∪ In(i<-j)
// - 严格校验：对同一节点 i，不允许出现同一个 /32 同时出现在多个 peer 的 AllowedIPs（否则 WG 行为不确定）

type AllowedIPsPlanner interface {
	// Compute 基于拓扑与链路指标，计算每个节点应配置到直连邻居的 AllowedIPs。
	// 输入的 peers 应包含同一 Network 下的所有 WireGuard 节点，links 为其有向链路。
	// 返回：节点ID->PeerConfig 列表，节点ID->Edge 列表（完整候选图，用于展示）。
	Compute(peers []*models.WireGuard, links []*models.WireGuardLink) (map[uint][]*pb.WireGuardPeerConfig, map[uint][]Edge, error)
	// BuildGraph 基于拓扑与链路指标，返回完整候选图（用于展示/诊断）。
	BuildGraph(peers []*models.WireGuard, links []*models.WireGuardLink) (map[uint][]Edge, error)
	// BuildFinalGraph 返回“最终下发的直连边”与其 routes（用于展示 SPF 结果）。
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

	idToPeer, order := buildNodeIndexSorted(peers)
	cidrByID, err := buildNodeCIDRMap(order, idToPeer)
	if err != nil {
		return nil, nil, err
	}

	adj := buildAdjacency(order, idToPeer, links, p.policy)
	// SPF 参与的边：显式边 + 已探测可达的推断边，并要求“可用于转发”的边必须双向存在
	spfAdj := filterAdjacencyForSPF(order, adj, p.policy)
	spfAdj = filterAdjacencyForSymmetricLinks(order, spfAdj)

	peerCfgs, finalEdges, err := computeAllowedIPs(order, idToPeer, cidrByID, spfAdj, adj, p.policy)
	if err != nil {
		return nil, nil, err
	}

	// 填充没有链路的节点（展示用）
	for _, id := range order {
		if _, ok := adj[id]; !ok {
			adj[id] = []Edge{}
		}
		if _, ok := finalEdges[id]; !ok {
			finalEdges[id] = []Edge{}
		}
		if _, ok := peerCfgs[id]; !ok {
			peerCfgs[id] = []*pb.WireGuardPeerConfig{}
		}
	}

	return peerCfgs, adj, nil
}

func (p *dijkstraAllowedIPsPlanner) BuildGraph(peers []*models.WireGuard, links []*models.WireGuardLink) (map[uint][]Edge, error) {
	if len(peers) == 0 {
		return map[uint][]Edge{}, nil
	}
	idToPeer, order := buildNodeIndexSorted(peers)
	adj := buildAdjacency(order, idToPeer, links, p.policy)
	for _, id := range order {
		if _, ok := adj[id]; !ok {
			adj[id] = []Edge{}
		}
	}
	return adj, nil
}

func (p *dijkstraAllowedIPsPlanner) BuildFinalGraph(peers []*models.WireGuard, links []*models.WireGuardLink) (map[uint][]Edge, error) {
	if len(peers) == 0 {
		return map[uint][]Edge{}, nil
	}

	idToPeer, order := buildNodeIndexSorted(peers)
	cidrByID, err := buildNodeCIDRMap(order, idToPeer)
	if err != nil {
		return nil, err
	}

	adj := buildAdjacency(order, idToPeer, links, p.policy)
	spfAdj := filterAdjacencyForSPF(order, adj, p.policy)
	spfAdj = filterAdjacencyForSymmetricLinks(order, spfAdj)

	_, finalEdges, err := computeAllowedIPs(order, idToPeer, cidrByID, spfAdj, adj, p.policy)
	if err != nil {
		return nil, err
	}
	for _, id := range order {
		if _, ok := finalEdges[id]; !ok {
			finalEdges[id] = []Edge{}
		}
	}
	return finalEdges, nil
}

// Edge 表示候选/最终图里的“有向直连边”。
type Edge struct {
	to         uint
	latency    uint32
	upMbps     uint32
	toEndpoint *models.Endpoint // 指定的目标端点，可能为 nil
	routes     []string         // 最终展示：该直连 peer 承载的路由（AllowedIPs）
	explicit   bool             // true: 显式 link；false: 推断/探测用 link
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

// buildNodeIndexSorted 返回：id->peer 映射 与 按 id 排序的 order（用于确定性）
func buildNodeIndexSorted(peers []*models.WireGuard) (map[uint]*models.WireGuard, []uint) {
	idToPeer := make(map[uint]*models.WireGuard, len(peers))
	order := make([]uint, 0, len(peers))
	for _, p := range peers {
		if p == nil {
			continue
		}
		id := uint(p.ID)
		idToPeer[id] = p
		order = append(order, id)
	}
	sort.Slice(order, func(i, j int) bool { return order[i] < order[j] })
	return idToPeer, order
}

func buildNodeCIDRMap(order []uint, idToPeer map[uint]*models.WireGuard) (map[uint]string, error) {
	out := make(map[uint]string, len(order))
	for _, id := range order {
		p := idToPeer[id]
		if p == nil {
			continue
		}
		base, err := p.AsBasePeerConfig(nil)
		if err != nil || len(base.GetAllowedIps()) == 0 {
			return nil, fmt.Errorf("invalid wireguard local address for id=%d", id)
		}
		out[id] = base.GetAllowedIps()[0]
	}
	return out, nil
}

// buildAdjacency 构建“候选直连边”：
// 1) 显式链路（管理员配置）直接加入
// 2) 若某节点具备 endpoint，则其他节点可按 ACL 推断直连它（用于探测/候选）
func buildAdjacency(order []uint, idToPeer map[uint]*models.WireGuard, links []*models.WireGuardLink, policy RoutingPolicy) map[uint][]Edge {
	adj := make(map[uint][]Edge, len(order))

	online := func(id uint) bool {
		if policy.CliMgr == nil {
			return true
		}
		p := idToPeer[id]
		if p == nil || p.ClientID == "" {
			return false
		}
		lastSeenAt, ok := policy.CliMgr.GetLastSeenAt(p.ClientID)
		if !ok {
			return false
		}
		if policy.OfflineThreshold > 0 && time.Since(lastSeenAt) > policy.OfflineThreshold {
			return false
		}
		return true
	}

	// 1) 显式链路
	for _, l := range links {
		if l == nil || !l.Active {
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
		if !online(from) || !online(to) {
			continue
		}
		// 如果两个 peer 都没有 endpoint，则不建立链路（无法直连）
		if len(idToPeer[from].AdvertisedEndpoints) == 0 && len(idToPeer[to].AdvertisedEndpoints) == 0 {
			continue
		}

		latency := l.LatencyMs
		if latency == 0 {
			if policy.NetworkTopologyCache != nil {
				if latencyMs, ok := policy.NetworkTopologyCache.GetLatencyMs(from, to); ok {
					latency = latencyMs
				}
			}
			if latency == 0 {
				latency = policy.DefaultEndpointLatencyMs
			}
		}

		adj[from] = append(adj[from], Edge{
			to:         to,
			latency:    latency,
			upMbps:     l.UpBandwidthMbps,
			toEndpoint: l.ToEndpoint,
			explicit:   true,
		})
	}

	// 2) 推断/探测用边：若某节点具备 endpoint，则所有其他节点可直连它
	edgeSet := make(map[[2]uint]struct{}, 64)
	for from, edges := range adj {
		for _, e := range edges {
			edgeSet[[2]uint{from, e.to}] = struct{}{}
		}
	}

	for _, to := range order {
		peerTo := idToPeer[to]
		if peerTo == nil || len(peerTo.AdvertisedEndpoints) == 0 {
			continue
		}
		for _, from := range order {
			if from == to {
				continue
			}
			if _, ok := idToPeer[from]; !ok {
				continue
			}
			if !online(from) || !online(to) {
				continue
			}

			latency := policy.DefaultEndpointLatencyMs
			if policy.NetworkTopologyCache != nil {
				if latencyMs, ok := policy.NetworkTopologyCache.GetLatencyMs(from, to); ok {
					latency = latencyMs
				}
			}

			// 注意：推断边需要按“两个方向”分别判断 ACL 并分别建边。
			// 这样即使 from 没有 endpoint，也能被 endpoint 节点纳入邻接（满足对称直连 peer 的要求）。

			// from -> to
			if policy.ACL == nil || policy.ACL.CanConnect(idToPeer[from], idToPeer[to]) {
				key := [2]uint{from, to}
				if _, exists := edgeSet[key]; !exists {
					adj[from] = append(adj[from], Edge{
						to:       to,
						latency:  latency,
						upMbps:   policy.DefaultEndpointUpMbps,
						explicit: false,
					})
					edgeSet[key] = struct{}{}
				}
			}

			// to -> from（反向边同样使用同一对节点的 latency 估计；GetLatencyMs 本身已做正反向兜底）
			if policy.ACL == nil || policy.ACL.CanConnect(idToPeer[to], idToPeer[from]) {
				key := [2]uint{to, from}
				if _, exists := edgeSet[key]; !exists {
					adj[to] = append(adj[to], Edge{
						to:       from,
						latency:  latency,
						upMbps:   policy.DefaultEndpointUpMbps,
						explicit: false,
					})
					edgeSet[key] = struct{}{}
				}
			}
		}
	}

	// 稳定排序：保证遍历顺序确定性
	for _, from := range order {
		if edges, ok := adj[from]; ok {
			sort.SliceStable(edges, func(i, j int) bool {
				if edges[i].explicit != edges[j].explicit {
					return edges[i].explicit // explicit 优先
				}
				return edges[i].to < edges[j].to
			})
			adj[from] = edges
		}
	}

	return adj
}

func isUnreachableLatency(latency uint32) bool {
	// 兼容两类不可达哨兵：
	// - math.MaxUint32（历史实现）
	// - math.MaxInt32（部分展示/转换链路里会出现 2147483647）
	return latency == math.MaxUint32 || latency == uint32(math.MaxInt32)
}

// filterAdjacencyForSPF：显式边直接保留；推断边必须有探测数据，且不可达哨兵值剔除
func filterAdjacencyForSPF(order []uint, adj map[uint][]Edge, policy RoutingPolicy) map[uint][]Edge {
	ret := make(map[uint][]Edge, len(order))
	for from, edges := range adj {
		for _, e := range edges {
			if e.explicit {
				ret[from] = append(ret[from], e)
				continue
			}
			if policy.NetworkTopologyCache == nil {
				continue
			}
			latency, ok := policy.NetworkTopologyCache.GetLatencyMs(from, e.to)
			if !ok {
				continue
			}
			if isUnreachableLatency(latency) {
				continue
			}
			e.latency = latency
			ret[from] = append(ret[from], e)
		}
	}
	for _, id := range order {
		if _, ok := ret[id]; !ok {
			ret[id] = []Edge{}
		}
	}
	return ret
}

// filterAdjacencyForSymmetricLinks 仅保留“存在反向直连边”的邻接（用于可转发 SPF）。
func filterAdjacencyForSymmetricLinks(order []uint, adj map[uint][]Edge) map[uint][]Edge {
	ret := make(map[uint][]Edge, len(order))
	edgeSet := make(map[[2]uint]struct{}, 64)
	for from, edges := range adj {
		for _, e := range edges {
			edgeSet[[2]uint{from, e.to}] = struct{}{}
		}
	}
	for from, edges := range adj {
		for _, e := range edges {
			if _, ok := edgeSet[[2]uint{e.to, from}]; !ok {
				continue
			}
			ret[from] = append(ret[from], e)
		}
	}
	for _, id := range order {
		if _, ok := ret[id]; !ok {
			ret[id] = []Edge{}
		}
	}
	return ret
}

type directedEdgeInfo struct {
	latency    uint32
	upMbps     uint32
	toEndpoint *models.Endpoint
	explicit   bool
}

type undirectedNeighbor struct {
	to     uint
	weight float64
}

// computeAllowedIPs 是“最终下发路由”的核心：
// - 在 spfAdj 上构建“对称权重的无向图”
// - 对每个 src 做一次 Dijkstra，得到最短路树 prev
// - 同时生成 Out(dst prefixes) 与 In(src prefixes) 并合并到每条直连 peer 的 AllowedIPs
func computeAllowedIPs(
	order []uint,
	idToPeer map[uint]*models.WireGuard,
	cidrByID map[uint]string,
	spfAdj map[uint][]Edge,
	fullAdj map[uint][]Edge, // 用于展示补齐 latency/up/endpoint
	policy RoutingPolicy,
) (map[uint][]*pb.WireGuardPeerConfig, map[uint][]Edge, error) {
	// 构建 directed edge info（用于 endpoint/展示），并构建 undirected graph（对称权重）
	dInfo := make(map[[2]uint]*directedEdgeInfo, 128)
	undir := make(map[uint][]undirectedNeighbor, len(order))

	// 先把 spfAdj 的 directed info 记下来
	for _, from := range order {
		for _, e := range spfAdj[from] {
			key := [2]uint{from, e.to}
			dInfo[key] = &directedEdgeInfo{
				latency:    e.latency,
				upMbps:     e.upMbps,
				toEndpoint: e.toEndpoint,
				explicit:   e.explicit,
			}
		}
	}

	// 无向图：只添加“成对存在”的边，weight 用 max(w_uv, w_vu) 保证对称
	added := make(map[[2]uint]struct{}, 128)
	for _, u := range order {
		for _, e := range spfAdj[u] {
			v := e.to
			if u == v {
				continue
			}
			// 只处理一次 pair(u,v)
			a, b := u, v
			if a > b {
				a, b = b, a
			}
			pair := [2]uint{a, b}
			if _, ok := added[pair]; ok {
				continue
			}
			// 需要双向边信息
			uv, ok1 := dInfo[[2]uint{u, v}]
			vu, ok2 := dInfo[[2]uint{v, u}]
			if !ok1 || !ok2 || uv == nil || vu == nil {
				continue
			}
			// 用 policy.EdgeWeight 计算双向权重并取 max 做对称
			wuv := policy.EdgeWeight(u, Edge{to: v, latency: uv.latency, upMbps: uv.upMbps, toEndpoint: uv.toEndpoint, explicit: uv.explicit}, idToPeer)
			wvu := policy.EdgeWeight(v, Edge{to: u, latency: vu.latency, upMbps: vu.upMbps, toEndpoint: vu.toEndpoint, explicit: vu.explicit}, idToPeer)
			w := math.Max(wuv, wvu)
			undir[a] = append(undir[a], undirectedNeighbor{to: b, weight: w})
			undir[b] = append(undir[b], undirectedNeighbor{to: a, weight: w})
			added[pair] = struct{}{}
		}
	}

	// 稳定排序
	for _, u := range order {
		neis := undir[u]
		sort.SliceStable(neis, func(i, j int) bool { return neis[i].to < neis[j].to })
		undir[u] = neis
	}

	// Out/ In 聚合：owner -> peer -> set[cidr]
	allowed := make(map[uint]map[uint]map[string]struct{}, len(order))

	for _, src := range order {
		dist := make(map[uint]float64, len(order))
		prev := make(map[uint]uint, len(order)) // prev[dst] = predecessor of dst on path from src
		visited := make(map[uint]bool, len(order))
		for _, id := range order {
			dist[id] = math.Inf(1)
		}
		dist[src] = 0

		// Dijkstra（O(n^2)，节点数通常不大；同时保证确定性）
		for {
			u, ok := pickNext(order, dist, visited)
			if !ok {
				break
			}
			visited[u] = true
			for _, nb := range undir[u] {
				v := nb.to
				if visited[v] {
					continue
				}
				alt := dist[u] + nb.weight
				if alt < dist[v] {
					dist[v] = alt
					prev[v] = u
					continue
				}
				// tie-break：相同距离时，选择更小的 predecessor，确保稳定
				if alt == dist[v] {
					if cur, ok := prev[v]; !ok || u < cur {
						prev[v] = u
					}
				}
			}
		}

		// 1) 出站目的集合：dstCIDR -> nextHop(src,dst)
		for _, dst := range order {
			if dst == src {
				continue
			}
			if _, ok := prev[dst]; !ok {
				continue // unreachable
			}
			next := findNextHop(src, dst, prev)
			if next == 0 {
				continue
			}
			cidr := cidrByID[dst]
			if cidr == "" {
				continue
			}
			ensureAllowedSet(allowed, src, next)[cidr] = struct{}{}
		}

		// 2) 入站源集合：srcCIDR -> prevHop(src,dst) 归到 dst 节点的 peer(prevHop)
		srcCIDR := cidrByID[src]
		if srcCIDR != "" {
			for _, dst := range order {
				if dst == src {
					continue
				}
				pred, ok := prev[dst]
				if !ok || pred == 0 {
					continue
				}
				ensureAllowedSet(allowed, dst, pred)[srcCIDR] = struct{}{}
			}
		}
	}

	// 构建 PeerConfigs，并做强校验（同一节点不允许 CIDR 分配到多个 peer）
	result := make(map[uint][]*pb.WireGuardPeerConfig, len(order))
	finalEdges := make(map[uint][]Edge, len(order))

	for _, owner := range order {
		peerToCIDRs := allowed[owner]
		if len(peerToCIDRs) == 0 {
			result[owner] = []*pb.WireGuardPeerConfig{}
			finalEdges[owner] = []Edge{}
			continue
		}

		seen := make(map[string]uint, 128)
		peerIDs := lo.Keys(peerToCIDRs)
		sort.Slice(peerIDs, func(i, j int) bool { return peerIDs[i] < peerIDs[j] })

		pcs := make([]*pb.WireGuardPeerConfig, 0, len(peerIDs))
		edges := make([]Edge, 0, len(peerIDs))

		for _, peerID := range peerIDs {
			cset := peerToCIDRs[peerID]
			if len(cset) == 0 {
				continue
			}
			remote := idToPeer[peerID]
			if remote == nil {
				continue
			}

			// endpoint：优先使用 spfAdj 的直连边的 toEndpoint（与实际更一致）
			var specifiedEndpoint *models.Endpoint
			if info := dInfo[[2]uint{owner, peerID}]; info != nil && info.toEndpoint != nil {
				specifiedEndpoint = info.toEndpoint
			}

			base, err := remote.AsBasePeerConfig(specifiedEndpoint)
			if err != nil {
				return nil, nil, errors.Join(errors.New("build peer base config failed"), err)
			}

			cidrs := make([]string, 0, len(cset))
			for c := range cset {
				if prevOwner, ok := seen[c]; ok && prevOwner != peerID {
					return nil, nil, fmt.Errorf("duplicate allowed ip on node %d: %s appears in peer %d and peer %d", owner, c, prevOwner, peerID)
				}
				seen[c] = peerID
				cidrs = append(cidrs, c)
			}
			sort.Strings(cidrs)
			base.AllowedIps = lo.Uniq(cidrs)
			pcs = append(pcs, base)

			// 用 fullAdj 补齐展示指标（latency/up/endpoint）
			lat, up, ep, explicit := lookupEdgeForDisplay(fullAdj, owner, peerID)
			edges = append(edges, Edge{
				to:         peerID,
				latency:    lat,
				upMbps:     up,
				toEndpoint: ep,
				routes:     base.AllowedIps,
				explicit:   explicit,
			})
		}

		// 按 client_id 稳定排序（保持原接口习惯）
		sort.SliceStable(pcs, func(i, j int) bool { return pcs[i].GetClientId() < pcs[j].GetClientId() })
		sort.SliceStable(edges, func(i, j int) bool { return edges[i].to < edges[j].to })

		result[owner] = pcs
		finalEdges[owner] = edges
	}

	return result, finalEdges, nil
}

func ensureAllowedSet(m map[uint]map[uint]map[string]struct{}, owner, peer uint) map[string]struct{} {
	if _, ok := m[owner]; !ok {
		m[owner] = make(map[uint]map[string]struct{}, 8)
	}
	if _, ok := m[owner][peer]; !ok {
		m[owner][peer] = make(map[string]struct{}, 32)
	}
	return m[owner][peer]
}

func lookupEdgeForDisplay(fullAdj map[uint][]Edge, from, to uint) (latency uint32, up uint32, ep *models.Endpoint, explicit bool) {
	edges := fullAdj[from]
	for _, e := range edges {
		if e.to == to {
			return e.latency, e.upMbps, e.toEndpoint, e.explicit
		}
	}
	return 0, 0, nil, false
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

// findNextHop 返回从 src 到 dst 的 nextHop（src 的直连邻居），依赖 prev[dst] = predecessor(dst)
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

// 仅用于测试/诊断：解析 /32 的 host ip（校验格式）
func parseHostFromCIDR(c string) (netip.Addr, bool) {
	p, err := netip.ParsePrefix(c)
	if err != nil {
		return netip.Addr{}, false
	}
	return p.Addr(), true
}

// getHandshakeAgeBetween 返回 a<->b 间 peer handshake 的“最大”年龄（只要任意一侧可观测到握手时间就生效）。
// 选择 max 的原因：如果任一方向握手过旧，都应抑制这对节点作为可靠转发 hop。
func getHandshakeAgeBetween(aWGID, bWGID uint, idToPeer map[uint]*models.WireGuard, policy RoutingPolicy) (time.Duration, bool) {
	ageA, okA := getOneWayHandshakeAge(aWGID, bWGID, idToPeer, policy)
	ageB, okB := getOneWayHandshakeAge(bWGID, aWGID, idToPeer, policy)
	if !okA && !okB {
		return 0, false
	}
	if !okA {
		return ageB, true
	}
	if !okB {
		return ageA, true
	}
	if ageA >= ageB {
		return ageA, true
	}
	return ageB, true
}

// getOneWayHandshakeAge 从 fromWGID 的 runtimeInfo 中查到 toWGID 对应 peer 的 last_handshake_time_sec/nsec，返回握手“距离现在”的时间差。
func getOneWayHandshakeAge(fromWGID, toWGID uint, idToPeer map[uint]*models.WireGuard, policy RoutingPolicy) (time.Duration, bool) {
	if policy.NetworkTopologyCache == nil {
		return 0, false
	}
	toPeer := idToPeer[toWGID]
	if toPeer == nil || toPeer.ClientID == "" {
		return 0, false
	}
	runtimeInfo, ok := policy.NetworkTopologyCache.GetRuntimeInfo(fromWGID)
	if !ok || runtimeInfo == nil {
		return 0, false
	}
	var hsSec uint64
	var hsNsec uint64
	for _, p := range runtimeInfo.GetPeers() {
		if p == nil {
			continue
		}
		if p.GetClientId() != toPeer.ClientID {
			continue
		}
		hsSec = p.GetLastHandshakeTimeSec()
		hsNsec = p.GetLastHandshakeTimeNsec()
		break
	}
	if hsSec == 0 {
		return 0, false
	}
	t := time.Unix(int64(hsSec), int64(hsNsec))
	age := time.Since(t)
	if age < 0 {
		age = 0
	}
	return age, true
}
