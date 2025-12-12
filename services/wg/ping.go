//go:build !windows
// +build !windows

package wg

import (
	"errors"
	"fmt"
	"math"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	probing "github.com/prometheus-community/pro-bing"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/pb"
)

const (
	endpointPingCount   = 5
	endpointPingTimeout = 10 * time.Second
)

func (w *wireGuard) pingPeers() {
	log := w.svcLogger.WithField("op", "pingPeers")

	ifceConfig, err := w.GetIfceConfig()
	if err != nil {
		log.WithError(err).Errorf("failed to get interface config")
		return
	}

	log.Debugf("start to ping peers, len: %d", len(ifceConfig.Peers))

	var waitGroup conc.WaitGroup
	w.scheduleEndpointPings(log, ifceConfig, &waitGroup)
	w.scheduleVirtualAddrPings(log, ifceConfig, &waitGroup)
	w.waitPingTasks(log, &waitGroup)
}

func (w *wireGuard) scheduleEndpointPings(log *logrus.Entry, ifceConfig *defs.WireGuardConfig, waitGroup *conc.WaitGroup) {
	targets := collectEndpointPingTargets(ifceConfig)
	if len(targets) == 0 {
		return
	}

	log.Debugf("schedule endpoint pings, targets=%d", len(targets))

	for peerID, endpoint := range targets {
		peerId := peerID
		ep := endpoint
		if ep == nil {
			continue
		}

		// ws endpoint 不走 ICMP ping，改为 TCP connect 探测，避免误报/不可达。
		if endpointTypeContainsWS(ep.GetType()) {
			tcpAddr, err := endpointTCPTarget(ep)
			if err != nil {
				log.WithError(err).Errorf("failed to resolve tcp target for endpoint, peer_id=%d, endpoint=%+v", peerId, ep)
				w.storeEndpointPing(peerId, math.MaxUint32)
				continue
			}

			waitGroup.Go(func() {
				avg, err := tcpPingAvg(tcpAddr, endpointPingCount, endpointPingTimeout)
				if err != nil {
					log.WithError(err).Errorf("tcp ping endpoint [%s] failed, peer_id=%d", tcpAddr, peerId)
					w.storeEndpointPing(peerId, math.MaxUint32)
					return
				}
				avgMs := uint32(avg.Milliseconds())
				if avgMs == 0 { // 0 means bug
					avgMs = 1
				}
				w.storeEndpointPing(peerId, avgMs)
				log.Debugf("tcp ping endpoint [%s] completed, peer_id=%d", tcpAddr, peerId)
			})
			continue
		}

		host := endpointICMPHost(ep)
		if host == "" {
			continue
		}

		epPinger, err := probing.NewPinger(host)
		if err != nil {
			log.WithError(err).Errorf("failed to create pinger for %s", host)
			continue
		}

		epPinger.Count = endpointPingCount
		epPinger.Timeout = endpointPingTimeout

		epPinger.OnFinish = func(stats *probing.Statistics) {
			// stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss
			// stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt
			log.Tracef("ping stats for %s: %v", host, stats)
			avgRttMs := uint32(stats.AvgRtt.Milliseconds())
			w.storeEndpointPing(peerId, avgRttMs)
		}

		epPinger.OnRecv = func(pkt *probing.Packet) {
			log.Tracef("recv from %s", pkt.IPAddr.String())
		}

		epPinger.OnSendError = func(_ *probing.Packet, err error) {
			log.WithError(err).Errorf("failed to send packet to %s", host)
			w.storeEndpointPing(peerId, math.MaxUint32)
		}

		waitGroup.Go(func() {
			if err := epPinger.Run(); err != nil {
				log.WithError(err).Errorf("failed to run pinger for %s", host)
				w.storeEndpointPing(peerId, math.MaxUint32)
				return
			}
			log.Debugf("ping endpoint [%s] completed, peer_id=%d", host, peerId)
		})
	}
}

func (w *wireGuard) scheduleVirtualAddrPings(log *logrus.Entry, ifceConfig *defs.WireGuardConfig, waitGroup *conc.WaitGroup) {
	if w.useGvisorNet {
		return
	}

	peers := ifceConfig.Peers
	for _, peer := range peers {
		p := peer
		addr := p.GetVirtualIp()
		if addr == "" {
			continue
		}

		tcpAddr, err := peerVirtualWSTCPTarget(p, addr)
		if err != nil {
			log.WithError(err).Errorf("failed to build tcp target for virt addr %s", addr)
			continue
		}

		waitGroup.Go(func() {
			avg, err := tcpPingAvg(tcpAddr, endpointPingCount, endpointPingTimeout)
			if err != nil {
				log.WithError(err).Errorf("failed to tcp ping virt addr %s via %s", addr, tcpAddr)
				if w.virtAddrPingMap != nil {
					w.virtAddrPingMap.Store(addr, math.MaxUint32)
				}
				return
			}

			log.Tracef("tcp ping stats for %s via %s: avg=%s", addr, tcpAddr, avg)
			avgRttMs := uint32(avg.Milliseconds())
			if w.virtAddrPingMap != nil {
				w.virtAddrPingMap.Store(addr, avgRttMs)
			}
			log.Debugf("tcp ping virt addr [%s] completed via %s", addr, tcpAddr)
		})
	}
}

func (w *wireGuard) waitPingTasks(log *logrus.Entry, waitGroup *conc.WaitGroup) {
	log.Debugf("wait for pingers to complete")
	rcs := waitGroup.WaitAndRecover()
	if rcs != nil {
		log.WithError(rcs.AsError()).Errorf("failed to wait for pingers")
	}
}

func (w *wireGuard) storeEndpointPing(peerID uint32, ms uint32) {
	if w.endpointPingMap == nil {
		return
	}
	w.endpointPingMap.Store(peerID, ms)
}

// collectEndpointPingTargets 收集所有“可能直连”的节点 endpoint（高内聚：只关注 ping 需要的目标集合）。
//
// 优先级：
// 1) adjs[localID] 中 link.to_endpoint（显式链路可指定 endpoint）
// 2) peers 中的 peer.endpoint（兜底：已下发 peer config 的 endpoint）
//
// 注意：当前 runtimeInfo.ping_map 的 key 只有 wireguardId，因此这里对同一 peerID 只保留一个 endpoint。
func collectEndpointPingTargets(ifceConfig *defs.WireGuardConfig) map[uint32]*pb.Endpoint {
	if ifceConfig == nil {
		return nil
	}

	targets := make(map[uint32]*pb.Endpoint, 32)

	// 1) 从 adj 图里拿：本节点（ifceConfig.Id）可直连的边的目标 endpoint
	localID := ifceConfig.GetId()
	if localID != 0 {
		if adjs := ifceConfig.GetAdjs(); adjs != nil {
			if links, ok := adjs[localID]; ok && links != nil {
				for _, l := range links.GetLinks() {
					toID := l.GetToWireguardId()
					if toID == 0 || toID == localID {
						continue
					}
					if l.GetToEndpoint() == nil {
						continue
					}
					// 显式链路优先：直接覆盖
					targets[toID] = l.GetToEndpoint()
				}
			}
		}
	}

	// 2) 从 peers 列表兜底补齐
	for _, peer := range ifceConfig.GetPeers() {
		if peer == nil {
			continue
		}
		peerID := peer.GetId()
		if peerID == 0 || peerID == localID {
			continue
		}
		if _, exists := targets[peerID]; exists {
			continue
		}
		if peer.GetEndpoint() == nil {
			continue
		}
		targets[peerID] = peer.GetEndpoint()
	}

	return targets
}

func endpointTypeContainsWS(endpointType string) bool {
	return strings.Contains(strings.ToLower(endpointType), "ws")
}

func endpointICMPHost(ep *pb.Endpoint) string {
	if ep == nil {
		return ""
	}

	// 优先使用 Host（兼容老数据），并剥离可能的端口。
	if ep.GetHost() != "" {
		if host, _, err := net.SplitHostPort(ep.GetHost()); err == nil && host != "" {
			return host
		}
		return ep.GetHost()
	}

	// 兜底：从 Uri 提取 hostname
	if ep.GetUri() != "" {
		u, err := url.Parse(ep.GetUri())
		if err == nil {
			if hn := u.Hostname(); hn != "" {
				return hn
			}
		}
	}

	return ""
}

func endpointTCPTarget(ep *pb.Endpoint) (string, error) {
	if ep == nil {
		return "", errors.New("nil endpoint")
	}

	// 优先使用 Uri（ws/wss 场景更准确）
	if ep.GetUri() != "" {
		u, err := url.Parse(ep.GetUri())
		if err != nil {
			return "", errors.Join(fmt.Errorf("parse uri '%s'", ep.GetUri()), err)
		}
		if u.Host == "" {
			return "", fmt.Errorf("empty host in uri '%s'", ep.GetUri())
		}
		if u.Port() != "" {
			return u.Host, nil // 已含端口
		}
		port := defaultPortForScheme(u.Scheme)
		if port == 0 && ep.GetPort() != 0 {
			port = ep.GetPort()
		}
		if port == 0 {
			return "", fmt.Errorf("missing port for uri '%s'", ep.GetUri())
		}
		return net.JoinHostPort(u.Hostname(), strconv.FormatUint(uint64(port), 10)), nil
	}

	host := strings.TrimSpace(ep.GetHost())
	if host == "" {
		return "", errors.New("empty endpoint host")
	}

	// host 已经带端口
	if _, _, err := net.SplitHostPort(host); err == nil {
		return host, nil
	}

	port := ep.GetPort()
	if port == 0 {
		// ws 类型未显式配置端口时，默认按 ws=80 处理
		port = 80
	}
	return net.JoinHostPort(host, strconv.FormatUint(uint64(port), 10)), nil
}

func defaultPortForScheme(scheme string) uint32 {
	switch strings.ToLower(scheme) {
	case "ws", "http":
		return 80
	case "wss", "https":
		return 443
	default:
		return 0
	}
}

func tcpPingAvg(addr string, count int, timeout time.Duration) (time.Duration, error) {
	if count <= 0 {
		return 0, errors.New("invalid count")
	}
	if timeout <= 0 {
		return 0, errors.New("invalid timeout")
	}

	var (
		ok    int
		sum   time.Duration
		dial  = &net.Dialer{Timeout: timeout}
		sleep = 100 * time.Millisecond
	)

	for i := 0; i < count; i++ {
		start := time.Now()
		conn, err := dial.Dial("tcp", addr)
		if err == nil {
			_ = conn.Close()
			sum += time.Since(start)
			ok++
		}
		// 轻微抖动，避免瞬时突刺；保持逻辑简单，不引入全局 rate limiter。
		if i != count-1 {
			time.Sleep(sleep)
		}
	}

	if ok == 0 {
		return 0, fmt.Errorf("all tcp probes failed for %s", addr)
	}
	return sum / time.Duration(ok), nil
}

func peerVirtualWSTCPTarget(peer *pb.WireGuardPeerConfig, virtIP string) (string, error) {
	if peer == nil {
		return "", errors.New("nil peer")
	}
	if virtIP == "" {
		return "", errors.New("empty virt ip")
	}

	wsPort := peer.GetWsListenPort()
	if wsPort == 0 {
		// 兼容“ws 端口未单独配置时复用 listen_port”的逻辑（见 initTransports）
		wsPort = peer.GetListenPort()
	}
	if wsPort == 0 {
		return "", fmt.Errorf("missing ws port for peer_id=%d virt_ip=%s", peer.GetId(), virtIP)
	}

	return net.JoinHostPort(virtIP, strconv.FormatUint(uint64(wsPort), 10)), nil
}
