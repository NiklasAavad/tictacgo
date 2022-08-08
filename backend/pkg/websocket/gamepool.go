package websocket

import (
	"fmt"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
	"github.com/gorilla/websocket"
)

type GamePool struct {
	Register   chan *Client
	Unregister chan *Client
	clients    map[*Client]bool
	Broadcast  chan Message
	game       game.GameService
}

func NewGamePool() *GamePool {
	return &GamePool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		game:       game.NewGame(),
	}
}

func (p *GamePool) registerClient(c *Client) {
	p.clients[c] = true // TODO  måske ændres til false for dem som ikke spiller

	body := c.Name + " is ready to play!"
	msg := Message{Type: websocket.TextMessage, Sender: GAME_INFO.String(), Body: body}
	p.broadcastMessage(msg)
}

func (p *GamePool) unregisterClient(c *Client) {
	delete(p.clients, c)

	body := c.Name + " will no longer play..."
	msg := Message{Type: websocket.TextMessage, Sender: GAME_INFO.String(), Body: body}
	p.broadcastMessage(msg)
}

func (p *GamePool) broadcastMessage(msg Message) error {
	for client := range p.clients {
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
