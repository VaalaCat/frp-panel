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
// cost = LatencyWeight*latency_ms + InverseBandwidthWeight*(1/max(up_mbps,1e-6)) + HopWeight + HandshakePenalty
type RoutingPolicy struct {
	LatencyWeight            float64
	InverseBandwidthWeight   float64
	HopWeight                float64
	MinUpMbps                uint32
	DefaultEndpointUpMbps    uint32
	DefaultEndpointLatencyMs uint32
	OfflineThreshold         time.Duration
	// HandshakeStaleThreshold/HandshakeStalePenalty 用于抑制“握手过旧”的链路被选为最短路。
	// 仅在能从 runtimeInfo 中找到对应 peer 的 last_handshake_time_sec 时生效；否则不惩罚（避免误伤）。
	HandshakeStaleThreshold time.Duration
	HandshakeStalePenalty   float64

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
		// 默认启用一个温和的“握手过旧惩罚”：优先选择近期有握手的链路，但不至于强制剔除路径。
		HandshakeStaleThreshold: 5 * time.Minute,
		HandshakeStalePenalty:   30.0,
		ACL:                     acl,
		NetworkTopologyCache:    networkTopologyCache,
		CliMgr:                  cliMgr,
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
	spfAdj := filterAdjacencyForSPF(order, adj, p.policy)
	// 路由（AllowedIPs）依赖 WireGuard 的“源地址校验”：下一跳收到的包会按“来自哪个 peer”做匹配，
	// 并校验 inner packet 的 source IP 是否落在该 peer 的 AllowedIPs 中。
	// 因此用于承载路由的直连边必须是双向的：若存在单向边，最短路会产生单向选路，导致中间节点丢包。
	spfAdj = filterAdjacencyForSymmetricLinks(order, spfAdj)
	aggByNode, edgeInfoMap := runAllPairsDijkstra(order, spfAdj, idToPeer, p.policy)
	result, err := assemblePeerConfigs(order, aggByNode, edgeInfoMap, idToPeer)
	if err != nil {
		return nil, nil, err
	}
	fillIsolates(order, result)
	if err := ensureRoutingPeerSymmetry(order, result, idToPeer); err != nil {
		return nil, nil, err
	}

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
	spfAdj := filterAdjacencyForSPF(order, adj, p.policy)
	spfAdj = filterAdjacencyForSymmetricLinks(order, spfAdj)
	routesInfoMap, edgeInfoMap := runAllPairsDijkstra(order, spfAdj, idToPeer, p.policy)

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

// filterAdjacencyForSymmetricLinks 仅保留“存在反向直连边”的邻接（用于 SPF）。
// 这样最短路产生的每一步转发 hop 都对应一个双向直连 peer，避免出现单向路由导致的丢包。
func filterAdjacencyForSymmetricLinks(order []uint, adj map[uint][]Edge) map[uint][]Edge {
	ret := make(map[uint][]Edge, len(order))
	edgeSet := make(map[[2]uint]struct{}, 16)

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

// ensureRoutingPeerSymmetry 确保：如果 src 的 peers 中存在 nextHop（承载路由），则 nextHop 的 peers 中也必须存在 src。
// 这里“对称”不是指两端 routes/AllowedIPs 集合一致，而是指两端都必须配置对方这个 peer，
// 以满足 WG 的解密与源地址校验（否则 nextHop 会丢弃来自 src 的转发包）。
func ensureRoutingPeerSymmetry(order []uint, peerCfgs map[uint][]*pb.WireGuardPeerConfig, idToPeer map[uint]*models.WireGuard) error {
	if len(order) == 0 {
		return nil
	}

	// 预计算每个节点自身的 /32 CIDR（AsBasePeerConfig 返回的 AllowedIps[0]）
	selfCIDR := make(map[uint]string, len(order))
	for _, id := range order {
		p := idToPeer[id]
		if p == nil {
			continue
		}
		base, err := p.AsBasePeerConfig(nil)
		if err != nil || len(base.GetAllowedIps()) == 0 {
			continue
		}
		selfCIDR[id] = base.GetAllowedIps()[0]
	}

	hasPeer := func(owner uint, peerID uint) bool {
		for _, pc := range peerCfgs[owner] {
			if pc == nil {
				continue
			}
			if uint(pc.GetId()) == peerID {
				return true
			}
		}
		return false
	}

	for _, src := range order {
		for _, pc := range peerCfgs[src] {
			if pc == nil {
				continue
			}
			if len(pc.GetAllowedIps()) == 0 {
				continue
			}
			nextHop := uint(pc.GetId())
			if nextHop == 0 || nextHop == src {
				continue
			}
			if hasPeer(nextHop, src) {
				continue
			}

			remote := idToPeer[src]
			if remote == nil {
				continue
			}
			base, err := remote.AsBasePeerConfig(nil)
			if err != nil {
				return err
			}
			if cidr := selfCIDR[src]; cidr != "" {
				base.AllowedIps = []string{cidr}
			}
			peerCfgs[nextHop] = append(peerCfgs[nextHop], base)
		}
	}

	for _, id := range order {
		sort.SliceStable(peerCfgs[id], func(i, j int) bool {
			return peerCfgs[id][i].GetClientId() < peerCfgs[id][j].GetClientId()
		})
	}
	return nil
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

		adj[from] = append(adj[from], Edge{
			to:         to,
			latency:    latency,
			upMbps:     l.UpBandwidthMbps,
			toEndpoint: l.ToEndpoint,
			explicit:   true,
		})
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

				adj[from] = append(adj[from], Edge{
					to:       to,
					latency:  latency,
					upMbps:   policy.DefaultEndpointUpMbps,
					explicit: false,
				})
				edgeSet[key1] = struct{}{}
			}

			if policy.ACL.CanConnect(idToPeer[to], idToPeer[from]) {
				key2 := [2]uint{to, from}
				if _, exists := edgeSet[key2]; exists {
					continue
				}
				adj[to] = append(adj[to], Edge{
					to:       from,
					latency:  latency,
					upMbps:   policy.DefaultEndpointUpMbps,
					explicit: false,
				})
				edgeSet[key2] = struct{}{}
			}
		}
	}
	return adj
}

// filterAdjacencyForSPF 将“用于探测的候选邻接(adj)”过滤为“允许进入 SPF 的邻接”。
//
// 参考 OSPF：新邻接必须先被确认可达（这里用 runtime ping/virt ping 的存在性作为信号）后，
// 才能参与最短路计算。否则在节点刚更新/刚加入时，会因为默认权重过低被误选，导致部分节点不可达。
func filterAdjacencyForSPF(order []uint, adj map[uint][]Edge, policy RoutingPolicy) map[uint][]Edge {
	ret := make(map[uint][]Edge, len(order))

	for from, edges := range adj {
		for _, e := range edges {
			// 显式 link：管理员配置的边，允许进入 SPF
			if e.explicit {
				ret[from] = append(ret[from], e)
				continue
			}

			// 推断/探测用 link：必须已存在探测数据，且不可达哨兵值要剔除
			latency, ok := policy.NetworkTopologyCache.GetLatencyMs(from, e.to)
			if !ok {
				continue
			}
			if latency == math.MaxUint32 {
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
				handshakePenalty := 0.0
				if policy.HandshakeStalePenalty > 0 && policy.HandshakeStaleThreshold > 0 {
					// 握手惩罚必须是“无方向”的，否则会导致 A->B 与 B->A 权重不一致，
					// 进而产生单向选路（WireGuard AllowedIPs 源地址校验下会丢包）。
					if age, ok := getHandshakeAgeBetween(u, e.to, idToPeer, policy); ok && age > policy.HandshakeStaleThreshold {
						handshakePenalty = policy.HandshakeStalePenalty
					}
				}
				w := policy.LatencyWeight*float64(e.latency) + policy.InverseBandwidthWeight*invBw + policy.HopWeight + handshakePenalty
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

// getOneWayHandshakeAge 从 fromWGID 的 runtimeInfo 中，查找到 toWGID 对应 peer 的 last_handshake_time_sec/nsec，返回握手“距离现在”的时间差。
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
