package websocket

import (
	"fmt"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
	"github.com/gorilla/websocket"
)

type GamePool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
	Game       game.GameService
}

func NewGamePool() *GamePool {
	return &GamePool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		Game:       game.NewGame(),
	}
}

func (p *GamePool) registerClient(c *Client) {
	p.Clients[c] = true // TODO  måske ændres til false for dem som ikke spiller

	body := c.Name + " just joined!"
	msg := Message{Type: websocket.TextMessage, Sender: "Chat Info", Body: body}
	p.broadcastMessage(msg)
}

func (p *GamePool) unregisterClient(c *Client) {
	delete(p.Clients, c)

	body := c.Name + " just left..."
	msg := Message{Type: websocket.TextMessage, Sender: "Chat Info", Body: body}
	p.broadcastMessage(msg)
}

func (p *GamePool) broadcastMessage(msg Message) error {
	for client := range p.Clients {
		if err := client.Conn.WriteJSON(msg); err != nil {
			return err
		}
	}
	return nil
}

func (pool *GamePool) Start() {
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
