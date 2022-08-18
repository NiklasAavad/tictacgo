package websocket

import "github.com/gorilla/websocket"

type Client interface {
	Read()
	Conn() *websocket.Conn
}
