package multibind

import (
	"net/netip"

	"golang.zx2c4.com/wireguard/conn"
)

var (
	_ conn.Endpoint = (*MultiEndpoint)(nil)
)

type MultiEndpoint struct {
	trans *Transport
	inner conn.Endpoint
}

func NewMultiEndpoint(trans *Transport, inner conn.Endpoint) *MultiEndpoint {
	return &MultiEndpoint{
		trans: trans,
		inner: inner,
	}
}

// ClearSrc implements conn.Endpoint.
func (m *MultiEndpoint) ClearSrc() {
	m.inner.ClearSrc()
}

// DstIP implements conn.Endpoint.
func (m *MultiEndpoint) DstIP() netip.Addr {
	return m.inner.DstIP()
}

// DstToBytes implements conn.Endpoint.
func (m *MultiEndpoint) DstToBytes() []byte {
	return m.inner.DstToBytes()
}

// DstToString implements conn.Endpoint.
func (m *MultiEndpoint) DstToString() string {
	return m.inner.DstToString()
}

// SrcIP implements conn.Endpoint.
func (m *MultiEndpoint) SrcIP() netip.Addr {
	return m.inner.SrcIP()
}

// SrcToString implements conn.Endpoint.
func (m *MultiEndpoint) SrcToString() string {
	return m.inner.SrcToString()
}
