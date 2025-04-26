package tunnel

import (
	"sync"

	"github.com/VaalaCat/frp-panel/app"
	"github.com/fatedier/frp/pkg/metrics"
)

type serverController struct {
	servers *sync.Map
}

func NewServerController() app.ServerController {
	metrics.EnableMem()
	metrics.EnablePrometheus()
	return &serverController{
		servers: &sync.Map{},
	}
}

func (c *serverController) Add(serverID string, serverHandler app.ServerHandler) {
	c.servers.Store(serverID, serverHandler)
}

func (c *serverController) Get(serverID string) app.ServerHandler {
	v, ok := c.servers.Load(serverID)
	if !ok {
		return nil
	}
	return v.(app.ServerHandler)
}

func (c *serverController) Delete(serverID string) {
	c.Stop(serverID)
	c.servers.Delete(serverID)
}

func (c *serverController) Set(serverID string, serverHandler app.ServerHandler) {
	c.servers.Store(serverID, serverHandler)
}

func (c *serverController) Run(serverID string) {
	v, ok := c.servers.Load(serverID)
	if !ok {
		return
	}
	go v.(app.ServerHandler).Run()
}

func (c *serverController) Stop(serverID string) {
	v, ok := c.servers.Load(serverID)
	if !ok {
		return
	}
	v.(app.ServerHandler).Stop()
}

func (c *serverController) List() []string {
	keys := make([]string, 0)
	c.servers.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})
	return keys
}
