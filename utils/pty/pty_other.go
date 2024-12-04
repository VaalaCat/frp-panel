//go:build !windows

package pty

import (
	"errors"
	"os"
	"os/exec"
	"syscall"

	opty "github.com/creack/pty"
)

var _ PTYInterface = (*Pty)(nil)

var defaultShells = []string{"fish", "zsh", "bash", "sh"}

type Pty struct {
	tty *os.File
	cmd *exec.Cmd
}

func DownloadDependency() error {
	return nil
}

func Start() (PTYInterface, error) {
	var shellPath string
	for i := 0; i < len(defaultShells); i++ {
		shellPath, _ = exec.LookPath(defaultShells[i])
		if shellPath != "" {
			break
		}
	}
	if shellPath == "" {
		return nil, errors.New("none of the default shells was found")
	}
	cmd := exec.Command(shellPath)
	cmd.Env = append(os.Environ(), "TERM=xterm")
	tty, err := opty.Start(cmd)
	return &Pty{tty: tty, cmd: cmd}, err
}

func (pty *Pty) Write(p []byte) (n int, err error) {
	return pty.tty.Write(p)
}

func (pty *Pty) Read(p []byte) (n int, err error) {
	return pty.tty.Read(p)
}

func (pty *Pty) Getsize() (uint16, uint16, error) {
	ws, err := opty.GetsizeFull(pty.tty)
	if err != nil {
		return 0, 0, err
	}
	return ws.Cols, ws.Rows, nil
}

func (pty *Pty) Setsize(cols, rows uint32) error {
	return opty.Setsize(pty.tty, &opty.Winsize{
		Cols: uint16(cols),
		Rows: uint16(rows),
	})
}

func (pty *Pty) killChildProcess(c *exec.Cmd) error {
	pgid, err := syscall.Getpgid(c.Process.Pid)
	if err != nil {
		// Fall-back on error. Kill the main process only.
		c.Process.Kill()
	}
	// Kill the whole process group.
	syscall.Kill(-pgid, syscall.SIGKILL) // SIGKILL 直接杀掉 SIGTERM 发送信号，等待进程自己退出
	return c.Wait()
}

func (pty *Pty) Close() error {
	if err := pty.tty.Close(); err != nil {
		return err
	}
	return pty.killChildProcess(pty.cmd)
}
