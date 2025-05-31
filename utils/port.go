package utils

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// GetFreePort asks the kernel for a free port for the given network.
// Valid networks: "tcp4", "tcp6", "udp4", "udp6".
func GetFreePort(network string) (uint32, error) {
	network = strings.ToLower(network)
	var (
		port int
		err  error
	)

	addr := ":0" // let OS choose
	switch network {
	case "tcp4", "tcp6":
		var ln net.Listener
		ln, err = net.Listen(network, addr)
		if err != nil {
			return 0, fmt.Errorf("listen %s failed: %w", network, err)
		}
		defer ln.Close()
		port, err = extractPort(ln.Addr().String())

	case "udp4", "udp6":
		var pc net.PacketConn
		pc, err = net.ListenPacket(network, addr)
		if err != nil {
			return 0, fmt.Errorf("listenpacket %s failed: %w", network, err)
		}
		defer pc.Close()
		port, err = extractPort(pc.LocalAddr().String())

	default:
		return 0, fmt.Errorf("unsupported network %q", network)
	}

	if err != nil {
		return 0, err
	}
	return uint32(port), nil
}

// extractPort splits "host:port" and returns port as int.
func extractPort(hostport string) (int, error) {
	_, portStr, err := net.SplitHostPort(hostport)
	if err != nil {
		return 0, fmt.Errorf("split hostport %q: %w", hostport, err)
	}
	p, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("invalid port %q: %w", portStr, err)
	}
	return p, nil
}
