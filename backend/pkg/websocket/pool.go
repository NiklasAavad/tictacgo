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

func registerClient(p *Pool, c *Client) {
	p.Clients[c] = true
	fmt.Println("Size of Connection Pool: ", len(p.Clients))

	body := c.Name + " just joined!"
	msg := Message{Type: websocket.TextMessage, Sender: "Chat Info", Body: body}

	for client := range p.Clients {
		fmt.Println(client)
		client.Conn.WriteJSON(msg)
	}
}

func unregisterClient(p *Pool, c *Client) {
	delete(p.Clients, c)
	fmt.Println("Size of Connection pool: ", len(p.Clients))

	body := c.Name + " just left..."
	msg := Message{Type: websocket.TextMessage, Sender: "Chat Info", Body: body}

	for client := range p.Clients {
		client.Conn.WriteJSON(msg)
	}
}

func broadcastMessage(p *Pool, msg Message) error {
	fmt.Println("Sending message to all clients in Pool")
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
			registerClient(pool, client)
			break
		case client := <-pool.Unregister:
			unregisterClient(pool, client)
			break
		case message := <-pool.Broadcast:
			if err := broadcastMessage(pool, message); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
