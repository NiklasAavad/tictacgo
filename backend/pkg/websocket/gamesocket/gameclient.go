package gamesocket

import (
	"fmt"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket"
)

type GameClient struct {
	ID   string
	name string
	conn websocket.Conn
	Pool *GamePool
}

var _ websocket.Client = new(GameClient)

func (c *GameClient) Conn() websocket.Conn {
	return c.conn
}

func (c *GameClient) Name() string {
	return c.name
}

func (c *GameClient) closeConn() {
	c.Pool.Unregister <- c
	c.Conn().Close()
}

func (c *GameClient) Read() {
	defer c.closeConn()

	for {
		var message GameMessage
		if err := c.Conn().ReadJSON(&message); err != nil {
			fmt.Println(err)
			return
		}

		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}
