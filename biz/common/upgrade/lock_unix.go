//go:build !windows

package upgrade

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

type unlockFn func()

// Unix 使用 flock：进程异常退出/被 SIGINT 杀死时锁会自动释放，不会残留“死锁文件”
func lock(workDir string) (unlockFn, error) {
	if len(workDir) == 0 {
		workDir = defaultWorkDir()
	}
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return nil, fmt.Errorf("create work dir failed: %w", err)
	}

	lockPath := filepath.Join(workDir, "upgrade.lock")
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, fmt.Errorf("open lock file failed: %w", err)
	}

	// LOCK_NB：如果已有升级在跑，直接返回错误（不会阻塞）
	if err := unix.Flock(int(f.Fd()), unix.LOCK_EX|unix.LOCK_NB); err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("upgrade already in progress: %w", err)
	}

	return func() {
		_ = unix.Flock(int(f.Fd()), unix.LOCK_UN)
		_ = f.Close()
	}, nil
}


