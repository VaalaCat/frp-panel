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
	virtAddrPingMap         *utils.SyncMap[string, uint32]                // fromWGID::toWGID -> pingMs
}

func NewNetworkTopologyCache() *networkTopologyCache {
	return &networkTopologyCache{
		wireguardRuntimeInfoMap: &utils.SyncMap[uint, *pb.WGDeviceRuntimeInfo]{},
		fromToLatencyMap:        &utils.SyncMap[string, uint32]{},
		virtAddrPingMap:         &utils.SyncMap[string, uint32]{},
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

	for virtAddr, pingMs := range runtimeInfo.GetVirtAddrPingMap() {
		c.virtAddrPingMap.Store(parseFromToLatencyKey(wireguardId, uint(runtimeInfo.PeerVirtAddrMap[virtAddr])), pingMs)
	}
}

func (c *networkTopologyCache) DeleteRuntimeInfo(wireguardId uint) {
	c.wireguardRuntimeInfoMap.Delete(wireguardId)
}

func (c *networkTopologyCache) GetLatencyMs(fromWGID, toWGID uint) (uint32, bool) {

	endpointLatency, ok := c.fromToLatencyMap.Load(parseFromToLatencyKey(fromWGID, toWGID))
	if !ok {
		endpointLatency, ok = c.fromToLatencyMap.Load(parseFromToLatencyKey(toWGID, fromWGID))
		if !ok {
			return 0, false
		}
	}

	virtAddrLatency, ok := c.virtAddrPingMap.Load(parseFromToLatencyKey(fromWGID, toWGID))
	if !ok {
		virtAddrLatency, ok = c.virtAddrPingMap.Load(parseFromToLatencyKey(toWGID, fromWGID))
		if !ok {
			return endpointLatency, false
		}
	}

	return (endpointLatency + virtAddrLatency) / 2, true
}

func parseFromToLatencyKey(fromWGID, toWGID uint) string {
	return fmt.Sprintf("%d::%d", fromWGID, toWGID)
}
