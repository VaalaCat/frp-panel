//go:build !windows

package pty

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

type shellCandidate struct {
	path string
	args []string
}

func Start() (PTYInterface, error) {
	var collectedErr error

	seen := map[string]struct{}{}
	addCandidate := func(path string, args ...string) []shellCandidate {
		if path == "" {
			return nil
		}
		key := filepath.Join(path, strings.Join(args, "\x00"))
		if _, ok := seen[key]; ok {
			return nil
		}
		seen[key] = struct{}{}
		return []shellCandidate{{path: path, args: args}}
	}

	var candidates []shellCandidate

	if shell := os.Getenv("SHELL"); shell != "" {
		if !strings.Contains(shell, string(os.PathSeparator)) {
			if resolved, err := exec.LookPath(shell); err == nil {
				shell = resolved
			} else {
				collectedErr = errors.Join(collectedErr, err)
			}
		}
		candidates = append(candidates, addCandidate(shell)...)
	}

	for _, shell := range defaultShells {
		candidate := shell
		if !strings.Contains(shell, string(os.PathSeparator)) {
			if resolved, err := exec.LookPath(shell); err == nil {
				candidate = resolved
			} else {
				collectedErr = errors.Join(collectedErr, err)
				continue
			}
		}
		candidates = append(candidates, addCandidate(candidate)...)
	}

	// Android fallbacks
	if _, err := os.Stat("/system/bin/sh"); err == nil {
		candidates = append(candidates, addCandidate("/system/bin/sh")...)
	} else {
		collectedErr = errors.Join(collectedErr, err)
	}
	if _, err := os.Stat("/system/bin/toybox"); err == nil {
		candidates = append(candidates, addCandidate("/system/bin/toybox", "sh")...)
	} else {
		collectedErr = errors.Join(collectedErr, err)
	}

	for _, candidate := range candidates {
		pty, err := startPTY(candidate.path, candidate.args...)
		if err != nil {
			collectedErr = errors.Join(collectedErr, err)
			continue
		}
		return pty, nil
	}

	if collectedErr != nil {
		return nil, collectedErr
	}

	return nil, errors.New("none of the default shells was found")
}

func startPTY(path string, args ...string) (PTYInterface, error) {
	cmd := exec.Command(path, args...)
	cmd.Env = append(os.Environ(), "TERM=xterm")

	tty, err := opty.Start(cmd)
	if err != nil {
		commandLine := path
		if len(args) > 0 {
			commandLine += " " + strings.Join(args, " ")
		}
		return nil, fmt.Errorf("%s: %w", commandLine, err)
	}

	return &Pty{tty: tty, cmd: cmd}, nil
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
