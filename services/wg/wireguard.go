package wg

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"os"
	"reflect"
	"sync"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun"
	"golang.zx2c4.com/wireguard/tun/netstack"
	"gvisor.dev/gvisor/pkg/tcpip"
	"gvisor.dev/gvisor/pkg/tcpip/network/ipv4"
	"gvisor.dev/gvisor/pkg/tcpip/network/ipv6"
	"gvisor.dev/gvisor/pkg/tcpip/stack"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/wg/multibind"
	"github.com/VaalaCat/frp-panel/services/wg/transport/ws"
	"github.com/VaalaCat/frp-panel/utils"
)

const (
	ReportInterval = time.Second * 60
)

var (
	_ app.WireGuard = (*wireGuard)(nil)
)

type wireGuard struct {
	sync.RWMutex

	ifce            *defs.WireGuardConfig
	endpointPingMap *utils.SyncMap[uint32, uint32] // ms
	virtAddrPingMap *utils.SyncMap[string, uint32] // ms

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

func NewWireGuard(ctx *app.Context, ifce defs.WireGuardConfig, logger *logrus.Entry) (app.WireGuard, error) {
	if logger == nil {
		defaultLog := logrus.New()
		logger = logrus.NewEntry(defaultLog)
	}
	cfg := ifce

	if err := InitAndValidateWGConfig(&cfg, logger); err != nil {
		return nil, errors.Join(errors.New("init and validate wg config error"), err)
	}

	svcCtx, cancel := ctx.CopyWithCancel()

	useGvisorNet := ctx.GetApp().GetConfig().App.UseGvisorNet
	if !useGvisorNet {
		useGvisorNet = cfg.GetUseGvisorNet()
	}

	fwManager := newFirewallManager(logger.WithField("component", "iptables"))

	return &wireGuard{
		RWMutex:         sync.RWMutex{},
		ifce:            &cfg,
		ctx:             svcCtx,
		cancel:          cancel,
		svcLogger:       logger,
		endpointPingMap: &utils.SyncMap[uint32, uint32]{},
		useGvisorNet:    useGvisorNet,
		virtAddrPingMap: &utils.SyncMap[string, uint32]{},
		fwManager:       fwManager,
	}, nil
}

// Start implements WireGuard.
func (w *wireGuard) Start() error {
	w.Lock()
	defer w.Unlock()

	log := w.svcLogger.WithField("op", "Start")

	if w.running {
		log.Warnf("wireguard is already running, skip start, ifce: %s", w.ifce.GetInterfaceName())
		return nil
	}

	if err := w.initTransports(); err != nil {
		return errors.Join(fmt.Errorf("init transports failed"), err)
	}

	if err := w.initWGDevice(); err != nil {
		return errors.Join(fmt.Errorf("init WG device failed"), err)
	}

	if err := w.applyPeerConfig(); err != nil {
		return errors.Join(fmt.Errorf("apply peer config failed"), err)
	}

	// 在应用配置后再启动设备
	if err := w.wgDevice.Up(); err != nil {
		return errors.Join(fmt.Errorf("wgDevice.Up '%s'", w.ifce.GetInterfaceName()), err)
	}

	if !w.useGvisorNet {
		if err := w.initNetwork(); err != nil {
			return errors.Join(errors.New("init network failed"), err)
		}
	}

	// 在 WireGuard 设备启动后配置 gvisor
	if w.useGvisorNet {
		if err := w.initGvisorNetwork(); err != nil {
			return errors.Join(errors.New("init gvisor network failed"), err)
		}
	}

	if err := w.applyFirewallRulesLocked(); err != nil {
		return errors.Join(errors.New("apply firewall rules failed"), err)
	}

	log.Infof("Started service done for iface '%s'", w.ifce.GetInterfaceName())
	w.running = true

	go w.reportStatusTask()

	return nil
}

// Stop implements WireGuard.
func (w *wireGuard) Stop() error {
	w.Lock()
	defer w.Unlock()

	log := w.svcLogger.WithField("op", "Stop")
	if !w.running {
		log.Info("Service already down.")
		return nil
	}

	log.Info("Stopping service...")

	if err := w.cleanupFirewallRulesLocked(); err != nil {
		log.WithError(err).Warn("cleanup firewall rules failed")
	}

	w.cleanupWGDevice()
	w.cleanupNetwork()
	w.cancel()

	w.running = false
	log.Info("Service stopped.")
	return nil
}

// AddPeer implements WireGuard.
func (w *wireGuard) AddPeer(peer *defs.WireGuardPeerConfig) error {
	log := w.svcLogger.WithField("op", "AddPeer")

	peerCfg, err := parseAndValidatePeerConfig(peer)
	if err != nil {
		return errors.Join(errors.New("parse and validate peer config"), err)
	}

	w.Lock()
	defer w.Unlock()

	w.ifce.Peers = append(w.ifce.Peers, peer.WireGuardPeerConfig)
	uapiBuilder := NewUAPIBuilder().AddPeerConfig(peerCfg)

	log.Debugf("uapiBuilder: %s", uapiBuilder.Build())

	if err = w.wgDevice.IpcSet(uapiBuilder.Build()); err != nil {
		return errors.Join(errors.New("add peer IpcSet error"), err)
	}

	return nil
}

// GenWGConfig implements WireGuard.
func (w *wireGuard) GenWGConfig() (string, error) {
	panic("unimplemented")
}

// GetIfceConfig implements WireGuard.
func (w *wireGuard) GetIfceConfig() (*defs.WireGuardConfig, error) {
	w.RLock()
	defer w.RUnlock()

	return w.ifce, nil
}

// GetBaseIfceConfig implements WireGuard.
func (w *wireGuard) GetBaseIfceConfig() *defs.WireGuardConfig {
	return w.ifce
}

// GetPeer implements WireGuard.
func (w *wireGuard) GetPeer(peerNameOrPk string) (*defs.WireGuardPeerConfig, error) {
	w.RLock()
	defer w.RUnlock()

	for _, p := range w.ifce.Peers {
		if p.ClientId == peerNameOrPk || p.PublicKey == peerNameOrPk {
			return &defs.WireGuardPeerConfig{WireGuardPeerConfig: p}, nil
		}
	}
	return nil, errors.New("peer not found")
}

// ListPeers implements WireGuard.
func (w *wireGuard) ListPeers() ([]*defs.WireGuardPeerConfig, error) {
	w.RLock()
	defer w.RUnlock()

	return w.ifce.GetParsedPeers(), nil
}

// RemovePeer implements WireGuard.
func (w *wireGuard) RemovePeer(peerNameOrPk string) error {
	log := w.svcLogger.WithField("op", "RemovePeer")

	w.Lock()
	defer w.Unlock()

	newPeers := []*pb.WireGuardPeerConfig{}
	var peerToRemove *defs.WireGuardPeerConfig
	for _, p := range w.ifce.Peers {
		if p.ClientId != peerNameOrPk && p.PublicKey != peerNameOrPk {
			newPeers = append(newPeers, p)
			continue
		}
		peerToRemove = &defs.WireGuardPeerConfig{WireGuardPeerConfig: p}
	}

	if len(newPeers) == len(w.ifce.Peers) {
		return errors.New("peer not found")
	}

	w.ifce.Peers = newPeers

	uapiBuilder := NewUAPIBuilder().RemovePeerByKey(peerToRemove.GetParsedPublicKey())

	log.Debugf("uapiBuilder: %s", uapiBuilder.Build())

	if err := w.wgDevice.IpcSet(uapiBuilder.Build()); err != nil {
		return errors.Join(errors.New("remove peer IpcSet error"), err)
	}

	return nil
}

// UpdatePeer implements WireGuard.
func (w *wireGuard) UpdatePeer(peer *defs.WireGuardPeerConfig) error {
	log := w.svcLogger.WithField("op", "UpdatePeer")

	peerCfg, err := parseAndValidatePeerConfig(peer)
	if err != nil {
		return errors.Join(errors.New("parse and validate peer config"), err)
	}

	w.Lock()
	defer w.Unlock()

	newPeers := []*pb.WireGuardPeerConfig{}
	for _, p := range w.ifce.Peers {
		if p.ClientId != peer.ClientId && p.PublicKey != peer.PublicKey {
			newPeers = append(newPeers, p)
			continue
		}
		newPeers = append(newPeers, peer.WireGuardPeerConfig)
	}

	w.ifce.Peers = newPeers

	uapiBuilder := NewUAPIBuilder().UpdatePeerConfig(peerCfg)

	log.Debugf("uapiBuilder: %s", uapiBuilder.Build())

	if err := w.wgDevice.IpcSet(uapiBuilder.Build()); err != nil {
		return errors.Join(errors.New("update peer IpcSet error"), err)
	}

	return nil
}

func (w *wireGuard) PatchPeers(newPeers []*defs.WireGuardPeerConfig) (*app.WireGuardDiffPeersResponse, error) {
	oldPeers := w.ifce.GetParsedPeers()

	diffResp := utils.Diff(oldPeers, newPeers)

	resp := &app.WireGuardDiffPeersResponse{
		AddPeers:    diffResp.NotInArr1,
		RemovePeers: diffResp.NotInArr2,
	}

	for _, peer := range resp.RemovePeers {
		w.RemovePeer(peer.GetPublicKey())
	}
	for _, peer := range resp.AddPeers {
		w.AddPeer(peer)
	}

	return resp, nil
}

func (w *wireGuard) GetWGRuntimeInfo() (*pb.WGDeviceRuntimeInfo, error) {
	runningInfo, err := w.wgDevice.IpcGet()
	if err != nil {
		return nil, errors.Join(errors.New("get WG running info error"), err)
	}

	runtimeInfo, err := ParseWGRunningInfo(runningInfo)
	if err != nil {
		return nil, err
	}

	runtimeInfo.PingMap = w.endpointPingMap.Export()
	runtimeInfo.VirtAddrPingMap = w.virtAddrPingMap.Export()

	if w.useGvisorNet {
		runtimeInfo.InterfaceName = w.ifce.GetInterfaceName()
	} else {
		link, err := netlink.LinkByName(w.ifce.GetInterfaceName())
		if err != nil {
			return nil, errors.Join(fmt.Errorf("get iface '%s' via netlink", w.ifce.GetInterfaceName()), err)
		}
		runtimeInfo.InterfaceName = link.Attrs().Name
	}

	parsedPeers := w.ifce.GetParsedPeers()
	parsedPublicKeysPeerMap := make(map[string]*defs.WireGuardPeerConfig)
	for _, peer := range parsedPeers {
		parsedPublicKeysPeerMap[peer.HexPublicKey()] = peer
	}

	runtimeInfo.PeerVirtAddrMap = make(map[string]uint32)
	for _, peer := range parsedPeers {
		runtimeInfo.PeerVirtAddrMap[peer.GetVirtualIp()] = peer.GetId()
	}

	for _, peerRuntimeInfo := range runtimeInfo.GetPeers() {
		peerConfig, ok := parsedPublicKeysPeerMap[peerRuntimeInfo.PublicKey]
		if !ok {
			continue
		}
		if runtimeInfo.PeerConfigMap == nil {
			runtimeInfo.PeerConfigMap = make(map[string]*pb.WireGuardPeerConfig)
		}
		runtimeInfo.PeerConfigMap[peerRuntimeInfo.PublicKey] = peerConfig.WireGuardPeerConfig
		peerRuntimeInfo.ClientId = peerConfig.GetClientId()
	}

	return runtimeInfo, nil
}

func (w *wireGuard) UpdateAdjs(adjs map[uint32]*pb.WireGuardLinks) error {
	w.Lock()
	defer w.Unlock()

	w.ifce.Adjs = adjs
	return nil
}

func (w *wireGuard) NeedRecreate(newCfg *defs.WireGuardConfig) bool {
	return w.ifce.GetId() != newCfg.GetId() ||
		w.ifce.GetInterfaceName() != newCfg.GetInterfaceName() ||
		w.ifce.GetPrivateKey() != newCfg.GetPrivateKey() ||
		w.ifce.GetLocalAddress() != newCfg.GetLocalAddress() ||
		w.ifce.GetListenPort() != newCfg.GetListenPort() ||
		w.ifce.GetWsListenPort() != newCfg.GetWsListenPort() ||
		w.ifce.GetInterfaceMtu() != newCfg.GetInterfaceMtu() ||
		w.ifce.GetUseGvisorNet() != newCfg.GetUseGvisorNet() ||
		w.ifce.GetNetworkId() != newCfg.GetNetworkId()
}

func (w *wireGuard) initTransports() error {
	log := w.svcLogger.WithField("op", "initTransports")

	wsTrans := ws.NewWSBind(w.ctx)
	w.multiBind = multibind.NewMultiBind(
		w.svcLogger,
		multibind.NewTransport(conn.NewDefaultBind(), "udp"),
		multibind.NewTransport(wsTrans, "ws"),
	)

	engine := gin.New()
	engine.Any(defs.DefaultWSHandlerPath, func(c *gin.Context) {
		err := wsTrans.HandleHTTP(c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	})

	// if ws listen port not set, use wg listen port, share tcp and udp port
	listenPort := w.ifce.GetWsListenPort()
	if listenPort == 0 {
		listenPort = w.ifce.GetListenPort()
	}
	go func() {
		if err := engine.Run(fmt.Sprintf(":%d", listenPort)); err != nil {
			w.svcLogger.WithError(err).Errorf("failed to run gin engine for ws transport on port %d", listenPort)
		}
	}()

	log.Infof("WS transport engine running on port %d", listenPort)

	return nil
}

func (w *wireGuard) initWGDevice() error {
	log := w.svcLogger.WithField("op", "initWGDevice")

	log.Debugf("start to create TUN device '%s' (MTU %d)", w.ifce.GetInterfaceName(), w.ifce.GetInterfaceMtu())

	var err error

	if w.useGvisorNet {
		log.Infof("using gvisor netstack for TUN device")
		prf, err := netip.ParsePrefix(w.ifce.GetLocalAddress())
		if err != nil {
			return errors.Join(fmt.Errorf("parse local addr '%s' for netip", w.ifce.GetLocalAddress()), err)
		}

		addrs := lo.Map(w.ifce.GetDnsServers(), func(s string, _ int) netip.Addr {
			addr, err := netip.ParseAddr(s)
			if err != nil {
				return netip.Addr{}
			}
			return addr
		})
		if len(addrs) == 0 {
			addrs = []netip.Addr{netip.AddrFrom4([4]byte{1, 2, 4, 8})}
		}
		log.Debugf("create netstack TUN with addr '%s' and dns servers '%v'", prf.Addr().String(), addrs)
		w.tunDevice, w.gvisorNet, err = netstack.CreateNetTUN([]netip.Addr{prf.Addr()}, addrs, 1200)
		if err != nil {
			return errors.Join(fmt.Errorf("create netstack TUN device '%s' (MTU %d) failed", w.ifce.GetInterfaceName(), w.ifce.GetInterfaceMtu()), err)
		}
	} else {
		w.tunDevice, err = tun.CreateTUN(w.ifce.GetInterfaceName(), int(w.ifce.GetInterfaceMtu()))
		if err != nil {
			return errors.Join(fmt.Errorf("create TUN device '%s' (MTU %d) failed", w.ifce.GetInterfaceName(), w.ifce.GetInterfaceMtu()), err)
		}
	}

	log.Debugf("TUN device '%s' (MTU %d) created successfully", w.ifce.GetInterfaceName(), w.ifce.GetInterfaceMtu())

	log.Debugf("start to create WireGuard device '%s'", w.ifce.GetInterfaceName())

	w.wgDevice = device.NewDevice(w.tunDevice, w.multiBind, &device.Logger{
		Verbosef: w.svcLogger.WithField("wg-dev-iface", w.ifce.GetInterfaceName()).Debugf,
		Errorf:   w.svcLogger.WithField("wg-dev-iface", w.ifce.GetInterfaceName()).Errorf,
	})

	log.Debugf("WireGuard device '%s' created successfully", w.ifce.GetInterfaceName())

	return nil
}

func (w *wireGuard) applyPeerConfig() error {
	log := w.svcLogger.WithField("op", "applyConfig")

	log.Debugf("start to apply config to WireGuard device '%s'", w.ifce.GetInterfaceName())

	if w.wgDevice == nil {
		return errors.New("wgDevice is nil, please init WG device first")
	}

	wgTypedPeerConfigs, err := parseAndValidatePeerConfigs(w.ifce.GetParsedPeers())
	if err != nil {
		return errors.Join(errors.New("parse/validate peers"), err)
	}

	log.Debugf("wgTypedPeerConfigs: %v", wgTypedPeerConfigs)

	uapiConfigString := generateUAPIConfigString(w.ifce, w.ifce.GetParsedPrivKey(), wgTypedPeerConfigs, !w.running, false)

	log.Debugf("uapiBuilder: %s", uapiConfigString)

	log.Debugf("calling IpcSet...")
	if err = w.wgDevice.IpcSet(uapiConfigString); err != nil {
		return errors.Join(errors.New("IpcSet error"), err)
	}
	log.Debugf("IpcSet completed successfully")

	return nil
}

func (w *wireGuard) initNetwork() error {
	log := w.svcLogger.WithField("op", "initNetwork")

	// 等待 TUN 设备在内核中完全注册,避免竞态条件
	var link netlink.Link
	var err error
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		link, err = netlink.LinkByName(w.ifce.GetInterfaceName())
		if err == nil {
			break
		}
		if i < maxRetries-1 {
			log.Debugf("attempt %d: waiting for iface '%s' to be ready, will retry...", i+1, w.ifce.GetInterfaceName())
			time.Sleep(100 * time.Millisecond)
		}
	}
	if err != nil {
		return errors.Join(fmt.Errorf("get iface '%s' via netlink after %d retries", w.ifce.GetInterfaceName(), maxRetries), err)
	}
	log.Debugf("successfully found interface '%s' via netlink", w.ifce.GetInterfaceName())

	addr, err := netlink.ParseAddr(w.ifce.GetLocalAddress())
	if err != nil {
		return errors.Join(fmt.Errorf("parse local addr '%s' for netlink", w.ifce.GetLocalAddress()), err)
	}

	if err = netlink.AddrAdd(link, addr); err != nil && !os.IsExist(err) {
		return errors.Join(fmt.Errorf("add IP '%s' to '%s'", w.ifce.GetLocalAddress(), w.ifce.GetInterfaceName()), err)
	} else if os.IsExist(err) {
		log.Infof("IP %s already on '%s'.", w.ifce.GetLocalAddress(), w.ifce.GetInterfaceName())
	} else {
		log.Infof("IP %s added to '%s'.", w.ifce.GetLocalAddress(), w.ifce.GetInterfaceName())
	}

	if err = netlink.LinkSetMTU(link, int(w.ifce.GetInterfaceMtu())); err != nil {
		log.Warnf("Set MTU %d on '%s' via netlink: %v. TUN MTU is %d.",
			w.ifce.GetInterfaceMtu(), w.ifce.GetInterfaceName(), err, w.ifce.GetInterfaceMtu())
	} else {
		log.Infof("Iface '%s' MTU %d set via netlink.", w.ifce.GetInterfaceName(), w.ifce.GetInterfaceMtu())
	}

	if err = netlink.LinkSetUp(link); err != nil {
		return errors.Join(fmt.Errorf("bring up iface '%s' via netlink", w.ifce.GetInterfaceName()), err)
	}

	log.Infof("Iface '%s' up via netlink.", w.ifce.GetInterfaceName())
	return nil
}

func (w *wireGuard) initGvisorNetwork() error {
	log := w.svcLogger.WithField("op", "initGvisorNetwork")

	if w.gvisorNet == nil {
		return errors.New("gvisorNet is nil, cannot initialize network")
	}

	// wg-go dose not expose the stack field, so we need to use reflection to access it
	netValue := reflect.ValueOf(w.gvisorNet).Elem()
	stackField := netValue.FieldByName("stack")

	if !stackField.IsValid() {
		return errors.New("cannot find stack field in gvisorNet")
	}

	stackPtrValue := reflect.NewAt(stackField.Type(), unsafe.Pointer(stackField.UnsafeAddr())).Elem()
	if !stackPtrValue.IsValid() || stackPtrValue.IsNil() {
		return errors.New("gvisor stack is nil or invalid")
	}

	gvisorStack := stackPtrValue.Interface().(*stack.Stack)
	if gvisorStack == nil {
		return errors.New("gvisor stack is nil after conversion")
	}

	log.Infof("successfully accessed gvisor stack, enabling IP forwarding")

	if err := gvisorStack.SetForwardingDefaultAndAllNICs(ipv4.ProtocolNumber, true); err != nil {
		log.Warnf("failed to enable IPv4 forwarding: %v, relay may not work", err)
	} else {
		log.Infof("IPv4 forwarding enabled for gvisor netstack")
	}

	if err := gvisorStack.SetForwardingDefaultAndAllNICs(ipv6.ProtocolNumber, true); err != nil {
		log.Warnf("failed to enable IPv6 forwarding: %v", err)
	} else {
		log.Infof("IPv6 forwarding enabled for gvisor netstack")
	}

	for _, peer := range w.ifce.Peers {
		for _, allowedIP := range peer.AllowedIps {
			prefix, err := netip.ParsePrefix(allowedIP)
			if err != nil {
				log.WithError(err).Warnf("failed to parse allowed IP: %s", allowedIP)
				continue
			}

			addr := tcpip.AddrFromSlice(prefix.Addr().AsSlice())

			ones := prefix.Bits()
			maskBytes := make([]byte, len(prefix.Addr().AsSlice()))
			for i := 0; i < len(maskBytes); i++ {
				if ones >= 8 {
					maskBytes[i] = 0xff
					ones -= 8
				} else if ones > 0 {
					maskBytes[i] = byte(0xff << (8 - ones))
					ones = 0
				}
			}

			subnet, err := tcpip.NewSubnet(addr, tcpip.MaskFromBytes(maskBytes))
			if err != nil {
				log.WithError(err).Warnf("failed to create subnet for %s", allowedIP)
				continue
			}

			route := tcpip.Route{
				Destination: subnet,
				NIC:         1,
			}

			gvisorStack.AddRoute(route)
			log.Debugf("added route for peer allowed IP: %s via NIC 1", allowedIP)
		}
	}

	log.Infof("gvisor netstack initialized with IP forwarding enabled")
	return nil
}

func (w *wireGuard) cleanupNetwork() {
	log := w.svcLogger.WithField("op", "cleanupNetwork")

	if w.useGvisorNet {
		log.Infof("skip network cleanup for gvisor netstack")
		return
	}

	link, err := netlink.LinkByName(w.ifce.GetInterfaceName())
	if err == nil {
		if err := netlink.LinkSetDown(link); err != nil {
			log.Warnf("Failed to LinkSetDown '%s' after wgDevice.Up() error: %v", w.ifce.GetInterfaceName(), err)
		}
	}
	log.Debug("Cleanup network complete.")
}

func (w *wireGuard) cleanupWGDevice() {
	log := w.svcLogger.WithField("op", "cleanupWGDevice")

	if w.wgDevice != nil {
		w.wgDevice.Close()
	} else if w.tunDevice != nil {
		w.tunDevice.Close()
	}
	w.wgDevice = nil
	w.tunDevice = nil
	log.Debug("Cleanup WG device complete.")
}

func (w *wireGuard) applyFirewallRulesLocked() error {
	if w.useGvisorNet || w.fwManager == nil {
		return nil
	}

	prefix, err := netip.ParsePrefix(w.ifce.GetLocalAddress())
	if err != nil {
		return errors.Join(fmt.Errorf("parse local address '%s' for firewall", w.ifce.GetLocalAddress()), err)
	}

	return w.fwManager.ApplyRelayRules(w.ifce.GetInterfaceName(), prefix.Masked().String())
}

func (w *wireGuard) cleanupFirewallRulesLocked() error {
	if w.useGvisorNet || w.fwManager == nil {
		return nil
	}
	return w.fwManager.Cleanup(w.ifce.GetInterfaceName())
}

func (w *wireGuard) reportStatusTask() {
	for {
		select {
		case <-w.ctx.Done():
			return
		default:
			w.pingPeers()
			time.Sleep(ReportInterval)
		}
	}
}
