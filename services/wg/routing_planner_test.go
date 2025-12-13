package wg

import (
	"testing"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
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

func TestRunAllPairsDijkstra_PreferFreshHandshake(t *testing.T) {
	// 1 -> 2 (stale handshake)
	// 1 -> 3 (fresh)
	// 3 -> 2 (fresh)
	// 期望：从 1 到 2 的 nextHop 选择 3，而不是 2
	now := time.Now().Unix()

	priv1, _ := wgtypes.GeneratePrivateKey()
	priv2, _ := wgtypes.GeneratePrivateKey()
	priv3, _ := wgtypes.GeneratePrivateKey()

	p1 := &models.WireGuard{WireGuardEntity: &models.WireGuardEntity{
		ClientID:     "c1",
		PrivateKey:   priv1.String(),
		LocalAddress: "10.0.0.1/32",
	}}
	p1.ID = 1
	p2 := &models.WireGuard{WireGuardEntity: &models.WireGuardEntity{
		ClientID:     "c2",
		PrivateKey:   priv2.String(),
		LocalAddress: "10.0.0.2/32",
	}}
	p2.ID = 2
	p3 := &models.WireGuard{WireGuardEntity: &models.WireGuardEntity{
		ClientID:     "c3",
		PrivateKey:   priv3.String(),
		LocalAddress: "10.0.0.3/32",
	}}
	p3.ID = 3

	idToPeer := map[uint]*models.WireGuard{1: p1, 2: p2, 3: p3}
	order := []uint{1, 2, 3}
	adj := map[uint][]Edge{
		1: {
			{to: 2, latency: 10, upMbps: 50, explicit: true},
			{to: 3, latency: 5, upMbps: 50, explicit: true},
		},
		3: {
			{to: 2, latency: 5, upMbps: 50, explicit: true},
		},
		2: {},
	}

	cache := &fakeTopologyCache{
		rt: map[uint]*pb.WGDeviceRuntimeInfo{
			1: {
				Peers: []*pb.WGPeerRuntimeInfo{
					{ClientId: "c2", LastHandshakeTimeSec: uint64(now - 3600)}, // stale
					{ClientId: "c3", LastHandshakeTimeSec: uint64(now)},        // fresh
				},
			},
			3: {
				Peers: []*pb.WGPeerRuntimeInfo{
					{ClientId: "c2", LastHandshakeTimeSec: uint64(now)}, // fresh
				},
			},
		},
	}

	policy := RoutingPolicy{
		LatencyWeight:           1,
		InverseBandwidthWeight:  0,
		HopWeight:               0,
		HandshakeStaleThreshold: 1 * time.Second,
		HandshakeStalePenalty:   100,
		NetworkTopologyCache:    cache,
	}

	aggByNode, _ := runAllPairsDijkstra(order, adj, idToPeer, policy)
	if aggByNode[1] == nil {
		t.Fatalf("aggByNode[1] should not be nil")
	}
	// dst=2 的 CIDR 应该被聚合到 nextHop=3 下（而不是 nextHop=2）
	if _, ok := aggByNode[1][3]; !ok {
		t.Fatalf("want nextHop=3 for src=1, got keys=%v", keysUint(aggByNode[1]))
	}
	if _, ok := aggByNode[1][2]; ok {
		t.Fatalf("did not expect nextHop=2 for src=1 when handshake is stale")
	}
}

func keysUint(m map[uint]map[string]struct{}) []uint {
	ret := make([]uint, 0, len(m))
	for k := range m {
		ret = append(ret, k)
	}
	return ret
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
	// 构造一个“1 直连 2，但 2 到 1 会更偏好走 3”的场景：
	// 1->2 成为承载路由的 nextHop，但 2 的路由结果中可能不包含 peer(1)，需要对称补齐。
	now := time.Now().Unix()

	priv1, _ := wgtypes.GeneratePrivateKey()
	priv2, _ := wgtypes.GeneratePrivateKey()
	priv3, _ := wgtypes.GeneratePrivateKey()

	p1 := &models.WireGuard{WireGuardEntity: &models.WireGuardEntity{
		ClientID:     "c1",
		PrivateKey:   priv1.String(),
		LocalAddress: "10.0.0.1/32",
	}}
	p1.ID = 1
	p2 := &models.WireGuard{WireGuardEntity: &models.WireGuardEntity{
		ClientID:     "c2",
		PrivateKey:   priv2.String(),
		LocalAddress: "10.0.0.2/32",
	}}
	p2.ID = 2
	p3 := &models.WireGuard{WireGuardEntity: &models.WireGuardEntity{
		ClientID:     "c3",
		PrivateKey:   priv3.String(),
		LocalAddress: "10.0.0.3/32",
	}}
	p3.ID = 3

	idToPeer := map[uint]*models.WireGuard{1: p1, 2: p2, 3: p3}
	order := []uint{1, 2, 3}

	// 全双向连通，但设置权重让 2->1 更偏好 2->3->1
	adj := map[uint][]Edge{
		1: {
			{to: 2, latency: 1, upMbps: 50, explicit: true},
			{to: 3, latency: 100, upMbps: 50, explicit: true},
		},
		2: {
			{to: 1, latency: 100, upMbps: 50, explicit: true},
			{to: 3, latency: 1, upMbps: 50, explicit: true},
		},
		3: {
			{to: 1, latency: 1, upMbps: 50, explicit: true},
			{to: 2, latency: 100, upMbps: 50, explicit: true},
		},
	}

	cache := &fakeTopologyCache{
		rt: map[uint]*pb.WGDeviceRuntimeInfo{
			1: {Peers: []*pb.WGPeerRuntimeInfo{{ClientId: "c2", LastHandshakeTimeSec: uint64(now)}}},
			2: {Peers: []*pb.WGPeerRuntimeInfo{{ClientId: "c1", LastHandshakeTimeSec: uint64(now)}}},
		},
	}

	policy := RoutingPolicy{
		LatencyWeight:           1,
		InverseBandwidthWeight:  0,
		HopWeight:               0,
		HandshakeStaleThreshold: 1 * time.Hour,
		HandshakeStalePenalty:   0,
		NetworkTopologyCache:    cache,
	}

	aggByNode, edgeInfo := runAllPairsDijkstra(order, adj, idToPeer, policy)
	peersMap, err := assemblePeerConfigs(order, aggByNode, edgeInfo, idToPeer)
	if err != nil {
		t.Fatalf("assemblePeerConfigs err: %v", err)
	}
	fillIsolates(order, peersMap)

	// 预期：在没有对称补齐前，2 可能不会包含 peer(1)
	_ = ensureRoutingPeerSymmetry(order, peersMap, idToPeer)

	found := false
	for _, pc := range peersMap[2] {
		if pc != nil && pc.GetId() == 1 {
			found = true
			if len(pc.GetAllowedIps()) == 0 || pc.GetAllowedIps()[0] != "10.0.0.1/32" {
				t.Fatalf("peer(1) on node2 should include 10.0.0.1/32, got=%v", pc.GetAllowedIps())
			}
		}
	}
	if !found {
		t.Fatalf("node2 should contain peer(1) after ensureRoutingPeerSymmetry")
	}
}
