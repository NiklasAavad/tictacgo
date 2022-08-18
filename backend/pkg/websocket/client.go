package websocket

import "github.com/gorilla/websocket"

type Client interface {
	Conn() *websocket.Conn
	Name() string
	Read()
}
