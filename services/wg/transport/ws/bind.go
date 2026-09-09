package ws

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

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
	defaultIncomingChanSize = 2048
	defaultBatchSize        = 128             // 批量处理大小
	wsReadBufferSize        = 64 * 1024       // 64KiB
	wsWriteBufferSize       = 64 * 1024       // 64KiB
	wsMaxMessageSize        = 4 * 1024 * 1024 // 4MiB
	maxPooledPayloadCap     = 64 * 1024       // 64KiB
)

type WSBind struct {
	ctx          *app.Context
	registerChan chan *serverIncoming
	incomingChan chan *incomingPacket

	epDialer *websocket.Dialer

	conns      map[*WSConn]struct{}
	connsMu    sync.RWMutex // 保护 conns map 的并发访问
	opened     atomic.Bool
	packetPool sync.Pool // incomingPacket 对象池
}

func NewWSBind(ctx *app.Context) *WSBind {
	wb := &WSBind{
		ctx: ctx,
		epDialer: &websocket.Dialer{
			ReadBufferSize:    wsReadBufferSize,
			WriteBufferSize:   wsWriteBufferSize,
			EnableCompression: false, // WireGuard 数据已加密，压缩无效且浪费 CPU
			HandshakeTimeout:  10 * time.Second,
			WriteBufferPool:   nil,
		},
		conns: make(map[*WSConn]struct{}),
	}

	wb.packetPool.New = func() interface{} {
		return &incomingPacket{
			payload: make([]byte, 0, 2048),
		}
	}

	return wb
}

// BatchSize implements conn.Bind.
func (w *WSBind) BatchSize() int {
	return defaultBatchSize
}

// Close implements conn.Bind.
func (w *WSBind) Close() error {
	w.opened.Store(false)

	w.connsMu.Lock()
	for conn := range w.conns {
		conn.close()
	}
	w.conns = make(map[*WSConn]struct{})
	w.connsMu.Unlock()

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

	wsConn.wLock.Lock()
	defer wsConn.wLock.Unlock()

	conn, err := wsConn.getConn(w.ctx)
	if err != nil {
		return err
	}

	// 需要限制单条 message 的大小。
	// 超大 message 可能导致对端有大分配
	var (
		writer     io.WriteCloser
		msgBytes   int
		openWriter = func() error { // 打开 writer
			wr, openErr := conn.NextWriter(websocket.BinaryMessage)
			if openErr != nil {
				return openErr
			}
			writer = wr
			msgBytes = 0
			return nil
		}
		flushWriter = func() error { // 关闭 writer
			if writer == nil {
				return nil
			}
			closeErr := writer.Close()
			writer = nil
			msgBytes = 0
			return closeErr
		}
	)

	if err = openWriter(); err != nil {
		conn.Close()
		return fmt.Errorf("ws get writer error: %w", err)
	}

	// 批量写包
	// TLV分割
	// 保证最大包大小不超过wsMaxMessageSize
	// 如果超过，分多次TLV写入
	for _, buf := range bufs {
		if len(buf) == 0 {
			continue
		}
		// TLV 长度字段为 2 字节，单包最大 65535
		// 如果超长，给wg-go报错，他不应该传这么长的包
		if len(buf) > 0xFFFF {
			_ = flushWriter()
			conn.Close()
			return fmt.Errorf("ws packet too large: %d > 65535", len(buf))
		}

		need := 2 + len(buf)
		if need > wsMaxMessageSize {
			_ = flushWriter()
			conn.Close()
			return fmt.Errorf("ws message too large for single packet: need=%d limit=%d", need, wsMaxMessageSize)
		}

		// 若追加后超过单条 message 上限，则先 flush，再开启新 message
		if msgBytes > 0 && msgBytes+need > wsMaxMessageSize {
			if err = flushWriter(); err != nil {
				conn.Close()
				return fmt.Errorf("ws flush error: %w", err)
			}
			if err = openWriter(); err != nil {
				conn.Close()
				return fmt.Errorf("ws get writer error: %w", err)
			}
		}

		// 写 TLV 长度
		lenBuf := [2]byte{byte(len(buf) >> 8), byte(len(buf))}
		if _, err = writer.Write(lenBuf[:]); err != nil {
			_ = flushWriter()
			conn.Close()
			return fmt.Errorf("ws write length error: %w", err)
		}
		// 写 TLV 内容
		if _, err = writer.Write(buf); err != nil {
			_ = flushWriter()
			conn.Close()
			return fmt.Errorf("ws write data error: %w", err)
		}

		msgBytes += need
	}

	// flush 最后一条 message
	if err = flushWriter(); err != nil {
		conn.Close()
		return fmt.Errorf("ws flush error: %w", err)
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
		ReadBufferSize:    wsReadBufferSize,
		WriteBufferSize:   wsWriteBufferSize,
		EnableCompression: false, // WireGuard 数据已加密，压缩无效且浪费 CPU
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(writer, r, nil)
	if err != nil {
		return fmt.Errorf("ws upgrade error: %w", err)
	}

	// 限制单条消息大小，避免 ReadMessage 触发 io.ReadAll 的超大分配
	conn.SetReadLimit(wsMaxMessageSize)

	// 禁用写入截止时间，避免在高负载下超时
	conn.SetWriteDeadline(time.Time{})

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
			w.connsMu.Lock()
			w.conns[wsConn] = struct{}{}
			w.connsMu.Unlock()
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
		// 避免把大 buffer 放回 sync.Pool
		if cap(p.payload) > maxPooledPayloadCap {
			p.payload = make([]byte, 0, 2048)
		} else {
			p.payload = p.payload[:0]
		}
		// 归还 packet 到对象池
		w.packetPool.Put(p)
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
