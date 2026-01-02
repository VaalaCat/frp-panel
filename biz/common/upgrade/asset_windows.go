//go:build windows

package upgrade

import (
	"fmt"
	"runtime"
)

func detectAssetName() (string, error) {
	switch runtime.GOARCH {
	case "amd64":
		return "frp-panel-windows-amd64.exe", nil
	case "arm64":
		return "frp-panel-windows-arm64.exe", nil
	default:
		return "", fmt.Errorf("暂不支持的系统/架构: %s %s", runtime.GOOS, runtime.GOARCH)
	}
}


