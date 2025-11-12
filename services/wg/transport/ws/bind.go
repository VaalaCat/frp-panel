package ws

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync/atomic"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/gorilla/websocket"
	"golang.zx2c4.com/wireguard/conn"
)

var (
	_ conn.Bind = (*WSBind)(nil)
)

const (
	defaultRegisterChanSize = 128
	defaultIncomingChanSize = 256
)

type WSBind struct {
	ctx          *app.Context
	registerChan chan *serverIncoming
	incomingChan chan *incomingPacket

	epDialer *websocket.Dialer

	conns  map[*WSConn]struct{}
	opened atomic.Bool
}

func NewWSBind(ctx *app.Context) *WSBind {
	return &WSBind{
		ctx:      ctx,
		epDialer: &websocket.Dialer{},
		conns:    make(map[*WSConn]struct{}),
	}
}

// BatchSize implements conn.Bind.
func (w *WSBind) BatchSize() int {
	return 1
}

// Close implements conn.Bind.
func (w *WSBind) Close() error {
	w.opened.Store(false)

	for conn := range w.conns {
		conn.close()
	}
	w.conns = make(map[*WSConn]struct{})

	// 关闭旧的 channels 以释放阻塞的 goroutines
	if w.registerChan != nil {
		close(w.registerChan)
	}
	if w.incomingChan != nil {
		close(w.incomingChan)
	}

	return nil
}

// Open implements conn.Bind.
func (w *WSBind) Open(port uint16) (fns []conn.ReceiveFunc, actualPort uint16, err error) {
	if w.opened.Load() {
		w.ctx.Logger().Debugf("ws bind already opened, closing and reopening for port %d", port)
		if closeErr := w.Close(); closeErr != nil {
			w.ctx.Logger().WithError(closeErr).Warnf("failed to close ws bind before reopening")
		}
	}

	// 重新创建 channels(在 Close() 后它们已经被关闭)
	w.registerChan = make(chan *serverIncoming, defaultRegisterChanSize)
	w.incomingChan = make(chan *incomingPacket, defaultIncomingChanSize)
	w.opened.Store(true)

	go w.serverLoop()
	return []conn.ReceiveFunc{w.recvFunc}, port, nil
}

// ParseEndpoint implements conn.Bind.
func (w *WSBind) ParseEndpoint(s string) (conn.Endpoint, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "ws" && u.Scheme != "wss" {
		return nil, fmt.Errorf("unsupported scheme %q", u.Scheme)
	}

	// if path is empty, use default path
	if u.Path == "" {
		u.Path = defs.DefaultWSHandlerPath
	}

	ep := &WSConn{
		wsBind:  w,
		dstText: u.String(),
		epUrl:   u,
	}
	return ep, nil
}

// Send implements conn.Bind.
func (w *WSBind) Send(bufs [][]byte, ep conn.Endpoint) error {
	wsConn, ok := ep.(*WSConn)
	if !ok {
		return fmt.Errorf("wrong endpoint type")
	}

	conn, err := wsConn.getConn(w.ctx)
	if err != nil {
		return err
	}

	for _, buf := range bufs {
		if len(buf) == 0 {
			continue
		}

		wsConn.wLock.Lock()
		err = conn.WriteMessage(websocket.BinaryMessage, buf)
		wsConn.wLock.Unlock()
		if err != nil {
			conn.Close()
			return fmt.Errorf("ws send message error: %w", err)
		}
	}

	return nil
}

// SetMark implements conn.Bind.
func (w *WSBind) SetMark(mark uint32) error {
	return nil
}

func (w *WSBind) HandleHTTP(writer http.ResponseWriter, r *http.Request) error {
	if !w.opened.Load() {
		return net.ErrClosed
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(writer, r, nil)
	if err != nil {
		return fmt.Errorf("ws upgrade error: %w", err)
	}

	select {
	case w.registerChan <- &serverIncoming{
		conn:   conn,
		remote: conn.RemoteAddr().String(),
	}:
	case <-w.ctx.Done():
		conn.Close()
		return net.ErrClosed
	}

	return nil
}

func (w *WSBind) serverLoop() {
	for {
		select {
		case <-w.ctx.Done():
			return
		case incoming, ok := <-w.registerChan:
			if !ok {
				return
			}
			wsConn := &WSConn{
				conn:    incoming.conn,
				dstText: incoming.remote,
				wsBind:  w,
			}
			w.conns[wsConn] = struct{}{}
			go wsConn.readLoop(w.ctx)
		}
	}
}

// recvFunc receive packets from incoming channel and copy to packets
// idx is the same between packets and sizes and eps
// return the number of batches received
func (w *WSBind) recvFunc(packets [][]byte, sizes []int, eps []conn.Endpoint) (int, error) {
	if !w.opened.Load() {
		return 0, net.ErrClosed
	}

	incoming := w.incomingChan
	if incoming == nil {
		return 0, net.ErrClosed
	}

	max := max(len(packets), len(sizes), len(eps))

	if max == 0 {
		return 0, nil
	}

	var done <-chan struct{}
	if w.ctx != nil {
		done = w.ctx.Done()
	}

	var (
		pkt *incomingPacket
		ok  bool
	)

	select {
	case <-done:
		return 0, net.ErrClosed
	case pkt, ok = <-incoming:
		if !ok {
			return 0, net.ErrClosed
		}
	}

	total := 0
	copyPacket := func(idx int, p *incomingPacket) {
		if idx >= max || p == nil {
			return
		}
		buf := packets[idx]
		payload := p.payload
		if len(buf) < len(payload) {
			// may truncate payload when wg-go give small buffer size
			copy(buf, payload[:len(buf)])
			sizes[idx] = len(buf)
		} else {
			copy(buf, payload)
			sizes[idx] = len(payload)
		}
		eps[idx] = p.endpoint
	}

	copyPacket(total, pkt)
	total++

	for total < max {
		select {
		case pkt, ok = <-incoming:
			if !ok {
				return total, nil
			}
			copyPacket(total, pkt)
			total++
		case <-done:
			return total, nil
		default:
			return total, nil
		}
	}

	return total, nil
}
