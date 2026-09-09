//go:build !windows
// +build !windows

package wg

import (
	"errors"
	"fmt"
	"net/netip"
)

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
