//go:build !windows

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
		Setsid: true,
	}
}
