//go:build windows && !arm64

package pty

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/UserExistsError/conpty"
	"github.com/iamacarpet/go-winpty"
	"github.com/shirou/gopsutil/v4/host"
)

var _ PTYInterface = (*winPTY)(nil)
var _ PTYInterface = (*conPty)(nil)

var isWin10 = IsWindows10()

type winPTY struct {
	tty    *winpty.WinPTY
	closed bool
}

type conPty struct {
	tty    *conpty.ConPty
	closed bool
}

func IsWindows10() bool {
	hi, err := host.Info()
	if err != nil {
		return false
	}

	re := regexp.MustCompile(`Build (\d+(\.\d+)?)`)
	match := re.FindStringSubmatch(hi.KernelVersion)
	if len(match) > 1 {
		versionStr := match[1]

		version, err := strconv.ParseFloat(versionStr, 64)
		if err != nil {
			return false
		}

		return version >= 17763
	}
	return false
}

func getExecutableFilePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}

func Start() (PTYInterface, error) {
	shellPath, err := exec.LookPath("powershell.exe")
	if err != nil || shellPath == "" {
		shellPath = "cmd.exe"
	}
	path, err := getExecutableFilePath()
	if err != nil {
		return nil, err
	}
	if !isWin10 {
		tty, err := winpty.OpenDefault(path, shellPath)
		return &winPTY{tty: tty}, err
	}
	tty, err := conpty.Start(shellPath, conpty.ConPtyWorkDir(path))
	return &conPty{tty: tty}, err
}

func (w *winPTY) Write(p []byte) (n int, err error) {
	return w.tty.StdIn.Write(p)
}

func (w *winPTY) Read(p []byte) (n int, err error) {
	return w.tty.StdOut.Read(p)
}

func (w *winPTY) Getsize() (uint16, uint16, error) {
	return 80, 40, nil
}

func (w *winPTY) Setsize(cols, rows uint32) error {
	w.tty.SetSize(cols, rows)
	return nil
}

func (w *winPTY) Close() error {
	if w.closed {
		return nil
	}
	w.closed = true
	w.tty.Close()
	return nil
}

func (c *conPty) Write(p []byte) (n int, err error) {
	return c.tty.Write(p)
}

func (c *conPty) Read(p []byte) (n int, err error) {
	return c.tty.Read(p)
}

func (c *conPty) Getsize() (uint16, uint16, error) {
	return 80, 40, nil
}

func (c *conPty) Setsize(cols, rows uint32) error {
	c.tty.Resize(int(cols), int(rows))
	return nil
}

func (c *conPty) Close() error {
	if c.closed {
		return nil
	}
	c.closed = true
	if err := c.tty.Close(); err != nil {
		return err
	}
	return nil
}
