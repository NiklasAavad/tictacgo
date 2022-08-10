package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Name string
	Conn *websocket.Conn
	Pool *ChatPool
}

type Message struct {
	Type   int    `json:"type"`
	Sender string `json:"sender"`
	Body   string `json:"body"`
}

func NewClient(r *http.Request, conn *websocket.Conn, pool *ChatPool) *Client {
	clientName := r.URL.Query().Get("name")
	if clientName == "" {
		clientName = "Unknown"
	}

	client := &Client{
		Name: clientName,
		Conn: conn,
		Pool: pool,
	}

	return client
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

		message := Message{Type: messageType, Sender: c.Name, Body: string(p)}
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}
