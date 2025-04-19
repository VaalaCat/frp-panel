package wsgrpc

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// ---------------------------------------
// 通用 websocketConn 实现 net.Conn 接口
// ---------------------------------------

type websocketConn struct {
	ws         *websocket.Conn
	readMutex  sync.Mutex
	writeMutex sync.Mutex
	// 缓存由于一次读取没有全部消耗完的数据
	readBuffer bytes.Buffer
}

// Read 实现对 websocket 消息的分段读取
func (c *websocketConn) Read(p []byte) (int, error) {
	c.readMutex.Lock()
	defer c.readMutex.Unlock()

	// 若缓冲区为空，则阻塞读取下一条消息
	if c.readBuffer.Len() == 0 {
		messageType, data, err := c.ws.ReadMessage()
		if err != nil {
			return 0, err
		}
		// 只接受二进制数据
		if messageType != websocket.BinaryMessage {
			return 0, fmt.Errorf("unexpected message type: %d", messageType)
		}
		c.readBuffer.Write(data)
	}

	return c.readBuffer.Read(p)
}

// Write 将数据作为单条二进制消息发送
func (c *websocketConn) Write(p []byte) (int, error) {
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()

	err := c.ws.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// Close 关闭 websocket 连接
func (c *websocketConn) Close() error {
	return c.ws.Close()
}

// LocalAddr 返回本地地址，通过 websocket 底层连接获取
func (c *websocketConn) LocalAddr() net.Addr {
	if conn := c.ws.UnderlyingConn(); conn != nil {
		return conn.LocalAddr()
	}
	return nil
}

// RemoteAddr 返回远端地址
func (c *websocketConn) RemoteAddr() net.Addr {
	if conn := c.ws.UnderlyingConn(); conn != nil {
		return conn.RemoteAddr()
	}
	return nil
}

// SetDeadline 同时设置读写超时
func (c *websocketConn) SetDeadline(t time.Time) error {
	if err := c.ws.SetReadDeadline(t); err != nil {
		return err
	}
	return c.ws.SetWriteDeadline(t)
}

// SetReadDeadline 设置读超时
func (c *websocketConn) SetReadDeadline(t time.Time) error {
	return c.ws.SetReadDeadline(t)
}

// SetWriteDeadline 设置写超时
func (c *websocketConn) SetWriteDeadline(t time.Time) error {
	return c.ws.SetWriteDeadline(t)
}

// ---------------------------------------
// 客户端 WebSocket Dialer
// ---------------------------------------

// WebsocketDialer 返回一个可以用于 grpc.WithContextDialer 的拨号函数；该函数通过 websocket 建立连接。
// 参数 url 表示 websocket 服务器地址；header 可用于传递额外的 header 参数。
func WebsocketDialer(url string, header http.Header, insecure bool) func(ctx context.Context, addr string) (net.Conn, error) {
	return func(ctx context.Context, addr string) (net.Conn, error) {
		dialer := websocket.Dialer{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		}
		ws, _, err := dialer.DialContext(ctx, url, header)
		if err != nil {
			return nil, err
		}
		return &websocketConn{ws: ws}, nil
	}
}

// ---------------------------------------
// 服务端 WebSocket Listener 及 Gin Handler
// ---------------------------------------

// WSListener 实现了 net.Listener 接口，用于接收 websocket 升级后的连接。
// gRPC server 可直接传入 WSListener 实例作为监听器调用 Serve 方法。
type WSListener struct {
	connCh chan net.Conn
	mu     sync.Mutex
	closed bool
	addr   net.Addr
	done   chan struct{}
}

// dummyAddr 用于 WSListener 的 Addr 实现
type dummyAddr struct {
	network string
	address string
}

func (d dummyAddr) Network() string {
	return d.network
}

func (d dummyAddr) String() string {
	return d.address
}

// NewWSListener 创建一个 WSListener 实例。
// 参数 addr 表示监听地址，network 建议为固定字符串（例如："ws"），bufSize 为连接队列大小。
func NewWSListener(addr, network string, bufSize int) *WSListener {
	return &WSListener{
		connCh: make(chan net.Conn, bufSize),
		addr:   dummyAddr{network: network, address: addr},
		done:   make(chan struct{}),
	}
}

// Accept 等待并返回下一个连接
func (l *WSListener) Accept() (net.Conn, error) {
	select {
	case conn, ok := <-l.connCh:
		if !ok {
			return nil, fmt.Errorf("listener closed")
		}
		return conn, nil
	case <-l.done:
		return nil, fmt.Errorf("listener closed")
	}
}

// Close 关闭 WSListener
func (l *WSListener) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.closed {
		return nil
	}
	l.closed = true
	close(l.done)
	close(l.connCh)
	return nil
}

// Addr 返回本监听器的地址
func (l *WSListener) Addr() net.Addr {
	return l.addr
}

// GinWSHandler 返回一个 Gin 的 HandlerFunc，用于处理 HTTP 请求，将其升级为 WebSocket 连接
// 并包装为 websocketConn 后推送到 WSListener 中，以供 gRPC server 使用。
// 参数 upgrader 可对 websocket 升级过程进行自定义配置。
func GinWSHandler(listener *WSListener, upgrader *websocket.Upgrader) gin.HandlerFunc {
	return func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.String(http.StatusInternalServerError, "ws upgrade error: %v", err)
			return
		}
		conn := &websocketConn{ws: ws}
		// 非阻塞方式将连接推送到 listener
		select {
		case listener.connCh <- conn:
			// 推送成功后，可选进行应答
		default:
			// 队列满则关闭连接
			ws.Close()
			c.String(http.StatusServiceUnavailable, "connection queue is full")
			return
		}
	}
}

// ------------------------------
// 使用示例
// ------------------------------

// 假设我们有这样一个 main 文件使用上述库：
/*
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"vaalacat/frp-panel/utils/wsgrpc"
	"github.com/gorilla/websocket"
)

// 服务端实例
func main() {
	// 创建 WebSocket Listener，缓冲队列大小为 100，地址和网络标识可自定义
	listener := wsgrpc.NewWSListener("ws-listener", "ws", 100)

	// 在单独的 goroutine 中启动 gRPC Server
	go func() {
		grpcServer := grpc.NewServer()
		// 在此注册你的 gRPC 服务…
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	// 使用 Gin 创建 HTTP 服务器，并在某个路径下提供 WebSocket 功能
	router := gin.Default()

	// 创建一个简单的 upgrader 实例；可根据需要自定义 CheckOrigin 等选项
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	// 注册 WebSocket 处理 handler，路径可自定义，例如 /ws
	router.GET("/ws", wsgrpc.GinWSHandler(listener, upgrader))

	// 启动 HTTP 服务
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}

	// 示例中，当 HTTP 请求升级为 WebSocket 后，会将连接推入 listener，
	// gRPC Server 的 Accept 就会获取到该 net.Conn 连接，实现 gRPC 请求的代理。
}

客户端示例：
func main() {
	// 定义 websocket 服务器地址和 header（如果有需要）
	wsURL := "ws://127.0.0.1:8080/ws" // 示例地址
	header := http.Header{}

	// 创建 websocket dialer
	dialer := wsgrpc.WebsocketDialer(wsURL, header)

	// 使用 grpc.WithContextDialer 配置 GRPC Dial
	conn, err := grpc.DialContext(context.Background(), "ignored",
		grpc.WithContextDialer(dialer),
		grpc.WithInsecure(), // 示例中禁用 TLS，生产环境建议使用安全连接
	)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	// 接下来可使用 conn 创建 GRPC 客户端进行调用
}
*/
