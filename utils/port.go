package utils

import (
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

func GetAvailablePort(addr string) (int, error) {
	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", addr))
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		return 0, err
	}

	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

func IsPortAvailable(port int, addr string) bool {

	address := fmt.Sprintf("%s:%d", addr, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logrus.Infof("port %s is taken: %s", address, err)
		return false
	}

	defer listener.Close()
	return true
}

func WaitForPort(host string, port int) {
	for {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
		if err == nil {
			conn.Close()
			break
		}

		logrus.Warnf("Target port %s:%d is not open yet, waiting...\n", host, port)
		time.Sleep(time.Second * 5)
	}
	logrus.Infof("Target port %s:%d is open", host, port)
	time.Sleep(time.Second * 1)
}
