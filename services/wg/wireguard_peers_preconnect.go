//go:build !windows
// +build !windows

package wg

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/pb"
)

// peer 预连接/常驻补齐逻辑
//
// 目标：在拓扑(adjs)变化时，确保本节点对“可直连/可连接”的 peer 已配置到 wg 设备，
// 但 AllowedIPs 为空（只保持连接，不承载路由），从而避免路由变化导致频繁 remove+add 造成断链。

func (w *wireGuard) onPeersChangedLocked(reason string) {
	if w == nil {
		return
	}
	if err := w.ensureConnectablePeersLocked(); err != nil {
		w.svcLogger.WithError(err).WithField("op", "onPeersChanged").Warnf("ensure connectable peers failed (%s)", reason)
	}
}

func (w *wireGuard) cleanupPreconnectPeersLocked() {
	if w == nil {
		return
	}
	if len(w.preconnectPeers) == 0 {
		return
	}
	exists := make(map[uint32]struct{}, len(w.ifce.Peers))
	for _, p := range w.ifce.GetParsedPeers() {
		if p == nil {
			continue
		}
		if id := p.GetId(); id != 0 {
			exists[id] = struct{}{}
		}
		if p.GetEndpoint() != nil {
			if id := p.GetEndpoint().GetWireguardId(); id != 0 {
				exists[id] = struct{}{}
			}
		}
	}
	for id := range w.preconnectPeers {
		if _, ok := exists[id]; !ok {
			delete(w.preconnectPeers, id)
		}
	}
}

func (w *wireGuard) indexPeerDirectoryLocked(p *defs.WireGuardPeerConfig) {
	if p == nil || p.WireGuardPeerConfig == nil {
		return
	}
	// 1) peer.id
	if id := p.GetId(); id != 0 {
		w.peerDirectory[id] = proto.Clone(p.WireGuardPeerConfig).(*pb.WireGuardPeerConfig)
	}
	// 2) endpoint.wireguard_id（部分场景 peer.id 可能未填，但 endpoint 带 wireguard_id）
	if p.GetEndpoint() != nil {
		if id := p.GetEndpoint().GetWireguardId(); id != 0 {
			w.peerDirectory[id] = proto.Clone(p.WireGuardPeerConfig).(*pb.WireGuardPeerConfig)
		}
	}
}

func (w *wireGuard) deletePeerDirectoryLocked(p *defs.WireGuardPeerConfig) {
	if p == nil {
		return
	}
	if id := p.GetId(); id != 0 {
		delete(w.peerDirectory, id)
	}
	if p.GetEndpoint() != nil {
		if id := p.GetEndpoint().GetWireguardId(); id != 0 {
			delete(w.peerDirectory, id)
		}
	}
}

// ensureConnectablePeersLocked 保证本节点在 wg 设备里已配置所有“当前可直连/可连接”的 peer，
// 但这些补齐 peer 的 AllowedIPs 为空（只保持连接，不承载路由）。
//
// 约束：必须在持有 w.Lock() 的情况下调用。
func (w *wireGuard) ensureConnectablePeersLocked() error {
	if w == nil || w.ifce == nil || w.wgDevice == nil {
		return nil
	}
	localID := w.ifce.GetId()
	if localID == 0 {
		return nil
	}
	adjs := w.ifce.GetAdjs()
	if adjs == nil {
		return nil
	}
	localLinks, ok := adjs[localID]
	if !ok || localLinks == nil || len(localLinks.GetLinks()) == 0 {
		return nil
	}

	// 当前可直连/可连接的 peer id 集合（来自 adj）
	connectable := make(map[uint32]struct{}, len(localLinks.GetLinks()))
	for _, l := range localLinks.GetLinks() {
		if l == nil {
			continue
		}
		toID := l.GetToWireguardId()
		if toID == 0 || toID == localID {
			continue
		}
		connectable[toID] = struct{}{}
	}

	log := w.svcLogger.WithField("op", "ensureConnectablePeers")
	log.Debugf("ensure connectable peers: local=%d connectable=%d peers=%d preconnect=%d directory=%d",
		localID, len(connectable), len(w.ifce.Peers), len(w.preconnectPeers), len(w.peerDirectory))

	// 当前已配置的 peer：用 peer.id 与 endpoint.wireguard_id 双索引，避免 peer.id 缺失导致重复补齐
	exists := make(map[uint32]struct{}, len(w.ifce.Peers))
	for _, p := range w.ifce.GetParsedPeers() {
		if p == nil {
			continue
		}
		if id := p.GetId(); id != 0 {
			exists[id] = struct{}{}
		}
		if p.GetEndpoint() != nil {
			if id := p.GetEndpoint().GetWireguardId(); id != 0 {
				exists[id] = struct{}{}
			}
		}
	}

	uapiBuilder := NewUAPIBuilder()
	added := 0
	removed := 0
	skippedNoBase := 0
	skippedAlready := 0

	// 先清理：相比上次，本次拓扑中已“完全不可直连”的 peer，需要彻底从设备移除
	// 仅清理“AllowedIPs 为空”的 peer（也就是不承载路由、只为保持连接而存在的 peer）
	newPeers := make([]*pb.WireGuardPeerConfig, 0, len(w.ifce.Peers))
	for _, raw := range w.ifce.GetParsedPeers() {
		if raw == nil || raw.WireGuardPeerConfig == nil {
			continue
		}
		// 仅对 AllowedIPs 为空的 peer 做自动清理
		if len(raw.GetAllowedIps()) != 0 {
			newPeers = append(newPeers, raw.WireGuardPeerConfig)
			continue
		}

		var peerID uint32
		if raw.GetId() != 0 {
			peerID = raw.GetId()
		} else if raw.GetEndpoint() != nil && raw.GetEndpoint().GetWireguardId() != 0 {
			peerID = raw.GetEndpoint().GetWireguardId()
		}

		// 无法识别 peer id：保守起见不清理
		if peerID == 0 {
			newPeers = append(newPeers, raw.WireGuardPeerConfig)
			continue
		}

		// 当前拓扑不可直连：彻底移除
		if _, ok := connectable[peerID]; !ok {
			log.Debugf("preconnect remove: peerID=%d pk=%s (reason=not_connectable)", peerID, truncate(raw.GetPublicKey(), 10))
			uapiBuilder.RemovePeerByKey(raw.GetParsedPublicKey())
			delete(w.preconnectPeers, peerID)
			delete(exists, peerID)
			removed++
			continue
		}

		newPeers = append(newPeers, raw.WireGuardPeerConfig)
	}
	// 如果有清理发生，先更新本地缓存（设备更新在最后统一 IpcSet）
	if removed > 0 {
		w.ifce.Peers = newPeers
	}

	for _, l := range localLinks.GetLinks() {
		if l == nil {
			continue
		}
		toID := l.GetToWireguardId()
		if toID == 0 || toID == localID {
			continue
		}
		if _, ok := exists[toID]; ok {
			skippedAlready++
			continue
		}

		base, ok := w.peerDirectory[toID]
		if !ok || base == nil || base.GetPublicKey() == "" {
			skippedNoBase++
			continue
		}
		cloned := &defs.WireGuardPeerConfig{WireGuardPeerConfig: proto.Clone(base).(*pb.WireGuardPeerConfig)}
		cloned.AllowedIps = nil
		if l.GetToEndpoint() != nil {
			cloned.Endpoint = l.GetToEndpoint()
		}
		if _, err := parseAndValidatePeerConfig(cloned); err != nil {
			continue
		}

		log.Debugf("preconnect add: peerID=%d pk=%s endpoint=%s",
			toID, truncate(cloned.GetPublicKey(), 10), endpointForLog(cloned.GetEndpoint()))
		uapiBuilder.AddPeerConfig(cloned)
		w.ifce.Peers = append(w.ifce.Peers, cloned.WireGuardPeerConfig)
		w.indexPeerDirectoryLocked(cloned)
		exists[toID] = struct{}{}
		w.preconnectPeers[toID] = struct{}{}
		added++
	}

	if added == 0 && removed == 0 {
		log.Debugf("ensure result: no-op (skippedAlready=%d skippedNoBase=%d)", skippedAlready, skippedNoBase)
		return nil
	}
	log.Debugf("ensure result: add=%d remove=%d skippedAlready=%d skippedNoBase=%d", added, removed, skippedAlready, skippedNoBase)
	if err := w.wgDevice.IpcSet(uapiBuilder.Build()); err != nil {
		log.WithError(err).Debugf("ensure IpcSet failed (add=%d remove=%d)", added, removed)
		return err
	}
	return nil
}

// mergeConnectablePeersFromAdj 将本节点 adj 图中可直连的 peer 合并进目标 peers。
//
// - **只在目标 peers 中缺失时才补齐**（按 PublicKey 去重），避免覆盖由路由规划器计算出的 AllowedIPs。
// - **补齐的 peer AllowedIPs 置空**，确保不会引入额外路由。
// - 如链路显式携带 to_endpoint，则优先用它覆盖 peer.endpoint（用于快速恢复直连）。
//
// knownPeers 用于在目标 peers 缺失时提供“可用的 peer 基础信息”（公钥/预共享密钥/端点等），通常传 oldPeers。
func mergeConnectablePeersFromAdj(ifce *defs.WireGuardConfig, desiredPeers []*defs.WireGuardPeerConfig, knownPeers []*defs.WireGuardPeerConfig) []*defs.WireGuardPeerConfig {
	if ifce == nil {
		return desiredPeers
	}
	localID := ifce.GetId()
	if localID == 0 {
		return desiredPeers
	}
	adjs := ifce.GetAdjs()
	if adjs == nil {
		return desiredPeers
	}
	localLinks, ok := adjs[localID]
	if !ok || localLinks == nil || len(localLinks.GetLinks()) == 0 {
		return desiredPeers
	}

	desiredByPK := make(map[string]*defs.WireGuardPeerConfig, len(desiredPeers))
	for _, p := range desiredPeers {
		if p == nil || p.GetPublicKey() == "" {
			continue
		}
		desiredByPK[p.GetPublicKey()] = p
	}

	// build id -> peer 基础信息索引（优先 desired，其次 known）
	idToPeer := make(map[uint32]*defs.WireGuardPeerConfig, len(desiredPeers)+len(knownPeers))
	putPeerIDs := func(p *defs.WireGuardPeerConfig) {
		if p == nil {
			return
		}
		// 1) peer.id
		if id := p.GetId(); id != 0 {
			if _, exists := idToPeer[id]; !exists {
				idToPeer[id] = p
			}
		}
		// 2) endpoint.wireguard_id（有些下发场景可能不填 peer.id，但 endpoint 里带 wireguard_id）
		if p.GetEndpoint() != nil {
			if id := p.GetEndpoint().GetWireguardId(); id != 0 {
				if _, exists := idToPeer[id]; !exists {
					idToPeer[id] = p
				}
			}
		}
	}
	for _, p := range desiredPeers {
		putPeerIDs(p)
	}
	for _, p := range knownPeers {
		putPeerIDs(p)
	}

	// 仅补齐：adj 中的直连节点（to_wireguard_id）
	for _, l := range localLinks.GetLinks() {
		if l == nil {
			continue
		}
		toID := l.GetToWireguardId()
		if toID == 0 || toID == localID {
			continue
		}

		base, ok := idToPeer[toID]
		if !ok || base == nil || base.GetPublicKey() == "" {
			continue
		}
		if _, exists := desiredByPK[base.GetPublicKey()]; exists {
			// 已在目标列表中（通常含有路由规划器计算出的 AllowedIPs），不覆盖。
			continue
		}

		// 复制一份（避免直接改 oldPeers / knownPeers 的底层 pb 指针）
		cloned := clonePeerConfig(base)
		// 不分配路由：AllowedIPs 置空
		cloned.AllowedIps = nil
		// 显式链路 endpoint 优先
		if l.GetToEndpoint() != nil {
			cloned.Endpoint = l.GetToEndpoint()
		}
		// 确保 keepalive/AllowedIPs 格式一致（复用现有校验逻辑）
		if _, err := parseAndValidatePeerConfig(cloned); err != nil {
			continue
		}

		desiredPeers = append(desiredPeers, cloned)
		desiredByPK[cloned.GetPublicKey()] = cloned
	}

	return desiredPeers
}

func clonePeerConfig(p *defs.WireGuardPeerConfig) *defs.WireGuardPeerConfig {
	if p == nil || p.WireGuardPeerConfig == nil {
		return &defs.WireGuardPeerConfig{}
	}

	// 使用 proto.Clone 避免直接拷贝 protoimpl.MessageState（内部含 mutex，会触发拷贝锁值的告警）
	cp, _ := proto.Clone(p.WireGuardPeerConfig).(*pb.WireGuardPeerConfig)
	if cp == nil {
		return &defs.WireGuardPeerConfig{}
	}
	return &defs.WireGuardPeerConfig{WireGuardPeerConfig: cp}
}

func endpointForLog(ep *pb.Endpoint) string {
	if ep == nil {
		return ""
	}
	if ep.GetUri() != "" {
		return ep.GetUri()
	}
	if ep.GetHost() != "" || ep.GetPort() != 0 {
		return fmt.Sprintf("%s:%d", ep.GetHost(), ep.GetPort())
	}
	return ""
}
