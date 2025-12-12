package wg

import (
	"errors"
	"fmt"
	"net/netip"
	"sync"

	"github.com/coreos/go-iptables/iptables"
	"github.com/sirupsen/logrus"
)

// firewallManager 负责为 WireGuard 接口配置必要的转发表规则。
// 目前通过 go-iptables 调用 iptables/ip6tables，互斥串行更新。
// 规则策略：允许接口 -> 网段、网段 -> 接口 的转发。
type firewallManager struct {
	mu      sync.Mutex
	logger  *logrus.Entry
	tracked map[string]string // iface -> cidr
}

func newFirewallManager(logger *logrus.Entry) *firewallManager {
	return &firewallManager{
		logger:  logger,
		tracked: make(map[string]string),
	}
}

// ApplyRelayRules 确保接口与网段的转发规则存在： -i iface -d cidr ACCEPT 与 -s cidr -o iface ACCEPT。
// cidr 形如 "10.10.0.0/24"。
func (f *firewallManager) ApplyRelayRules(iface, cidr string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if cidr == "" {
		return errors.New("cidr is empty")
	}

	_, prefixErr := netip.ParsePrefix(cidr)
	if prefixErr != nil {
		return errors.Join(fmt.Errorf("parse cidr '%s' failed", cidr), prefixErr)
	}

	// 如果 CIDR 未变化则跳过
	if old, ok := f.tracked[iface]; ok && old == cidr {
		return nil
	}

	if err := f.ensureRelayRules(iface, cidr); err != nil {
		return err
	}

	f.tracked[iface] = cidr
	return nil
}

// Cleanup 会删除指定接口下记录过的所有规则。
func (f *firewallManager) Cleanup(iface string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	cidr, ok := f.tracked[iface]
	if !ok {
		return nil
	}

	err := f.deleteRelayRules(iface, cidr)
	delete(f.tracked, iface)
	return err
}

func (f *firewallManager) ensureRelayRules(iface, cidr string) error {
	prefix, prefixErr := netip.ParsePrefix(cidr)
	if prefixErr != nil {
		return errors.Join(fmt.Errorf("parse cidr '%s' failed", cidr), prefixErr)
	}

	ipt, err := newIPT(prefix.Addr().Is6())
	if err != nil {
		return err
	}

	if err := ipt.AppendUnique("filter", "FORWARD", "-i", iface, "-d", cidr, "-j", "ACCEPT"); err != nil {
		return errors.Join(fmt.Errorf("append forward rule iface->cidr failed"), err)
	}

	if err := ipt.AppendUnique("filter", "FORWARD", "-s", cidr, "-o", iface, "-j", "ACCEPT"); err != nil {
		return errors.Join(fmt.Errorf("append forward rule cidr->iface failed"), err)
	}

	return nil
}

func (f *firewallManager) deleteRelayRules(iface, cidr string) error {
	prefix, prefixErr := netip.ParsePrefix(cidr)
	if prefixErr != nil {
		return errors.Join(fmt.Errorf("parse cidr '%s' failed", cidr), prefixErr)
	}

	ipt, err := newIPT(prefix.Addr().Is6())
	if err != nil {
		return err
	}

	var errs error
	if err := ipt.DeleteIfExists("filter", "FORWARD", "-i", iface, "-d", cidr, "-j", "ACCEPT"); err != nil {
		errs = errors.Join(errs, errors.Join(fmt.Errorf("delete rule iface->cidr failed"), err))
	}
	if err := ipt.DeleteIfExists("filter", "FORWARD", "-s", cidr, "-o", iface, "-j", "ACCEPT"); err != nil {
		errs = errors.Join(errs, errors.Join(fmt.Errorf("delete rule cidr->iface failed"), err))
	}
	return errs
}

func newIPT(isIPv6 bool) (*iptables.IPTables, error) {
	proto := iptables.ProtocolIPv4
	if isIPv6 {
		proto = iptables.ProtocolIPv6
	}
	ipt, err := iptables.NewWithProtocol(proto)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("create iptables handler (ipv6=%v)", isIPv6), err)
	}
	return ipt, nil
}
