package websocket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type GameClient struct {
	ID   string
	Name string
	Conn *websocket.Conn
	Pool *GamePool
}

type GameMessage struct {
	Instruction string `json:"instruction"`
	Content     int    `json:"content"`
}

func NewGameClient(r *http.Request, conn *websocket.Conn, pool *GamePool) *GameClient {
	clientName := r.URL.Query().Get("name")
	if clientName == "" {
		clientName = "Unknown"
	}

	client := &GameClient{
		Name: clientName,
		Conn: conn,
		Pool: pool,
	}

	return client
}

func closeGameConn(c *GameClient) {
	c.Pool.Unregister <- c
	c.Conn.Close()
}

func (c *GameClient) Read() {
	defer closeGameConn(c)

	for {
		var message GameMessage
		if err := c.Conn.ReadJSON(&message); err != nil {
			fmt.Printf("Message did not match json schema: %+v\n", message)
			continue
		}

		// messageType, p, err := c.Conn.ReadMessage()
		// if err != nil {
		// log.Println(err)
		// return
		// }

		// message := Message{Type: messageType, Sender: c.Name, Body: string(p)}
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}
