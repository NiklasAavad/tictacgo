package chat

import (
	"fmt"
	"net/http"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket"
)

type ChatPool struct {
	register   chan websocket.Client
	Unregister chan websocket.Client
	clients    map[websocket.Client]bool
	Broadcast  chan Message
}

var _ websocket.Pool = new(ChatPool)

func NewChatPool() *ChatPool {
	return &ChatPool{
		register:   make(chan websocket.Client),
		Unregister: make(chan websocket.Client),
		clients:    make(map[websocket.Client]bool),
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

func (p *ChatPool) registerClient(c websocket.Client) {
	p.clients[c] = true

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
