package utils_test

import (
	"testing"

	"github.com/VaalaCat/frp-panel/utils"
)

func TestAllocateIP(t *testing.T) {
	ip, err := utils.AllocateIP("192.168.1.0/24", []string{"192.168.1.1/24"}, "192.168.1.1")
	if err != nil {
		t.Errorf("AllocateIP() failed: %v", err)
	}
	t.Errorf("AllocateIP() = %v", ip)
}
