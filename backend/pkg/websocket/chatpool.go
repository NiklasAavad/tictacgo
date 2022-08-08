package websocket

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type ChatPool struct {
	Register   chan *Client
	Unregister chan *Client
	clients    map[*Client]bool
	Broadcast  chan Message
}

func NewChatPool() *ChatPool {
	return &ChatPool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (p *ChatPool) registerClient(c *Client) {
	p.clients[c] = true

	body := c.Name + " just joined!"
	msg := Message{Type: websocket.TextMessage, Sender: CHAT_INFO.String(), Body: body}
	p.broadcastMessage(msg)
}

func (p *ChatPool) unregisterClient(c *Client) {
	delete(p.clients, c)

	body := c.Name + " just left..."
	msg := Message{Type: websocket.TextMessage, Sender: CHAT_INFO.String(), Body: body}
	p.broadcastMessage(msg)
}

func (p *ChatPool) broadcastMessage(msg Message) error {
	for client := range p.clients {
		if err := client.Conn.WriteJSON(msg); err != nil {
			return err
		}
	}
	return nil
}

func (pool *ChatPool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.registerClient(client)
			break
		case client := <-pool.Unregister:
			pool.unregisterClient(client)
			break
		case message := <-pool.Broadcast:
			if err := pool.broadcastMessage(message); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
