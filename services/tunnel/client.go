package tunnel

import (
	"context"

	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/VaalaCat/frp-panel/utils/logger"
)

type clientController struct {
	clients *utils.SyncMap[string, *utils.SyncMap[string, app.ClientHandler]]
}

func NewClientController() app.ClientController {
	return &clientController{
		clients: &utils.SyncMap[string, *utils.SyncMap[string, app.ClientHandler]]{},
	}
}

func (c *clientController) Add(clientID string, serverID string, clientHandler app.ClientHandler) {
	m, _ := c.clients.LoadOrStore(clientID, &utils.SyncMap[string, app.ClientHandler]{})
	oldClientHandler, loaded := m.LoadAndDelete(serverID)
	if loaded {
		oldClientHandler.Stop()
	}
	m.Store(serverID, clientHandler)
}

func (c *clientController) Get(clientID string, serverID string) app.ClientHandler {
	v, ok := c.clients.Load(clientID)
	if !ok {
		return nil
	}
	vv, ok := v.Load(serverID)
	if !ok {
		return nil
	}
	return vv
}

func (c *clientController) GetByClient(clientID string) *utils.SyncMap[string, app.ClientHandler] {
	v, ok := c.clients.Load(clientID)
	if !ok {
		return nil
	}

	return v
}

func (c *clientController) Delete(clientID string, serverID string) {
	c.Stop(clientID, serverID)
	v, ok := c.clients.Load(clientID)
	if !ok {
		return
	}
	v.Delete(serverID)
}

func (c *clientController) DeleteByClient(clientID string) {
	c.clients.Delete(clientID)
}

func (c *clientController) Set(clientID string, serverID string, clientHandler app.ClientHandler) {
	v, _ := c.clients.LoadOrStore(clientID, &utils.SyncMap[string, app.ClientHandler]{})
	v.Store(serverID, clientHandler)
}

func (c *clientController) Run(clientID string, serverID string) {
	ctx := context.Background()
	v, ok := c.clients.Load(clientID)
	if !ok {
		logger.Logger(ctx).Errorf("cannot get client by clientID, clientID: [%s] serverID: [%s]", clientID, serverID)
		return
	}
	vv, ok := v.Load(serverID)
	if !ok {
		logger.Logger(ctx).Errorf("cannot load client connected server, clientID: [%s] serverID: [%s]", clientID, serverID)
		return
	}

	go vv.Run()
}

func (c *clientController) RunByClient(clientID string) {
	v, ok := c.clients.Load(clientID)
	if !ok {
		return
	}
	v.Range(func(k string, v app.ClientHandler) bool {
		v.Run()
		return true
	})
}

func (c *clientController) Stop(clientID string, serverID string) {
	v, ok := c.clients.Load(clientID)
	if !ok {
		return
	}
	vv, ok := v.Load(serverID)
	if !ok {
		return
	}
	vv.Stop()
}

func (c *clientController) StopByClient(clientID string) {
	v, ok := c.clients.Load(clientID)
	if !ok {
		return
	}
	v.Range(func(k string, v app.ClientHandler) bool {
		v.Stop()
		return true
	})
}

func (c *clientController) StopAll() {
	c.clients.Range(func(k string, v *utils.SyncMap[string, app.ClientHandler]) bool {
		v.Range(func(k string, v app.ClientHandler) bool {
			v.Stop()
			return true
		})
		return true
	})
}

func (c *clientController) DeleteAll() {
	c.clients.Range(func(k string, v *utils.SyncMap[string, app.ClientHandler]) bool {
		c.DeleteByClient(k)
		return true
	})
	c.clients = &utils.SyncMap[string, *utils.SyncMap[string, app.ClientHandler]]{}
}

func (c *clientController) RunAll() {
	c.clients.Range(func(k string, v *utils.SyncMap[string, app.ClientHandler]) bool {
		c.RunByClient(k)
		return true
	})
}

func (c *clientController) List() []string {
	return c.clients.Keys()
}
