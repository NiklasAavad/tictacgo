package websocket

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (p *Pool) registerClient(c *Client) {
	p.Clients[c] = true

	body := c.Name + " just joined!"
	msg := Message{Type: websocket.TextMessage, Sender: "Chat Info", Body: body}
	p.broadcastMessage(msg)
}

func (p *Pool) unregisterClient(c *Client) {
	delete(p.Clients, c)

	body := c.Name + " just left..."
	msg := Message{Type: websocket.TextMessage, Sender: "Chat Info", Body: body}
	p.broadcastMessage(msg)
}

func (p *Pool) broadcastMessage(msg Message) error {
	for client := range p.Clients {
		if err := client.Conn.WriteJSON(msg); err != nil {
			return err
		}
	}
	return nil
}

func (pool *Pool) Start() {
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
