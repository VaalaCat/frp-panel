package ws

import "github.com/gorilla/websocket"

type serverIncoming struct {
	conn   *websocket.Conn
	remote string
}

type incomingPacket struct {
	payload  []byte
	endpoint *WSConn
}
