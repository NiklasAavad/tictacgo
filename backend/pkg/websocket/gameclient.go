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
	Instruction GameInstruction `json:"instruction"`
	Content     int             `json:"content"`
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

func (c *GameClient) closeConn() {
	c.Pool.Unregister <- c
	c.Conn.Close()
}

func (c *GameClient) Read() {
	defer c.closeConn()

	for {
		var message GameMessage
		if err := c.Conn.ReadJSON(&message); err != nil {
			fmt.Println(err)
			return
		}

		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}
