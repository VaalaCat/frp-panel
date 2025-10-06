package wg

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	probing "github.com/prometheus-community/pro-bing"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc"
	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
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

	ifce    *defs.WireGuardConfig
	pingMap *utils.SyncMap[uint32, uint32] // ms

	wgDevice  *device.Device
	tunDevice tun.Device

	running bool

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

	return &wireGuard{
		RWMutex:   sync.RWMutex{},
		ifce:      &cfg,
		ctx:       svcCtx,
		cancel:    cancel,
		svcLogger: logger,
		pingMap:   &utils.SyncMap[uint32, uint32]{},
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

	if err := w.initWGDevice(); err != nil {
		return errors.Join(fmt.Errorf("init WG device failed"), err)
	}

	if err := w.applyPeerConfig(); err != nil {
		return errors.Join(fmt.Errorf("apply peer config failed"), err)
	}

	if err := w.initNetwork(); err != nil {
		return errors.Join(errors.New("init network failed"), err)
	}

	if err := w.wgDevice.Up(); err != nil {
		return errors.Join(fmt.Errorf("wgDevice.Up '%s'", w.ifce.GetInterfaceName()), err)
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

	runtimeInfo.PingMap = w.pingMap.Export()
	runtimeInfo.InterfaceName = w.ifce.GetInterfaceName()

	return runtimeInfo, nil
}

func (w *wireGuard) initWGDevice() error {
	log := w.svcLogger.WithField("op", "initWGDevice")

	log.Debugf("start to create TUN device '%s' (MTU %d)", w.ifce.GetInterfaceName(), w.ifce.GetInterfaceMtu())

	var err error
	w.tunDevice, err = tun.CreateTUN(w.ifce.GetInterfaceName(), int(w.ifce.GetInterfaceMtu()))
	if err != nil {
		return errors.Join(fmt.Errorf("create TUN device '%s' (MTU %d) failed", w.ifce.GetInterfaceName(), w.ifce.GetInterfaceMtu()), err)
	}

	log.Debugf("TUN device '%s' (MTU %d) created successfully", w.ifce.GetInterfaceName(), w.ifce.GetInterfaceMtu())
	log.Debugf("start to create WireGuard device '%s'", w.ifce.GetInterfaceName())

	w.wgDevice = device.NewDevice(w.tunDevice, conn.NewDefaultBind(), &device.Logger{
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

	uapiConfigString := generateUAPIConfigString(w.ifce, w.ifce.GetParsedPrivKey(), wgTypedPeerConfigs, !w.running)

	log.Debugf("uapiBuilder: %s", uapiConfigString)

	if err = w.wgDevice.IpcSet(uapiConfigString); err != nil {
		return errors.Join(errors.New("IpcSet error"), err)
	}

	return nil
}

func (w *wireGuard) initNetwork() error {
	log := w.svcLogger.WithField("op", "initNetwork")

	link, err := netlink.LinkByName(w.ifce.GetInterfaceName())
	if err != nil {
		return errors.Join(fmt.Errorf("get iface '%s' via netlink", w.ifce.GetInterfaceName()), err)
	}

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

func (w *wireGuard) cleanupNetwork() {
	log := w.svcLogger.WithField("op", "cleanupNetwork")

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

func (w *wireGuard) pingPeers() {

	log := w.svcLogger.WithField("op", "pingPeers")

	ifceConfig, err := w.GetIfceConfig()
	if err != nil {
		log.WithError(err).Errorf("failed to get interface config")
		return
	}

	peers := ifceConfig.Peers

	var waitGroup conc.WaitGroup

	for _, peer := range peers {

		addr := ""

		if peer.Endpoint != nil && peer.Endpoint.Host != "" {
			addr = peer.Endpoint.Host
		}

		if addr == "" {
			continue
		}

		pinger, err := probing.NewPinger(addr)
		if err != nil {
			log.WithError(err).Errorf("failed to create pinger for %s", addr)
			return
		}

		pinger.Count = 5

		pinger.OnFinish = func(stats *probing.Statistics) {
			// stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss
			// stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt
			if w.pingMap != nil {
				log.Tracef("ping stats for %s: %v", addr, stats)
				avgRttMs := uint32(stats.AvgRtt.Milliseconds())
				if avgRttMs == 0 { // 0 means bug
					avgRttMs = 1
				}
				w.pingMap.Store(peer.Id, avgRttMs)
			}
		}

		pinger.OnRecv = func(pkt *probing.Packet) {
			log.Tracef("recv from %s", pkt.IPAddr.String())
		}

		waitGroup.Go(func() {
			if err := pinger.Run(); err != nil {
				log.WithError(err).Errorf("failed to run pinger for %s", addr)
				return
			}
		})
	}

	rcs := waitGroup.WaitAndRecover()
	if rcs != nil {
		log.WithError(rcs.AsError()).Errorf("failed to wait for pingers")
	}
}
