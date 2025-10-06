package wg_test

import (
	"testing"

	"github.com/VaalaCat/frp-panel/services/wg"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func TestGenerateKeys(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want wg.WGKeys
	}{
		{"test", wg.WGKeys{
			PrivateKeyBase64: "test",
			PublicKeyBase64:  "test",
			PrivateKey:       wgtypes.Key{},
			PublicKey:        wgtypes.Key{},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wg.GenerateKeys()
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("GenerateKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
