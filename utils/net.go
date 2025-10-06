package utils

import (
	"fmt"
	"math/big"
	"net"
	"strings"
)

// GetLocalIPv4s 返回本地所有活跃网络接口的 IPv4 地址列表。
// 忽略未启用、回环以及非 IPv4 地址。
func GetLocalIPv4s() ([]net.IP, error) {
	var ips []net.IP
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("list interfaces failed: %w", err)
	}

	for _, iface := range ifaces {
		// 跳过未启用或回环接口
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			// 某些接口可能无权限，此处跳过
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// 仅保留 IPv4，过滤回环
			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}

			ipstr := ip.String()
			if ipstr == "<nil>" {
				continue
			}

			ips = append(ips, ip)
		}
	}

	return ips, nil
}

// nextIP returns ip + 1
func nextIP(ip net.IP) net.IP {
	ip = ip.To4()
	ipInt := big.NewInt(0).SetBytes(ip)
	ipInt.Add(ipInt, big.NewInt(1))
	return net.IP(ipInt.FillBytes(make([]byte, 4)))
}

// ipLE returns true if ip1 ≤ ip2
func ipLE(ip1, ip2 net.IP) bool {
	a := big.NewInt(0).SetBytes(ip1.To4())
	b := big.NewInt(0).SetBytes(ip2.To4())
	return a.Cmp(b) <= 0
}

// parseNetwork parses CIDR and returns network IP and broadcast IP
func parseNetwork(cidr string) (net.IP, net.IP, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid network CIDR %q: %w", cidr, err)
	}
	network := ip.Mask(ipnet.Mask)
	ones, bits := ipnet.Mask.Size()
	total := big.NewInt(1)
	total.Lsh(total, uint(bits-ones))
	// broadcast = network + total - 1
	bcastInt := big.NewInt(0).SetBytes(network.To4())
	bcastInt.Add(bcastInt, big.NewInt(0).Sub(total, big.NewInt(1)))
	bcast := net.IP(bcastInt.FillBytes(make([]byte, 4)))
	return network, bcast, nil
}

// buildUsedSet builds a set of used IP strings.
// If an entry is CIDR, only its base IP is used.
func buildUsedSet(allocated []string) (map[string]bool, error) {
	used := make(map[string]bool)
	for _, a := range allocated {
		ipStr := a
		if idx := strings.Index(a, "/"); idx != -1 {
			ipStr = a[:idx]
		}
		ip := net.ParseIP(ipStr)
		if ip == nil {
			return nil, fmt.Errorf("invalid allocated IP %q", a)
		}
		used[ip.String()] = true
	}
	return used, nil
}

// IsIPFree checks whether ipStr is free in the network given allocated list.
// allocated entries with CIDR ignore the mask (only base IP).
func IsIPFree(networkCidr string, allocated []string, ipStr string) (bool, error) {
	network, bcast, err := parseNetwork(networkCidr)
	if err != nil {
		return false, nil
	}
	ip := net.ParseIP(ipStr).To4()
	if ip == nil {
		return false, fmt.Errorf("invalid IP %q", ipStr)
	}
	// must be inside (network, broadcast)
	if !ipLE(nextIP(network), ip) || !ipLE(ip, nextIP(bcast)) {
		return false, nil
	}
	used, err := buildUsedSet(allocated)
	if err != nil {
		return false, err
	}
	if used[ip.String()] {
		return false, nil
	}
	return true, nil
}

// AllocateIP tries to allocate desiredStr; if not free or empty, picks the first free IP in network.
func AllocateIP(networkCidr string, allocated []string, desiredStr string) (string, error) {
	network, bcast, err := parseNetwork(networkCidr)
	if err != nil {
		return "", err
	}
	used, err := buildUsedSet(allocated)
	if err != nil {
		return "", err
	}
	// try desired
	if desiredStr != "" {
		free, err := IsIPFree(networkCidr, allocated, desiredStr)
		if err != nil {
			return "", err
		}
		if free {
			return desiredStr, nil
		}
	}
	// scan for first free
	for ip := nextIP(network); ipLE(ip, nextIP(bcast)); ip = nextIP(ip) {
		if ip.Equal(bcast) {
			break
		}
		s := ip.String()
		if !used[s] {
			return s, nil
		}
	}
	return "", fmt.Errorf("no available IP")
}
