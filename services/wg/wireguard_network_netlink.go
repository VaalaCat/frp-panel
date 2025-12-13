//go:build !windows
// +build !windows

package wg

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/vishvananda/netlink"
)

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
