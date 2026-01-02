//go:build windows

package upgrade

import (
	"fmt"
	"os"
	"path/filepath"
)

type unlockFn func()

// Windows 简化处理：用 O_EXCL 创建锁文件
// 说明：Windows 不支持 unix flock；进程异常退出可能残留 lock 文件，必要时可手动删除。
func lock(workDir string) (unlockFn, error) {
	if len(workDir) == 0 {
		workDir = defaultWorkDir()
	}
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return nil, fmt.Errorf("create work dir failed: %w", err)
	}

	lockPath := filepath.Join(workDir, "upgrade.lock")
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("upgrade already in progress (lock exists): %w", err)
	}
	_ = f.Close()

	return func() { _ = os.Remove(lockPath) }, nil
}


