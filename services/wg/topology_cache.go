package wg

import (
	"fmt"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils"
)

var (
	_ app.NetworkTopologyCache = (*networkTopologyCache)(nil)
)

// networkTopologyCache 目前只给服务端用
type networkTopologyCache struct {
	wireguardRuntimeInfoMap *utils.SyncMap[uint, *pb.WGDeviceRuntimeInfo] // wireguardId -> peerRuntimeInfo
	fromToLatencyMap        *utils.SyncMap[string, uint32]                // fromWGID::toWGID -> latencyMs
}

func NewNetworkTopologyCache() *networkTopologyCache {
	return &networkTopologyCache{
		wireguardRuntimeInfoMap: &utils.SyncMap[uint, *pb.WGDeviceRuntimeInfo]{},
		fromToLatencyMap:        &utils.SyncMap[string, uint32]{},
	}
}

func (c *networkTopologyCache) GetRuntimeInfo(wireguardId uint) (*pb.WGDeviceRuntimeInfo, bool) {
	return c.wireguardRuntimeInfoMap.Load(wireguardId)
}

func (c *networkTopologyCache) SetRuntimeInfo(wireguardId uint, runtimeInfo *pb.WGDeviceRuntimeInfo) {
	c.wireguardRuntimeInfoMap.Store(wireguardId, runtimeInfo)
	for toWireGuardId, latencyMs := range runtimeInfo.GetPingMap() {
		c.fromToLatencyMap.Store(parseFromToLatencyKey(wireguardId, uint(toWireGuardId)), latencyMs)
	}
}

func (c *networkTopologyCache) DeleteRuntimeInfo(wireguardId uint) {
	c.wireguardRuntimeInfoMap.Delete(wireguardId)
}

func (c *networkTopologyCache) GetLatencyMs(fromWGID, toWGID uint) (uint32, bool) {
	v1, ok := c.fromToLatencyMap.Load(parseFromToLatencyKey(fromWGID, toWGID))
	if !ok {
		return c.fromToLatencyMap.Load(parseFromToLatencyKey(toWGID, fromWGID))
	}
	return v1, true
}

func parseFromToLatencyKey(fromWGID, toWGID uint) string {
	return fmt.Sprintf("%d::%d", fromWGID, toWGID)
}
