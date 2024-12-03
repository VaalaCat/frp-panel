//go:build windows && arm64

package pty

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/UserExistsError/conpty"
)

var _ PTYInterface = (*Pty)(nil)

type Pty struct {
	tty *conpty.ConPty
}

func DownloadDependency() error {
	return nil
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
	tty, err := conpty.Start(shellPath, conpty.ConPtyWorkDir(path))
	return &Pty{tty: tty}, err
}

func (pty *Pty) Write(p []byte) (n int, err error) {
	return pty.tty.Write(p)
}

func (pty *Pty) Read(p []byte) (n int, err error) {
	return pty.tty.Read(p)
}

func (pty *Pty) Getsize() (uint16, uint16, error) {
	return 80, 40, nil
}

func (pty *Pty) Setsize(cols, rows uint32) error {
	return pty.tty.Resize(int(cols), int(rows))
}

func (pty *Pty) Close() error {
	if err := pty.tty.Close(); err != nil {
		return err
	}
	return nil
}
