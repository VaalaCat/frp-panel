package tunnel

import (
	"sync"

	"github.com/VaalaCat/frp-panel/services/server"
	"github.com/fatedier/frp/pkg/metrics"
)

type ServerController interface {
	Add(serverID string, serverHandler server.ServerHandler)
	Get(serverID string) server.ServerHandler
	Delete(serverID string)
	Set(serverID string, serverHandler server.ServerHandler)
	Run(serverID string) // 不阻塞
	Stop(serverID string)
	List() []string
}

type serverController struct {
	servers *sync.Map
}

var (
	serverControllerInstance *serverController
)

func NewServerController() ServerController {
	metrics.EnableMem()
	metrics.EnablePrometheus()
	return &serverController{
		servers: &sync.Map{},
	}
}

func GetServerController() ServerController {
	if serverControllerInstance == nil {
		serverControllerInstance = NewServerController().(*serverController)
	}
	return serverControllerInstance
}

func (c *serverController) Add(serverID string, serverHandler server.ServerHandler) {
	c.servers.Store(serverID, serverHandler)
}

func (c *serverController) Get(serverID string) server.ServerHandler {
	v, ok := c.servers.Load(serverID)
	if !ok {
		return nil
	}
	return v.(server.ServerHandler)
}

func (c *serverController) Delete(serverID string) {
	c.Stop(serverID)
	c.servers.Delete(serverID)
}

func (c *serverController) Set(serverID string, serverHandler server.ServerHandler) {
	c.servers.Store(serverID, serverHandler)
}

func (c *serverController) Run(serverID string) {
	v, ok := c.servers.Load(serverID)
	if !ok {
		return
	}
	go v.(server.ServerHandler).Run()
}

func (c *serverController) Stop(serverID string) {
	v, ok := c.servers.Load(serverID)
	if !ok {
		return
	}
	v.(server.ServerHandler).Stop()
}

func (c *serverController) List() []string {
	keys := make([]string, 0)
	c.servers.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})
	return keys
}
