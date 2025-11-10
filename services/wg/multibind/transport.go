package multibind

import (
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/samber/lo"
	"golang.zx2c4.com/wireguard/conn"
)

type Transport struct {
	bind      conn.Bind
	name      string
	endpoints utils.SyncMap[conn.Endpoint, *MultiEndpoint]
}

func (t *Transport) loadOrNewEndpoint(inner conn.Endpoint) conn.Endpoint {
	if lo.IsNil(inner) {
		return &MultiEndpoint{
			trans: t,
			inner: inner,
		}
	}
	if cached, ok := t.endpoints.Load(inner); ok {
		return cached
	}

	newEndpoint := &MultiEndpoint{
		trans: t,
		inner: inner,
	}
	t.endpoints.Store(inner, newEndpoint)

	return newEndpoint
}

func NewTransport(bind conn.Bind, name string) *Transport {
	return &Transport{
		bind:      bind,
		name:      name,
		endpoints: utils.SyncMap[conn.Endpoint, *MultiEndpoint]{},
	}
}
