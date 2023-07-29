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
	c.Pool.Unregister(c)
	c.Conn().Close()
}

type GameMessage struct {
	Instruction string `json:"instruction"`
	Content     any    `json:"content"`
}

func (client *GameClient) Read() {
	defer client.closeConn()

	for {
		var message GameMessage
		if err := client.Conn().ReadJSON(&message); err != nil {
			fmt.Println(err)
			return
		}

		command, err := NewCommand(message, client)
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := client.Pool.Broadcast(command); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Message Received: %+v\n", message)
	}
}
