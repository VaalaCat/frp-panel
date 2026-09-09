package rpc

import (
    "net"
    "strings"
    "time"

    "github.com/VaalaCat/frp-panel/defs"
    "github.com/VaalaCat/frp-panel/pb"
    "github.com/VaalaCat/frp-panel/services/app"
    "github.com/VaalaCat/frp-panel/utils"
    "google.golang.org/grpc/metadata"
    "google.golang.org/grpc/peer"
)

type ClientsManager interface {
	Get(cliID string) *defs.Connector
	Set(cliID, clientType string, sender pb.Master_ServerSendServer)
	Remove(cliID string)
	ClientAddr(cliID string) string
	ConnectTime(cliID string) (time.Time, bool)
}

type ClientsManagerImpl struct {
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
func (c *ClientsManagerImpl) Set(cliID, clientType string, sender pb.Master_ServerSendServer) {
	c.senders.Store(cliID, &defs.Connector{
		CliID:   cliID,
		Conn:    sender,
		CliType: clientType,
	})
	c.connectTime.Store(cliID, time.Now())
	c.lastSeenAt.Store(cliID, time.Now())
}

func (c *ClientsManagerImpl) Remove(cliID string) {
	c.senders.Delete(cliID)
	c.connectTime.Delete(cliID)
	c.lastSeenAt.Delete(cliID)
}

func (c *ClientsManagerImpl) ClientAddr(cliID string) string {
    connector := c.Get(cliID)
    if connector == nil {
        return ""
    }

    ctx := connector.Conn.Context()

    if md, ok := metadata.FromIncomingContext(ctx); ok {
        // Prefer the real client IP set by a trusted reverse proxy.
        if values := md.Get("x-real-ip"); len(values) > 0 {
            if ip := strings.TrimSpace(values[0]); net.ParseIP(ip) != nil {
                return ip
            }
        }

        // Fall back to the first valid address in X-Forwarded-For.
        if values := md.Get("x-forwarded-for"); len(values) > 0 {
            for part := range strings.SplitSeq(values[0], ",") {
                if ip := strings.TrimSpace(part); net.ParseIP(ip) != nil {
                    return ip
                }
            }
        }
    }

    // Preserve the original behavior for direct gRPC connections.
    p, ok := peer.FromContext(ctx)
    if !ok || p == nil || p.Addr == nil {
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
