package rpc

import (
	"time"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils"
	"google.golang.org/grpc/peer"
)

type ClientsManager interface {
	Get(cliID string) *app.Connector
	Set(cliID, clientType string, sender pb.Master_ServerSendServer)
	Remove(cliID string)
	ClientAddr(cliID string) string
	ConnectTime(cliID string) (time.Time, bool)
}

type ClientsManagerImpl struct {
	senders     *utils.SyncMap[string, *app.Connector]
	connectTime *utils.SyncMap[string, time.Time]
}

// Get implements ClientsManager.
func (c *ClientsManagerImpl) Get(cliID string) *app.Connector {
	cliAny, ok := c.senders.Load(cliID)
	if !ok {
		return nil
	}
	return cliAny
}

// Set implements ClientsManager.
func (c *ClientsManagerImpl) Set(cliID, clientType string, sender pb.Master_ServerSendServer) {
	c.senders.Store(cliID, &app.Connector{
		CliID:   cliID,
		Conn:    sender,
		CliType: clientType,
	})
	c.connectTime.Store(cliID, time.Now())
}

func (c *ClientsManagerImpl) Remove(cliID string) {
	c.senders.Delete(cliID)
	c.connectTime.Delete(cliID)
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

func NewClientsManager() *ClientsManagerImpl {
	return &ClientsManagerImpl{
		senders:     &utils.SyncMap[string, *app.Connector]{},
		connectTime: &utils.SyncMap[string, time.Time]{},
	}
}
