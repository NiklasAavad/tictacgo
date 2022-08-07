package websocket

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	name string
	Conn *websocket.Conn
	Pool *Pool
}

type Message struct {
	Type   int    `json:"type"`
	Sender string `json:"sender"`
	Body   string `json:"body"`
}

func closeConn(c *Client) {
	c.Pool.Unregister <- c
	c.Conn.Close()
}

func (c *Client) Read() {
	defer closeConn(c)

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		message := Message{Type: messageType, Sender: c.name, Body: string(p)}
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}
