package utils

import (
	"fmt"
	"net"
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
