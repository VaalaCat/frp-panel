//go:build !windows

package pty

import (
	"errors"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	opty "github.com/creack/pty"
)

var _ PTYInterface = (*Pty)(nil)

var defaultShells = []string{"fish", "zsh", "bash", "sh", "mksh"}

type Pty struct {
	tty    *os.File
	cmd    *exec.Cmd
	closed bool
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
	env := envSliceToMap(os.Environ())
	ensureFishEnv(env)
	cmd.Env = envMapToSlice(env)
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
	if pty.closed {
		return nil
	}
	pty.closed = true
	if err := pty.tty.Close(); err != nil {
		return err
	}
	return pty.killChildProcess(pty.cmd)
}

func envSliceToMap(env []string) map[string]string {
	m := make(map[string]string, len(env))
	for _, kv := range env {
		k, v, ok := strings.Cut(kv, "=")
		if !ok || k == "" {
			continue
		}
		m[k] = v
	}
	return m
}

func envMapToSlice(env map[string]string) []string {
	out := make([]string, 0, len(env))
	for k, v := range env {
		if k == "" || strings.Contains(k, "=") {
			continue
		}
		out = append(out, k+"="+v)
	}
	return out
}

func ensureFishEnv(env map[string]string) {
	home := strings.TrimSpace(env["HOME"])
	if home == "" {
		// Prefer user database (e.g. /etc/passwd) when HOME is not provided.
		if u, err := user.Current(); err == nil && strings.TrimSpace(u.HomeDir) != "" {
			home = u.HomeDir
		}
	}
	if home == "" {
		// Fallback to os.UserHomeDir (may still work in some environments).
		if h, err := os.UserHomeDir(); err == nil && strings.TrimSpace(h) != "" {
			home = h
		}
	}
	if home == "" {
		uid := os.Geteuid()
		home = filepath.Join(os.TempDir(), "frp-panel-home-"+strconv.Itoa(uid))
	}
	env["HOME"] = home
	env["TERM"] = "xterm"
}
