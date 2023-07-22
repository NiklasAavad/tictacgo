package chat

import (
	"fmt"
	"net/http"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket"
)

type ChatPool struct {
	register   chan websocket.Client
	unregister chan websocket.Client
	clients    map[websocket.Client]game.SquareCharacter
	Broadcast  chan Message
}

var _ websocket.Pool = new(ChatPool)

func NewChatPool() *ChatPool {
	return &ChatPool{
		register:   make(chan websocket.Client),
		unregister: make(chan websocket.Client),
		clients:    make(map[websocket.Client]game.SquareCharacter),
		Broadcast:  make(chan Message),
	}
}

const CHAT_INFO = "Chat Info"

func (p *ChatPool) NewClient(w http.ResponseWriter, r *http.Request) websocket.Client {
	clientName := r.URL.Query().Get("name")
	if clientName == "" {
		clientName = "Unknown"
	}

	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+V\n", err)
	}

	client := &ChatClient{
		name: clientName,
		conn: conn,
		Pool: p,
	}

	return client
}

func (p *ChatPool) Register(c websocket.Client) {
	p.register <- c
}

func (p *ChatPool) Unregister(c websocket.Client) {
	p.unregister <- c
}

func (p *ChatPool) Clients() map[websocket.Client]game.SquareCharacter {
	return p.clients
}

func (p *ChatPool) registerClient(c websocket.Client) {
	p.clients[c] = game.EMPTY_CHARACTER

	body := c.Name() + " just joined!"
	msg := Message{Type: websocket.TextMessage, Sender: CHAT_INFO, Body: body}
	p.broadcastMessage(msg)
}

func (p *ChatPool) unregisterClient(c websocket.Client) {
	delete(p.clients, c)

	body := c.Name() + " just left..."
	msg := Message{Type: websocket.TextMessage, Sender: CHAT_INFO, Body: body}
	p.broadcastMessage(msg)
}

func (p *ChatPool) broadcastMessage(msg Message) error {
	for client := range p.clients {
		if err := client.Conn().WriteJSON(msg); err != nil {
			return err
		}
	}
	return nil
}

func (pool *ChatPool) Start() {
	for {
		select {
		case client := <-pool.register:
			pool.registerClient(client)
			break
		case client := <-pool.unregister:
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
