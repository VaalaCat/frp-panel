package models_test

import (
	"net/netip"
	"testing"

	"github.com/VaalaCat/frp-panel/models"
	"github.com/stretchr/testify/assert"
)

func TestParseIPOrCIDRWithNetip(t *testing.T) {
	ip, cidr, _ := models.ParseIPOrCIDRWithNetip("192.168.1.1/24")
	t.Logf("ip: %v, cidr: %v", ip, cidr)
	assert.Equal(t, ip, netip.MustParseAddr("192.168.1.1"))
	assert.Equal(t, cidr, netip.MustParsePrefix("192.168.1.1/24"))
	newcidr := netip.PrefixFrom(ip, 32)
	assert.Equal(t, newcidr, netip.MustParsePrefix("192.168.1.1/32"))
	t.Logf("newcidr: %v", newcidr)
}
