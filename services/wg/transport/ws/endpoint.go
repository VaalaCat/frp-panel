package ws

import (
	"fmt"
	"net/http"
	"net/netip"
	"net/url"
	"sync"
	"time"

	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/gorilla/websocket"
	"golang.zx2c4.com/wireguard/conn"
)

var (
	_ conn.Endpoint = (*WSConn)(nil)
)

type WSConn struct {
	conn    *websocket.Conn
	srcAddr netip.Addr
	dstAddr netip.Addr
	dstText string

	epUrl *url.URL // endpoint URL, used to create a new connection if needed, only used by client

	wsBind *WSBind

	wLock sync.RWMutex
}

// ClearSrc implements conn.Endpoint.
func (w *WSConn) ClearSrc() {
	w.srcAddr = netip.Addr{}
}

// DstIP implements conn.Endpoint.
func (w *WSConn) DstIP() netip.Addr {
	return w.dstAddr
}

// DstToBytes implements conn.Endpoint.
func (w *WSConn) DstToBytes() []byte {
	return w.dstAddr.AsSlice()
}

// DstToString implements conn.Endpoint.
func (w *WSConn) DstToString() string {
	return w.dstText
}

// SrcIP implements conn.Endpoint.
func (w *WSConn) SrcIP() netip.Addr {
	return w.srcAddr
}

// SrcToString implements conn.Endpoint.
func (w *WSConn) SrcToString() string {
	return w.srcAddr.String()
}

// readLoop read messages from websocket connection and send to incoming channel
// client readLoop is started by getConn
// server readLoop is started by serverLoop
func (w *WSConn) readLoop(ctx *app.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if w.wsBind.incomingChan == nil {
				ctx.Logger().Error("ws recv channel is nil, skip read")
				return
			}

			if w.conn == nil {
				ctx.Logger().Error("ws connection is nil, skip read")
				time.Sleep(time.Second)
				continue
			}

			msgType, data, err := w.conn.ReadMessage()
			if err != nil {
				ctx.Logger().WithError(err).Error("ws read message error, close connection")
				w.close()
				return
			}
			if msgType != websocket.BinaryMessage {
				ctx.Logger().Debugf("ws read message type %d, data length %d, skip", msgType, len(data))
				continue
			}

			pkt := &incomingPacket{
				payload:  data,
				endpoint: w,
			}

			select {
			case w.wsBind.incomingChan <- pkt:
			case <-ctx.Done():
				return
			default:
				ctx.Logger().Debugf("ws recv channel is full, drop packet, length %d", len(data))
			}
		}
	}
}

// getConn create a new connection if needed
// intended to be used by server and client
func (w *WSConn) getConn(ctx *app.Context) (*websocket.Conn, error) {
	if w.conn != nil {
		return w.conn, nil
	}

	if w.epUrl == nil {
		return nil, fmt.Errorf("ws endpoint missing URL")
	}

	conn, resp, err := w.wsBind.epDialer.DialContext(ctx, w.epUrl.String(), http.Header{})
	if err != nil {
		return nil, err
	}
	if resp != nil {
		_ = resp.Body.Close()
	}

	w.conn = conn

	go w.readLoop(ctx)
	return conn, nil
}

func (w *WSConn) close() {
	if w.conn != nil {
		_ = w.conn.Close()
		w.conn = nil
	}
	delete(w.wsBind.conns, w)
}
