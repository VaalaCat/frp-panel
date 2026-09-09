//go:build windows
// +build windows

package wg

import (
	"errors"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/services/app"
)

var errWireGuardNotSupported = errors.New("wireguard is not supported on windows build")

type wireGuardManagerWindows struct{}

func NewWireGuardManager(appInstance app.Application) app.WireGuardManager {
	return &wireGuardManagerWindows{}
}

func (m *wireGuardManagerWindows) CreateService(cfg *defs.WireGuardConfig) (app.WireGuard, error) {
	return nil, errWireGuardNotSupported
}

func (m *wireGuardManagerWindows) StartService(interfaceName string) error {
	return errWireGuardNotSupported
}

func (m *wireGuardManagerWindows) GetService(interfaceName string) (app.WireGuard, bool) {
	return nil, false
}

func (m *wireGuardManagerWindows) StopService(interfaceName string) error {
	return errWireGuardNotSupported
}

func (m *wireGuardManagerWindows) GetAllServices() []app.WireGuard {
	return nil
}

func (m *wireGuardManagerWindows) RemoveService(interfaceName string) error {
	return errWireGuardNotSupported
}

func (m *wireGuardManagerWindows) StopAllServices() map[string]error {
	return map[string]error{}
}

func (m *wireGuardManagerWindows) RestartService(interfaceName string) error {
	return errWireGuardNotSupported
}

func (m *wireGuardManagerWindows) Start() {}

func (m *wireGuardManagerWindows) Stop() {}
