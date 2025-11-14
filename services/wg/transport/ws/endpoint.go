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
	conn := w.conn
	incomingChan := w.wsBind.incomingChan
	done := ctx.Done()

	if incomingChan == nil {
		ctx.Logger().Error("ws recv channel is nil, skip read")
		return
	}

	if conn == nil {
		ctx.Logger().Error("ws connection is nil, skip read")
		return
	}

	for {
		// data 是 TLV 格式
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			ctx.Logger().WithError(err).Error("ws read message error, close connection")
			w.close()
			return
		}

		select {
		case <-done:
			return
		default:
		}

		if msgType != websocket.BinaryMessage {
			continue
		}

		// TLV解包
		offset := 0
		for offset+2 <= len(data) {
			// 读取包长度
			pktLen := int(data[offset])<<8 | int(data[offset+1])
			offset += 2

			if offset+pktLen > len(data) {
				ctx.Logger().Errorf("invalid packet length: %d, remaining: %d", pktLen, len(data)-offset)
				break
			}

			pkt := w.wsBind.packetPool.Get().(*incomingPacket)

			// 直接引用原始数据的切片，避免拷贝
			if cap(pkt.payload) >= pktLen {
				pkt.payload = pkt.payload[:pktLen]
				copy(pkt.payload, data[offset:offset+pktLen])
			} else {
				pkt.payload = make([]byte, pktLen)
				copy(pkt.payload, data[offset:offset+pktLen])
			}
			pkt.endpoint = w

			select {
			case incomingChan <- pkt:
				// 成功发送
			case <-done:
				// context 已取消，归还并退出
				w.wsBind.packetPool.Put(pkt)
				return
			default:
				// channel 满了，丢弃包并归还到对象池
				w.wsBind.packetPool.Put(pkt)
			}

			offset += pktLen
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

	// 禁用写入截止时间，避免在高负载下超时
	conn.SetWriteDeadline(time.Time{})

	w.conn = conn

	go w.readLoop(ctx)
	return conn, nil
}

func (w *WSConn) close() {
	if w.conn != nil {
		_ = w.conn.Close()
		w.conn = nil
	}
	w.wsBind.connsMu.Lock()
	delete(w.wsBind.conns, w)
	w.wsBind.connsMu.Unlock()
}
