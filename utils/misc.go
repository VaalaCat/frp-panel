package utils

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/jackpal/gateway"
)

func IsSameDay(first time.Time, second time.Time) bool {
	return first.YearDay() == second.YearDay() && first.Year() == second.Year()
}

func GetHostnameWithIP() string {
	hostname, _ := os.Hostname()
	interfaces, err := net.Interfaces()
	if err != nil {
		return hostname
	}
	ipGateway, err := gateway.DiscoverGateway()
	if err != nil {
		return hostname
	}

	stop := false
	for _, iface := range interfaces {
		if stop {
			break
		}

		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			addrIP, ok := addr.(*net.IPNet)
			if !ok || addrIP == nil {
				continue
			}
			if !addrIP.Contains(ipGateway) {
				continue
			}
			hostname = fmt.Sprintf("%s-%s", hostname, addrIP.IP.String())
			stop = true
			break
		}
	}
	return hostname
}
