//go:build !windows
// +build !windows

package wg

import (
	"errors"
	"net/netip"
	"reflect"
	"unsafe"

	"gvisor.dev/gvisor/pkg/tcpip"
	"gvisor.dev/gvisor/pkg/tcpip/network/ipv4"
	"gvisor.dev/gvisor/pkg/tcpip/network/ipv6"
	"gvisor.dev/gvisor/pkg/tcpip/stack"
)

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
