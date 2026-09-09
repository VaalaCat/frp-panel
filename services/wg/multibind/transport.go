package multibind

import (
	"golang.zx2c4.com/wireguard/conn"
)

type Transport struct {
	bind conn.Bind
	name string
}

func (t *Transport) loadOrNewEndpoint(inner conn.Endpoint) conn.Endpoint {
	return &MultiEndpoint{
		trans: t,
		inner: inner,
	}
}

func NewTransport(bind conn.Bind, name string) *Transport {
	return &Transport{
		bind: bind,
		name: name,
	}
}
