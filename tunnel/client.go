package tunnel

import (
	"context"

	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/services/client"
	"github.com/VaalaCat/frp-panel/utils"
)

type ClientController interface {
	Add(clientID string, serverID string, clientHandler client.ClientHandler)
	Get(clientID string, serverID string) client.ClientHandler
	Delete(clientID string, serverID string)
	Set(clientID string, serverID string, clientHandler client.ClientHandler)
	Run(clientID string, serverID string) // 不阻塞
	Stop(clientID string, serverID string)
	GetByClient(clientID string) *utils.SyncMap[string, client.ClientHandler]
	DeleteByClient(clientID string)
	RunByClient(clientID string) // 不阻塞
	StopByClient(clientID string)
	StopAll()
	DeleteAll()
	RunAll()
	List() []string
}

type clientController struct {
	clients *utils.SyncMap[string, *utils.SyncMap[string, client.ClientHandler]]
}

var (
	clientControllerInstance *clientController
)

func NewClientController() ClientController {
	return &clientController{
		clients: &utils.SyncMap[string, *utils.SyncMap[string, client.ClientHandler]]{},
	}
}

func GetClientController() ClientController {
	if clientControllerInstance == nil {
		clientControllerInstance = NewClientController().(*clientController)
	}
	return clientControllerInstance
}

func (c *clientController) Add(clientID string, serverID string, clientHandler client.ClientHandler) {
	m, _ := c.clients.LoadOrStore(clientID, &utils.SyncMap[string, client.ClientHandler]{})
	oldClientHandler, loaded := m.LoadAndDelete(serverID)
	if loaded {
		oldClientHandler.Stop()
	}
	m.Store(serverID, clientHandler)
}

func (c *clientController) Get(clientID string, serverID string) client.ClientHandler {
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

func (c *clientController) GetByClient(clientID string) *utils.SyncMap[string, client.ClientHandler] {
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

func (c *clientController) Set(clientID string, serverID string, clientHandler client.ClientHandler) {
	v, _ := c.clients.LoadOrStore(clientID, &utils.SyncMap[string, client.ClientHandler]{})
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
	v.Range(func(k string, v client.ClientHandler) bool {
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
	v.Range(func(k string, v client.ClientHandler) bool {
		v.Stop()
		return true
	})
}

func (c *clientController) StopAll() {
	c.clients.Range(func(k string, v *utils.SyncMap[string, client.ClientHandler]) bool {
		v.Range(func(k string, v client.ClientHandler) bool {
			v.Stop()
			return true
		})
		return true
	})
}

func (c *clientController) DeleteAll() {
	c.clients.Range(func(k string, v *utils.SyncMap[string, client.ClientHandler]) bool {
		c.DeleteByClient(k)
		return true
	})
	c.clients = &utils.SyncMap[string, *utils.SyncMap[string, client.ClientHandler]]{}
}

func (c *clientController) RunAll() {
	c.clients.Range(func(k string, v *utils.SyncMap[string, client.ClientHandler]) bool {
		c.RunByClient(k)
		return true
	})
}

func (c *clientController) List() []string {
	return c.clients.Keys()
}
