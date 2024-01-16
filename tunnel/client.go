package tunnel

import (
	"sync"

	"github.com/VaalaCat/frp-panel/services/client"
)

type ClientController interface {
	Add(clientID string, clientHandler client.ClientHandler)
	Get(clientID string) client.ClientHandler
	Delete(clientID string)
	Set(clientID string, clientHandler client.ClientHandler)
	Run(clientID string) // 不阻塞
	Stop(clientID string)
	List() []string
}

type clientController struct {
	clients *sync.Map
}

var (
	clientControllerInstance *clientController
)

func NewClientController() ClientController {
	return &clientController{
		clients: &sync.Map{},
	}
}

func GetClientController() ClientController {
	if clientControllerInstance == nil {
		clientControllerInstance = NewClientController().(*clientController)
	}
	return clientControllerInstance
}

func (c *clientController) Add(clientID string, clientHandler client.ClientHandler) {
	c.clients.Store(clientID, clientHandler)
}

func (c *clientController) Get(clientID string) client.ClientHandler {
	v, ok := c.clients.Load(clientID)
	if !ok {
		return nil
	}
	return v.(client.ClientHandler)
}

func (c *clientController) Delete(clientID string) {
	c.Stop(clientID)
	c.clients.Delete(clientID)
}

func (c *clientController) Set(clientID string, clientHandler client.ClientHandler) {
	c.clients.Store(clientID, clientHandler)
}

func (c *clientController) Run(clientID string) {
	v, ok := c.clients.Load(clientID)
	if !ok {
		return
	}
	go v.(client.ClientHandler).Run()
}

func (c *clientController) Stop(clientID string) {
	v, ok := c.clients.Load(clientID)
	if !ok {
		return
	}
	v.(client.ClientHandler).Stop()
}

func (c *clientController) List() []string {
	keys := make([]string, 0)
	c.clients.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})
	return keys
}
