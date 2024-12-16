package utils

import (
	"testing"

	v1 "github.com/fatedier/frp/pkg/config/v1"
)

func TestLoadConfigureFromContent(t *testing.T) {
	content := []byte(`[[proxies]]
name = "ssh"
type = "tcp"
localIP = "127.0.0.1"
localPort = 22
remotePort = 6000`)

	allCfg := &v1.ClientConfig{}

	if err := LoadConfigureFromContent(content, allCfg, true); err != nil {
		t.Error(err)
	}
	t.Errorf("%+v", allCfg)
}
