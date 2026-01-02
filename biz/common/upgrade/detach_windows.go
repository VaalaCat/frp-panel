//go:build windows

package upgrade

import (
	"os"
	"syscall"
)

func applyDetachAttr(attr *os.ProcAttr) {
	if attr == nil {
		return
	}
	attr.Sys = &syscall.SysProcAttr{
		// 0x00000008: DETACHED_PROCESS（避免引入额外 windows 依赖，同时满足“不使用 exec 包”）
		CreationFlags: 0x00000008,
	}
}
