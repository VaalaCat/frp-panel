package wg

import (
	"testing"

	"github.com/VaalaCat/frp-panel/pb"
)

type fakeTopologyCache struct {
	lat map[[2]uint]uint32
}

func (c *fakeTopologyCache) GetRuntimeInfo(_ uint) (*pb.WGDeviceRuntimeInfo, bool) { return nil, false }
func (c *fakeTopologyCache) SetRuntimeInfo(_ uint, _ *pb.WGDeviceRuntimeInfo)      {}
func (c *fakeTopologyCache) DeleteRuntimeInfo(_ uint)                              {}
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
