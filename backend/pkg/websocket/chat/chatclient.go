package chat

import (
	"fmt"
	"log"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket"
)

var _ websocket.Client = new(ChatClient)

type ChatClient struct {
	ID   string
	name string
	conn websocket.Conn
	Pool *ChatPool
}

type Message struct {
	Type   int    `json:"type"`
	Sender string `json:"sender"`
	Body   string `json:"body"`
}

func (c *ChatClient) Conn() websocket.Conn {
	return c.conn
}

func (c *ChatClient) Name() string {
	return c.name
}

func (c *ChatClient) closeConn() {
	c.Pool.Unregister(c)
	c.Conn().Close()
}

func (c *ChatClient) Read() {
	defer c.closeConn()

	for {
		messageType, p, err := c.Conn().ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		message := Message{Type: messageType, Sender: c.name, Body: string(p)}
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}
