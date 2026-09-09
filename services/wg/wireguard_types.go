//go:build !windows
// +build !windows

package wg

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun"
	"golang.zx2c4.com/wireguard/tun/netstack"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/wg/multibind"
	"github.com/VaalaCat/frp-panel/utils"
)

var (
	_ app.WireGuard = (*wireGuard)(nil)
)

type wireGuard struct {
	sync.RWMutex

	ifce            *defs.WireGuardConfig
	endpointPingMap *utils.SyncMap[uint32, uint32] // ms
	virtAddrPingMap *utils.SyncMap[string, uint32] // ms
	// ping 平滑器：对“瞬时探测值”做 EWMA 聚合，降低抖动
	pingAggMu         sync.Mutex
	endpointPingEWMA  map[uint32]float64 // peerID -> ema(ms)
	virtAddrPingEWMA  map[string]float64 // virtAddr -> ema(ms)
	peerDirectory   map[uint32]*pb.WireGuardPeerConfig
	// 仅用于“预连接/保持连接”的 peer（AllowedIPs 为空），用于后续根据拓扑变化做增删
	preconnectPeers map[uint32]struct{}

	wgDevice  *device.Device
	tunDevice tun.Device
	multiBind *multibind.MultiBind
	gvisorNet *netstack.Net
	fwManager *firewallManager

	running      bool
	useGvisorNet bool // if true, use gvisor netstack

	svcLogger *logrus.Entry
	ctx       *app.Context
	cancel    context.CancelFunc
}
