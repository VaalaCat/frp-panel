package wg

import (
	"testing"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/samber/lo"
)

type fakeTopologyCache struct {
	lat map[[2]uint]uint32
	rt  map[uint]*pb.WGDeviceRuntimeInfo
}

func (c *fakeTopologyCache) GetRuntimeInfo(id uint) (*pb.WGDeviceRuntimeInfo, bool) {
	if c == nil || c.rt == nil {
		return nil, false
	}
	v, ok := c.rt[id]
	return v, ok
}
func (c *fakeTopologyCache) SetRuntimeInfo(_ uint, _ *pb.WGDeviceRuntimeInfo) {}
func (c *fakeTopologyCache) DeleteRuntimeInfo(_ uint)                         {}
func (c *fakeTopologyCache) GetLatencyMs(fromWGID, toWGID uint) (uint32, bool) {
	if c == nil || c.lat == nil {
		return 0, false
	}
	v, ok := c.lat[[2]uint{fromWGID, toWGID}]
	return v, ok
}

func TestFilterAdjacencyForSPF(t *testing.T) {
	cache := &fakeTopologyCache{
		lat: map[[2]uint]uint32{
			{1, 2}: 10,
			{2, 1}: 10,
			{1, 4}: ^uint32(0), // MaxUint32: 明确不可达，应被剔除
		},
	}

	policy := RoutingPolicy{NetworkTopologyCache: cache}

	adj := map[uint][]Edge{
		1: {
			{to: 2, latency: 30, explicit: false},              // implicit + 有探测 => 保留，且 latency 覆盖为 10
			{to: 3, latency: 30, explicit: false},              // implicit + 无探测 => 剔除
			{to: 4, latency: 30, explicit: false},              // implicit + 不可达哨兵 => 剔除
			{to: 5, latency: 1, upMbps: 1, explicit: true},     // explicit => 保留
			{to: 6, latency: 999, upMbps: 1, explicit: true},   // explicit => 保留
			{to: 7, latency: 30, upMbps: 50, explicit: false},  // implicit + 无探测 => 剔除
			{to: 8, latency: 30, upMbps: 50, explicit: false},  // implicit + 无探测 => 剔除
			{to: 9, latency: 30, upMbps: 50, explicit: false},  // implicit + 无探测 => 剔除
			{to: 10, latency: 30, upMbps: 50, explicit: false}, // implicit + 无探测 => 剔除
		},
	}

	ret := filterAdjacencyForSPF([]uint{1, 2}, adj, policy)

	edges1 := ret[1]
	if len(edges1) != 3 {
		t.Fatalf("want 3 edges for node 1, got %d: %#v", len(edges1), edges1)
	}

	// 校验 implicit edge(1->2) 的 latency 被覆盖为真实探测值 10
	found12 := false
	for _, e := range edges1 {
		if e.to == 2 {
			found12 = true
			if e.latency != 10 {
				t.Fatalf("want latency=10 for edge 1->2, got %d", e.latency)
			}
		}
		if e.to == 4 {
			t.Fatalf("edge 1->4 should be filtered out (MaxUint32)")
		}
		if e.to == 3 {
			t.Fatalf("edge 1->3 should be filtered out (no probe data)")
		}
	}
	if !found12 {
		t.Fatalf("edge 1->2 should exist")
	}

	// order 中的节点必须存在 key（即便没有边）
	if _, ok := ret[2]; !ok {
		t.Fatalf("node 2 should exist in return map")
	}
}

func TestPlanAllowedIPs_PreferFreshHandshake(t *testing.T) {
	// 1 <-> 2：低延迟但握手过旧（应被惩罚）
	// 1 <-> 3 <-> 2：略高延迟但握手新（应被选为 1->2 的 nextHop=3）
	now := time.Now().Unix()

	priv1, _ := wgtypes.GeneratePrivateKey()
	priv2, _ := wgtypes.GeneratePrivateKey()
	priv3, _ := wgtypes.GeneratePrivateKey()

	p1 := &models.WireGuard{WireGuardEntity: &models.WireGuardEntity{ClientID: "c1", PrivateKey: priv1.String(), LocalAddress: "10.0.0.1/32"}}
	p1.ID = 1
	p2 := &models.WireGuard{WireGuardEntity: &models.WireGuardEntity{ClientID: "c2", PrivateKey: priv2.String(), LocalAddress: "10.0.0.2/32"}}
	p2.ID = 2
	p3 := &models.WireGuard{WireGuardEntity: &models.WireGuardEntity{ClientID: "c3", PrivateKey: priv3.String(), LocalAddress: "10.0.0.3/32"}}
	p3.ID = 3

	// 显式链路也要求至少一侧存在 endpoint（符合真实运行时：需要可连接入口）
	p1.AdvertisedEndpoints = []*models.Endpoint{{EndpointEntity: &models.EndpointEntity{Host: "redacted.example", Port: 61820, Type: "ws", WireGuardID: 1, ClientID: "c1"}}}
	p2.AdvertisedEndpoints = []*models.Endpoint{{EndpointEntity: &models.EndpointEntity{Host: "redacted.example", Port: 61820, Type: "ws", WireGuardID: 2, ClientID: "c2"}}}
	p3.AdvertisedEndpoints = []*models.Endpoint{{EndpointEntity: &models.EndpointEntity{Host: "redacted.example", Port: 61820, Type: "ws", WireGuardID: 3, ClientID: "c3"}}}

	peers := []*models.WireGuard{p1, p2, p3}
	link := func(from, to uint, latency uint32) *models.WireGuardLink {
		return &models.WireGuardLink{WireGuardLinkEntity: &models.WireGuardLinkEntity{
			FromWireGuardID: from,
			ToWireGuardID:   to,
			UpBandwidthMbps: 50,
			LatencyMs:       latency,
			Active:          true,
		}}
	}
	links := []*models.WireGuardLink{
		link(1, 2, 5), link(2, 1, 5),
		link(1, 3, 8), link(3, 1, 8),
		link(3, 2, 8), link(2, 3, 8),
	}

	cache := &fakeTopologyCache{
		rt: map[uint]*pb.WGDeviceRuntimeInfo{
			1: {Peers: []*pb.WGPeerRuntimeInfo{
				{ClientId: "c2", LastHandshakeTimeSec: uint64(now - 3600)}, // stale
				{ClientId: "c3", LastHandshakeTimeSec: uint64(now)},        // fresh
			}},
			2: {Peers: []*pb.WGPeerRuntimeInfo{
				{ClientId: "c1", LastHandshakeTimeSec: uint64(now - 3600)}, // stale (对称)
				{ClientId: "c3", LastHandshakeTimeSec: uint64(now)},        // fresh
			}},
			3: {Peers: []*pb.WGPeerRuntimeInfo{
				{ClientId: "c1", LastHandshakeTimeSec: uint64(now)},
				{ClientId: "c2", LastHandshakeTimeSec: uint64(now)},
			}},
		},
	}

	policy := DefaultRoutingPolicy(NewACL(), cache, nil)
	policy.HandshakeStaleThreshold = 1 * time.Second
	policy.HandshakeStalePenalty = 1000
	policy.InverseBandwidthWeight = 0
	policy.HopWeight = 0
	policy.LatencyLogScale = 0

	peerCfgs, _, err := PlanAllowedIPs(peers, links, policy)
	if err != nil {
		t.Fatalf("PlanAllowedIPs err: %v", err)
	}

	// 对 node1：10.0.0.2/32 应走 peer(3) 而不是 peer(2)
	wantDst := "10.0.0.2/32"
	var gotPeer uint32
	for _, pc := range peerCfgs[1] {
		if pc == nil {
			continue
		}
		if lo.Contains(pc.GetAllowedIps(), wantDst) {
			gotPeer = pc.GetId()
		}
	}
	if gotPeer != 3 {
		t.Fatalf("want node1 route %s via peer 3, got peer %d", wantDst, gotPeer)
	}
}

func TestSymmetrizeAdjacencyForPeers_FillReverseEdge(t *testing.T) {
	t.Skip("symmetrizeAdjacencyForPeers 已移除：路由承载的边必须双向存在，不应自动补齐单向边")
}

func TestFilterAdjacencyForSymmetricLinks_DropOneWay(t *testing.T) {
	order := []uint{1, 2}
	adj := map[uint][]Edge{
		1: {{to: 2, latency: 10, upMbps: 50, explicit: true}}, // 单向
		2: {},
	}
	ret := filterAdjacencyForSymmetricLinks(order, adj)
	if len(ret[1]) != 0 {
		t.Fatalf("want 0 edges for node 1 after symmetric filter, got %d: %#v", len(ret[1]), ret[1])
	}
	if _, ok := ret[2]; !ok {
		t.Fatalf("node 2 should exist in return map")
	}
}

func TestEnsureRoutingPeerSymmetry_AddReversePeer(t *testing.T) {
	t.Skip("routing planner rewritten: inbound-source-set generation replaces old symmetry patching")
}

func TestPlanAllowedIPs_Regression_NoDuplicateAllowedIPs_And_TransitSourceValidation(t *testing.T) {
	// 复现 & 防回归：
	// 1) 同一节点的 AllowedIPs 不允许在多个 peer 间重复（例如 10.10.0.4/32 只能分配给一个 nextHop）
	// 2) 多跳转发时，入站 source validation 需要允许“原始源地址”：
	//    构造 21(10.10.0.8) -> 16(10.10.0.2) 走 24 中转，
	//    期望 16 的 peer(24) AllowedIPs 包含 10.10.0.8/32（否则 16 会丢弃来自 24 的转发包）。

	type node struct {
		id    uint
		cid   string
		addr  string
		tags  []string
		hasEP bool
	}

	nodes := []node{
		{id: 4, cid: "c4", addr: "10.10.0.4/24", tags: []string{"cn", "bj"}, hasEP: true},
		{id: 11, cid: "c11", addr: "10.10.0.1/24", tags: []string{"cn", "wh"}, hasEP: false},
		{id: 16, cid: "c16", addr: "10.10.0.2/24", tags: []string{"cn", "bj", "ali"}, hasEP: true},
		{id: 17, cid: "c17", addr: "10.10.0.3/24", tags: []string{"cn", "wh"}, hasEP: false},
		{id: 18, cid: "c18", addr: "10.10.0.6/24", tags: []string{"us"}, hasEP: true},
		{id: 20, cid: "c20", addr: "10.10.0.7/24", tags: []string{"us"}, hasEP: false},
		{id: 21, cid: "c21", addr: "10.10.0.8/24", tags: []string{"cn", "nc"}, hasEP: false},
		{id: 22, cid: "c22", addr: "10.10.0.9/24", tags: []string{"cn", "nc"}, hasEP: false},
		{id: 24, cid: "c24", addr: "10.10.0.5/24", tags: []string{"cn", "nc"}, hasEP: true},
	}

	makePeer := func(n node) *models.WireGuard {
		priv, _ := wgtypes.GeneratePrivateKey()
		wg := &models.WireGuard{WireGuardEntity: &models.WireGuardEntity{
			ClientID:     n.cid,
			PrivateKey:   priv.String(),
			LocalAddress: n.addr,
			Tags:         n.tags,
		}}
		wg.ID = n.id
		if n.hasEP {
			wg.AdvertisedEndpoints = []*models.Endpoint{
				{EndpointEntity: &models.EndpointEntity{
					Host:        "redacted.example",
					Port:        61820,
					Type:        "ws",
					WireGuardID: n.id,
					ClientID:    n.cid,
				}},
			}
		}
		return wg
	}

	peers := lo.Map(nodes, func(n node, _ int) *models.WireGuard { return makePeer(n) })

	// 构造 ACL（与用户提供一致：只验证 tag 匹配逻辑正确，不涉及公网信息）
	acl := NewACL().LoadFromPB(&pb.AclConfig{Acls: []*pb.AclRuleConfig{
		{Action: "allow", Src: []string{"bj", "wh"}, Dst: []string{"bj", "wh"}},
		{Action: "allow", Src: []string{"nc", "wh"}, Dst: []string{"nc", "wh"}},
		{Action: "allow", Src: []string{"nc", "ali"}, Dst: []string{"nc", "ali"}},
		{Action: "allow", Src: []string{"wh", "ali"}, Dst: []string{"wh", "ali"}},
		{Action: "allow", Src: []string{"us"}, Dst: []string{"us"}},
	}})

	// 只需要 latency cache 为推断边提供“探测存在性”，这里直接手动构造显式 links，更可控
	// 关键：让 21->16 走 24 中转（21-24-16 低延迟，21-16 高延迟）
	link := func(from, to uint, latency uint32) *models.WireGuardLink {
		return &models.WireGuardLink{WireGuardLinkEntity: &models.WireGuardLinkEntity{
			FromWireGuardID: from,
			ToWireGuardID:   to,
			UpBandwidthMbps: 50,
			LatencyMs:       latency,
			Active:          true,
		}}
	}
	links := []*models.WireGuardLink{
		link(21, 24, 10), link(24, 21, 10),
		link(24, 16, 10), link(16, 24, 10),
		link(21, 16, 200), link(16, 21, 200),

		// 再补一些连通边，确保能算出包含 4 的路由
		link(11, 16, 30), link(16, 11, 30),
		link(16, 4, 5), link(4, 16, 5),
		link(11, 4, 50), link(4, 11, 50),
	}

	policy := DefaultRoutingPolicy(acl, &fakeTopologyCache{lat: map[[2]uint]uint32{}}, nil)
	policy.HandshakeStalePenalty = 0
	policy.HandshakeStaleThreshold = 0
	policy.InverseBandwidthWeight = 0
	policy.HopWeight = 0
	policy.LatencyLogScale = 0

	peerCfgs, _, err := PlanAllowedIPs(peers, links, policy)
	if err != nil {
		t.Fatalf("PlanAllowedIPs err: %v", err)
	}

	// 1) 断言：每个节点的 AllowedIPs 在不同 peer 间不重复
	for owner, pcs := range peerCfgs {
		seen := map[string]uint32{}
		for _, pc := range pcs {
			if pc == nil {
				continue
			}
			for _, cidr := range pc.GetAllowedIps() {
				if prev, ok := seen[cidr]; ok && prev != pc.GetId() {
					t.Fatalf("node %d has duplicate cidr %s on peer %d and peer %d", owner, cidr, prev, pc.GetId())
				}
				seen[cidr] = pc.GetId()
			}
		}
	}

	// 2) 断言：16 的 peer(24) 必须包含 10.10.0.8/32（21 的 /32），用于入站 source validation
	wantSrc := "10.10.0.8/32"
	found := false
	for _, pc := range peerCfgs[16] {
		if pc == nil || pc.GetId() != 24 {
			continue
		}
		if lo.Contains(pc.GetAllowedIps(), wantSrc) {
			found = true
		}
	}
	if !found {
		t.Fatalf("node 16 peer(24) should contain %s for transit source validation", wantSrc)
	}

	// 3) 断言：11 节点的 10.10.0.4/32 不能同时出现在多个 peer
	wantC4 := "10.10.0.4/32"
	var peersWithC4 []uint32
	for _, pc := range peerCfgs[11] {
		if pc == nil {
			continue
		}
		if lo.Contains(pc.GetAllowedIps(), wantC4) {
			peersWithC4 = append(peersWithC4, pc.GetId())
		}
	}
	if len(peersWithC4) != 1 {
		t.Fatalf("node 11 should have exactly one peer carrying %s, got peers=%v", wantC4, peersWithC4)
	}
}

func TestBuildAdjacency_InferredEdgesAreBidirectionalWhenACLAllows(t *testing.T) {
	// 回归：推断边必须支持 to(with endpoint) -> from(no endpoint) 的反向补齐，
	// 否则 filterAdjacencyForSymmetricLinks 会把所有 “no-endpoint 节点” 剔除，导致 SPF 结果为空。

	privA, _ := wgtypes.GeneratePrivateKey()
	privB, _ := wgtypes.GeneratePrivateKey()

	a := &models.WireGuard{WireGuardEntity: &models.WireGuardEntity{
		ClientID:     "ca",
		PrivateKey:   privA.String(),
		LocalAddress: "10.0.0.1/24",
		Tags:         []string{"t1"},
	}}
	a.ID = 1 // no endpoint

	b := &models.WireGuard{WireGuardEntity: &models.WireGuardEntity{
		ClientID:     "cb",
		PrivateKey:   privB.String(),
		LocalAddress: "10.0.0.2/24",
		Tags:         []string{"t1"},
	}}
	b.ID = 2
	b.AdvertisedEndpoints = []*models.Endpoint{{EndpointEntity: &models.EndpointEntity{
		Host:        "redacted.example",
		Port:        61820,
		Type:        "ws",
		WireGuardID: 2,
		ClientID:    "cb",
	}}}

	idToPeer, order := buildNodeIndexSorted([]*models.WireGuard{a, b})
	acl := NewACL().LoadFromPB(&pb.AclConfig{Acls: []*pb.AclRuleConfig{
		{Action: "allow", Src: []string{"t1"}, Dst: []string{"t1"}},
	}})
	policy := DefaultRoutingPolicy(acl, &fakeTopologyCache{lat: map[[2]uint]uint32{
		{1, 2}: 10,
		{2, 1}: 10,
	}}, nil)

	adj := buildAdjacency(order, idToPeer, nil, policy)
	// 期望：1->2 与 2->1 都存在（推断边双向）
	has12 := false
	for _, e := range adj[1] {
		if e.to == 2 {
			has12 = true
		}
	}
	has21 := false
	for _, e := range adj[2] {
		if e.to == 1 {
			has21 = true
		}
	}
	if !has12 || !has21 {
		t.Fatalf("want inferred edges 1->2 and 2->1, got has12=%v has21=%v adj=%#v", has12, has21, adj)
	}
}
