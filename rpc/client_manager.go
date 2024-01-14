package rpc

import (
	"sync"

	"github.com/VaalaCat/frp-panel/pb"
)

type ClientsManager interface {
	Get(cliID string) *Connector
	Set(cliID, clientType string, sender pb.Master_ServerSendServer)
	Remove(cliID string)
}

type Connector struct {
	CliID   string
	Conn    pb.Master_ServerSendServer
	CliType string
}

type ClientsManagerImpl struct {
	senders *sync.Map
}

// Get implements ClientsManager.
func (c *ClientsManagerImpl) Get(cliID string) *Connector {
	cliAny, ok := c.senders.Load(cliID)
	if !ok {
		return nil
	}

	cli, ok := cliAny.(*Connector)
	if !ok {
		return nil
	}

	return cli
}

// Set implements ClientsManager.
func (c *ClientsManagerImpl) Set(cliID, clientType string, sender pb.Master_ServerSendServer) {
	c.senders.Store(cliID, &Connector{
		CliID:   cliID,
		Conn:    sender,
		CliType: clientType,
	})
}

func (c *ClientsManagerImpl) Remove(cliID string) {
	c.senders.Delete(cliID)
}

var (
	clientsManager *ClientsManagerImpl
)

func NewClientsManager() *ClientsManagerImpl {
	return &ClientsManagerImpl{
		senders: &sync.Map{},
	}
}

func MustInitClientsManager() {
	if clientsManager != nil {
		return
	}

	clientsManager = NewClientsManager()
}

func GetClientsManager() ClientsManager {
	return clientsManager
}
