package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Pool interface {
	Register(c Client)
	NewClient(r *http.Request, conn *websocket.Conn) Client
}
