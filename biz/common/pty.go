package common

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/rpcclient"
	"github.com/creack/pty"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc"
)

var (
	LinuxShells = []string{
		"/usr/bin/fish",
		"/usr/local/bin/fish",
		"/bin/fish",
		"/opt/homebrew/bin/fish",
		"/usr/bin/zsh",
		"/bin/zsh",
		"/bin/bash",
		"/bin/sh",
	}
	WindowsShells = []string{"cmd.exe", "powershell.exe"}
)

func StartPTYConnect(c context.Context, req *pb.CommonRequest, initMsg *pb.PTYClientMessage) (*pb.CommonResponse, error) {
	conn, err := rpcclient.GetClientRPCSerivce().GetCli().PTYConnect(c)
	if err != nil {
		logger.Logger(c).WithError(err).Infof("rpc connect master error")
		return nil, err
	}

	sessionID := uuid.New().String()
	initMsg.SessionId = sessionID

	if err := conn.Send(initMsg); err != nil {
		logger.Logger(c).WithError(err).Infof("send server base error")
		return nil, err
	}

	ack, err := conn.Recv()
	if err != nil {
		logger.Logger(c).WithError(err).Infof("recv ack error")
		return nil, err
	}

	if ack.GetData() != "ok" {
		logger.Logger(c).Infof("ack error")
		return nil, fmt.Errorf("ack error")
	}

	go func() {
		HandlePTY(c, conn, sessionID)
	}()

	return &pb.CommonResponse{Data: &sessionID}, nil
}

func HandlePTY(c context.Context, conn pb.Master_PTYConnectClient, sessionID string) {
	connectionErrorLimit := 10
	maxBufferSizeBytes := 4096

	cmd := exec.Command(GetShell())
	cmd.Env = os.Environ()
	tty, err := pty.Start(cmd)
	if err != nil {
		msg := fmt.Sprintf("failed to start tty: %s", err)
		logger.Logger(c).WithError(err).Warn(msg)
		conn.Send(&pb.PTYClientMessage{Data: &msg, SessionId: sessionID})
		return
	}

	defer func() {
		logger.Logger(c).Info("gracefully stopping spawned tty...")
		if err := cmd.Process.Kill(); err != nil {
			logger.Logger(c).Warnf("failed to kill process: %s", err)
		}
		if _, err := cmd.Process.Wait(); err != nil {
			logger.Logger(c).Warnf("failed to wait for process to exit: %s", err)
		}
		if err := tty.Close(); err != nil {
			logger.Logger(c).Warnf("failed to close spawned tty gracefully: %s", err)
		}
		if err := conn.CloseSend(); err != nil {
			logger.Logger(c).Warnf("failed to close webscoket connection: %s", err)
		}
	}()

	var connectionClosed bool
	var wg conc.WaitGroup

	// tty >> master
	wg.Go(func() {
		errorCounter := 0
		for {
			if errorCounter > connectionErrorLimit {
				break
			}
			buffer := make([]byte, maxBufferSizeBytes)
			readLength, err := tty.Read(buffer)
			if err != nil {
				logger.Logger(c).Warnf("failed to read from tty: %s", err)
				if err := conn.Send(&pb.PTYClientMessage{Data: lo.ToPtr("bye!"), SessionId: sessionID}); err != nil {
					logger.Logger(c).Warnf("failed to send termination message from tty to master: %s", err)
				}
				if err := conn.CloseSend(); err != nil {
					logger.Logger(c).Warnf("failed to close grpc stream connection: %s", err)
				}
				return
			}
			str := string(buffer[:readLength])
			if err := conn.Send(&pb.PTYClientMessage{Data: lo.ToPtr(str), SessionId: sessionID}); err != nil {
				logger.Logger(c).Warnf("failed to send %v bytes from tty to master", readLength)
				errorCounter++
				continue
			}
			logger.Logger(c).Tracef("sent message of size %v bytes from tty to master", readLength)
			errorCounter = 0
		}
	})

	// tty << master
	wg.Go(func() {
		for {
			// data processing
			msg, err := conn.Recv()
			if err != nil {
				if !connectionClosed {
					logger.Logger(c).Warnf("failed to get next reader: %s", err)
				}
				if err := conn.CloseSend(); err != nil {
					logger.Logger(c).Warnf("failed to close grpc stream connection: %s", err)
				}
				return
			}
			if msg.GetDone() {
				if err := conn.CloseSend(); err != nil {
					logger.Logger(c).Warnf("failed to close grpc stream connection: %s", err)
				}
				logger.Logger(c).Info("gracefully stopping spawned tty...")
				if err := cmd.Process.Kill(); err != nil {
					logger.Logger(c).Warnf("failed to kill process: %s", err)
				}
				if _, err := cmd.Process.Wait(); err != nil {
					logger.Logger(c).Warnf("failed to wait for process to exit: %s", err)
				}
				if err := tty.Close(); err != nil {
					logger.Logger(c).Warnf("failed to close spawned tty gracefully: %s", err)
				}
				logger.Logger(c).Info("recv server msg done, closing conn...")
				return
			}
			data := msg.GetData()

			// handle resizing
			if msg.Height != nil && msg.Width != nil {
				logger.Logger(c).Infof("resizing tty to use %+v rows and %+v columns...", *msg.Height, *msg.Width)
				if err := pty.Setsize(tty, &pty.Winsize{
					Rows: uint16(*msg.Height),
					Cols: uint16(*msg.Width),
				}); err != nil {
					logger.Logger(c).Warnf("failed to resize tty, error: %s", err)
				}
				continue
			}

			// write to tty
			bytesWritten, err := tty.Write([]byte(data))
			if err != nil {
				logger.Logger(c).Warn(fmt.Sprintf("failed to write %v bytes to tty: %s", len(data), err))
				continue
			}
			logger.Logger(c).Tracef("%v bytes written to tty...", bytesWritten)
		}
	})

	wg.Wait()
	logger.Logger(c).Info("closing conn...")
	connectionClosed = true
}

func GetShell() string {
	if runtime.GOOS == "windows" {
		for sh := range WindowsShells {
			if _, err := exec.LookPath(WindowsShells[sh]); err == nil {
				return WindowsShells[sh]
			}
		}
		return WindowsShells[len(WindowsShells)-1]
	}
	for sh := range LinuxShells {
		if _, err := exec.LookPath(LinuxShells[sh]); err == nil {
			return LinuxShells[sh]
		}
	}
	return LinuxShells[len(LinuxShells)-1]
}
