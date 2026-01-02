package upgrade

import (
	"os"
	"path/filepath"
	"runtime"
)

func defaultWorkDir() string {
	// 对于 systemd 场景，优先使用持久目录，方便 upgrader service 读取 plan/status
	if runtime.GOOS == "linux" && os.Geteuid() == 0 {
		return filepath.Join(string(os.PathSeparator), "etc", "frpp", "upgrade")
	}
	return filepath.Join(os.TempDir(), "vaala-frp-panel-upgrade")
}
