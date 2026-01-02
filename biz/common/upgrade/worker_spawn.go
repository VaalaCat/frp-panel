package upgrade

import (
	"fmt"
	"os"
	"path/filepath"
)

// spawnWorker 使用 os.StartProcess 拉起独立 worker（不使用 exec 包）
func spawnWorker(exePath string, planPath string) error {
	if len(exePath) == 0 || len(planPath) == 0 {
		return fmt.Errorf("exePath/planPath is empty")
	}

	argv := []string{exePath, "__upgrade-worker", "--plan", planPath}
	devnull, _ := os.OpenFile(os.DevNull, os.O_RDWR, 0)
	files := []*os.File{os.Stdin, os.Stdout, os.Stderr}
	if devnull != nil {
		files = []*os.File{devnull, devnull, devnull}
	}
	attr := &os.ProcAttr{
		Dir:   filepath.Dir(exePath),
		Env:   os.Environ(),
		Files: files,
	}
	applyDetachAttr(attr)

	p, err := os.StartProcess(exePath, argv, attr)
	if err != nil {
		return err
	}
	_ = p.Release()
	if devnull != nil {
		_ = devnull.Close()
	}
	return nil
}
