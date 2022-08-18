package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// overvej at lave flere metoder / konstanter, så vi ikke behøver at importere gorilla andre steder end her.
const (
	TextMessage = 1
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// No checking for now, just allow any connection
	CheckOrigin: func(r *http.Request) bool { return true },
}

func Upgrade(w http.ResponseWriter, r *http.Request) (Conn, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return ws, err
	}
	return ws, nil
}
