package defs

import (
	"context"
	"sync"

	"github.com/VaalaCat/frp-panel/pb"
)

type Connector struct {
	CliID     string
	Conn      pb.Master_ServerSendServer
	CliType   string
	sendGate  chan struct{}
	done      chan struct{}
	closeOnce sync.Once
}

func NewConnector(id, clientType string, conn pb.Master_ServerSendServer) *Connector {
	return &Connector{CliID: id, CliType: clientType, Conn: conn, sendGate: make(chan struct{}, 1), done: make(chan struct{})}
}

func (c *Connector) Done() <-chan struct{} { return c.done }

func (c *Connector) Close() { c.closeOnce.Do(func() { close(c.done) }) }

// Send serializes gRPC writes. Closing the connector lets its owner end the
// server stream, which also releases a Send blocked by a disconnected peer.
func (c *Connector) Send(ctx context.Context, msg *pb.ServerMessage) error {
	select {
	case c.sendGate <- struct{}{}:
	case <-ctx.Done():
		return ctx.Err()
	case <-c.done:
		return context.Canceled
	case <-c.Conn.Context().Done():
		return c.Conn.Context().Err()
	}
	result := make(chan error, 1)
	go func() {
		defer func() { <-c.sendGate }()
		if err := ctx.Err(); err != nil {
			result <- err
			return
		}
		result <- c.Conn.Send(msg)
	}()
	select {
	case err := <-result:
		if err != nil {
			c.Close()
		}
		return err
	case <-ctx.Done():
		c.Close()
		return ctx.Err()
	case <-c.done:
		return context.Canceled
	case <-c.Conn.Context().Done():
		return c.Conn.Context().Err()
	}
}
