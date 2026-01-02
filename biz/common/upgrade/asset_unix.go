//go:build !windows

package upgrade

import (
	"fmt"
	"runtime"
)

func detectAssetName() (string, error) {
	osName := runtime.GOOS
	machine := unameMachine()
	if len(machine) == 0 {
		// fallback
		machine = runtime.GOARCH
	}

	switch osName {
	case "linux":
		switch machine {
		case "x86_64", "amd64":
			return "frp-panel-linux-amd64", nil
		case "aarch64", "arm64":
			return "frp-panel-linux-arm64", nil
		case "armv7l":
			return "frp-panel-linux-armv7l", nil
		case "armv6l":
			return "frp-panel-linux-armv6l", nil
		}
	case "darwin":
		switch machine {
		case "x86_64", "amd64":
			return "frp-panel-darwin-amd64", nil
		case "arm64":
			return "frp-panel-darwin-arm64", nil
		}
	}
	return "", fmt.Errorf("暂不支持的系统/架构: %s %s", osName, machine)
}


