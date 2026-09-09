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

    // 优先读取反向代理传递的真实客户端 IP。
    // gRPC metadata 实际通过 HTTP/2 headers 传输。
    if md, ok := metadata.FromIncomingContext(ctx); ok {
        // 1. 优先 X-Real-IP
        if values := md.Get("x-real-ip"); len(values) > 0 {
            ip := strings.TrimSpace(values[0])
            if net.ParseIP(ip) != nil {
                return ip
            }
        }

        // 2. 其次 X-Forwarded-For，取最左边的客户端 IP
        if values := md.Get("x-forwarded-for"); len(values) > 0 {
            parts := strings.Split(values[0], ",")
            for _, part := range parts {
                ip := strings.TrimSpace(part)
                if net.ParseIP(ip) != nil {
                    return ip
                }
            }
        }
    }

    // 3. 没经过反向代理或没有真实 IP header 时，
    // 保持原项目行为，返回实际 gRPC 对端地址。
    p, ok := peer.FromContext(ctx)
    if !ok || p == nil || p.Addr == nil {
        return ""
    }

    return p.Addr.String()
}
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
