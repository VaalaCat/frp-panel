package models_test

import (
	"net/netip"
	"testing"

	"github.com/VaalaCat/frp-panel/models"
)

func TestParseIPOrCIDRWithNetip(t *testing.T) {
	ip, cidr, _ := models.ParseIPOrCIDRWithNetip("192.168.1.1/24")
	t.Errorf("ip: %v, cidr: %v", ip, cidr)
	newcidr := netip.PrefixFrom(ip, 32)
	t.Errorf("newcidr: %v", newcidr)
}
