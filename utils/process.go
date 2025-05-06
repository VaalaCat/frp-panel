package utils

import (
	"os"
	"strings"

	"github.com/shirou/gopsutil/v4/process"
)

func ProcessExistsBySelf(target string) (bool, error) {
	selfPID := int32(os.Getpid())

	procs, err := process.Processes()
	if err != nil {
		return false, err
	}

	for _, p := range procs {
		ppid, err := p.Ppid()
		if err != nil || ppid != selfPID {
			continue
		}
		cmdline, err := p.Cmdline()
		if err != nil {
			continue
		}
		if strings.Contains(cmdline, target) {
			return true, nil
		}
	}
	return false, nil
}
