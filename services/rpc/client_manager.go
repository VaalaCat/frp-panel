package rpc

import (
	"sync"
	"time"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils"
	"google.golang.org/grpc/peer"
)

type ClientsManager interface {
	Get(cliID string) *defs.Connector
	Set(cliID, clientType string, sender pb.Master_ServerSendServer) *defs.Connector
	Remove(cliID string)
	ClientAddr(cliID string) string
	ConnectTime(cliID string) (time.Time, bool)
}

type ClientsManagerImpl struct {
	mu          sync.Mutex
	senders     *utils.SyncMap[string, *defs.Connector]
	connectTime *utils.SyncMap[string, time.Time]
	lastSeenAt  *utils.SyncMap[string, time.Time]
}

// Get implements ClientsManager.
func (c *ClientsManagerImpl) Get(cliID string) *defs.Connector {
	cliAny, ok := c.senders.Load(cliID)
	if !ok {
		return nil
	}
	return cliAny
}

// Set implements ClientsManager.
func (c *ClientsManagerImpl) Set(cliID, clientType string, sender pb.Master_ServerSendServer) *defs.Connector {
	c.mu.Lock()
	defer c.mu.Unlock()
	if old := c.Get(cliID); old != nil {
		old.Close()
	}
	connector := defs.NewConnector(cliID, clientType, sender)
	c.senders.Store(cliID, connector)
	c.connectTime.Store(cliID, time.Now())
	c.lastSeenAt.Store(cliID, time.Now())
	return connector
}

func (c *ClientsManagerImpl) Remove(cliID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if old := c.Get(cliID); old != nil {
		old.Close()
	}
	c.remove(cliID)
}

func (c *ClientsManagerImpl) RemoveIf(cliID string, connector *defs.Connector) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Get(cliID) != connector {
		return
	}
	connector.Close()
	c.remove(cliID)
}

func (c *ClientsManagerImpl) remove(cliID string) {
	c.senders.Delete(cliID)
	c.connectTime.Delete(cliID)
	c.lastSeenAt.Delete(cliID)
}

func (c *ClientsManagerImpl) ClientAddr(cliID string) string {
	connector := c.Get(cliID)
	if connector == nil {
		return ""
	}
	p, ok := peer.FromContext(connector.Conn.Context())
	if !ok || p == nil {
		return ""
	}
	return p.Addr.String()
}

func (c *ClientsManagerImpl) ConnectTime(cliID string) (time.Time, bool) {
	t, ok := c.connectTime.Load(cliID)
	if !ok {
		return time.Time{}, false
	}
	return t, true
}

func (c *ClientsManagerImpl) UpdateLastSeenAt(cliID string) {
	c.lastSeenAt.Store(cliID, time.Now())
}

func (c *ClientsManagerImpl) GetLastSeenAt(cliID string) (time.Time, bool) {
	t, ok := c.lastSeenAt.Load(cliID)
	if !ok {
		return time.Time{}, false
	}
	return t, true
}

func NewClientsManager() app.ClientsManager {
	return &ClientsManagerImpl{
		senders:     &utils.SyncMap[string, *defs.Connector]{},
		connectTime: &utils.SyncMap[string, time.Time]{},
		lastSeenAt:  &utils.SyncMap[string, time.Time]{},
	}
}
