//go:build !windows
// +build !windows

package wg

import (
	"errors"
	"fmt"
	"net/netip"

	"github.com/samber/lo"

	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun"
	"golang.zx2c4.com/wireguard/tun/netstack"
)

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
