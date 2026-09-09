package wg

import (
	"math"
	"time"

	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/services/app"
)

// RoutingPolicy 决定边权重的计算方式。
// cost = LatencyTerm + InverseBandwidthTerm + HopWeight + HandshakePenalty
type RoutingPolicy struct {
	LatencyWeight          float64
	InverseBandwidthWeight float64
	HopWeight              float64
	MinUpMbps              uint32

	// LatencyBucketMs 用于对 latency 做“分桶/量化”，降低抖动导致的最短路频繁切换。
	// 例如 bucket=5ms，则 31/32/33ms 都会被量化为 30/35ms 附近的同一档。
	LatencyBucketMs uint32
	// MinLatencyMs/MaxLatencyMs 用于对 latency 做限幅，避免异常值对最短路产生过强扰动。
	MinLatencyMs uint32
	MaxLatencyMs uint32
	// LatencyLogScale>0 时，对 latency 使用 log1p 变换并乘以该系数，使权重对小幅抖动更不敏感。
	// 若为 0，则回退为线性 latency。
	LatencyLogScale float64

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
		MinUpMbps:                1,
		LatencyBucketMs:          5,
		MinLatencyMs:             1,
		MaxLatencyMs:             1500,
		LatencyLogScale:          10.0,
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

// EdgeWeight 计算一条“有向边”的权重（越小越优）。
// 为了抑制延迟探测的噪声导致路由频繁抖动，这里对 latency 做了：限幅 + 分桶（可选）+ log1p（可选）。
func (p *RoutingPolicy) EdgeWeight(fromWGID uint, e Edge, idToPeer map[uint]*models.WireGuard) float64 {
	lat := float64(e.latency)

	// 1) 延迟限幅
	minLat := float64(p.MinLatencyMs)
	maxLat := float64(p.MaxLatencyMs)
	if minLat <= 0 {
		minLat = 1
	}
	if maxLat <= 0 {
		maxLat = 1500
	}
	if lat < minLat {
		lat = minLat
	}
	if lat > maxLat {
		lat = maxLat
	}

	// 2) 延迟分桶（量化）
	if p.LatencyBucketMs > 0 {
		b := float64(p.LatencyBucketMs)
		// 四舍五入到最近桶
		lat = math.Floor((lat+b/2)/b) * b
		if lat < minLat {
			lat = minLat
		}
		if lat > maxLat {
			lat = maxLat
		}
	}

	// 3) 延迟项：log1p（可选）+ scale
	latencyTerm := 0.0
	if p.LatencyWeight != 0 {
		if p.LatencyLogScale > 0 {
			latencyTerm = p.LatencyWeight * math.Log1p(lat) * p.LatencyLogScale
		} else {
			latencyTerm = p.LatencyWeight * lat
		}
	}

	// 4) 带宽项：对低带宽更敏感；使用 MinUpMbps 做下限避免极端值
	minUp := float64(p.MinUpMbps)
	if minUp <= 0 {
		minUp = 1
	}
	up := math.Max(float64(e.upMbps), minUp)
	invBw := 1.0 / math.Max(up, 1e-6)
	bwTerm := p.InverseBandwidthWeight * invBw

	// 5) hop 项
	hopTerm := p.HopWeight

	// 6) 握手过旧惩罚：必须无方向
	handshakePenalty := 0.0
	if p.HandshakeStalePenalty > 0 && p.HandshakeStaleThreshold > 0 {
		if age, ok := getHandshakeAgeBetween(fromWGID, e.to, idToPeer, *p); ok && age > p.HandshakeStaleThreshold {
			handshakePenalty = p.HandshakeStalePenalty
		}
	}

	return latencyTerm + bwTerm + hopTerm + handshakePenalty
}
