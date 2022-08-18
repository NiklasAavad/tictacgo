package chat

import (
	"fmt"
	"net/http"

	ws "github.com/NiklasPrograms/tictacgo/backend/pkg/websocket"
	"github.com/gorilla/websocket"
)

type ChatPool struct {
	register   chan ws.Client
	Unregister chan ws.Client
	clients    map[ws.Client]bool
	Broadcast  chan Message
}

var _ ws.Pool = new(ChatPool)

func NewChatPool() *ChatPool {
	return &ChatPool{
		register:   make(chan ws.Client),
		Unregister: make(chan ws.Client),
		clients:    make(map[ws.Client]bool),
		Broadcast:  make(chan Message),
	}
}

const CHAT_INFO = "Chat Info"

func (p *ChatPool) NewClient(r *http.Request, conn *websocket.Conn) ws.Client {
	clientName := r.URL.Query().Get("name")
	if clientName == "" {
		clientName = "Unknown"
	}

	client := &ChatClient{
		name: clientName,
		conn: conn,
		Pool: p,
	}

	return client
}

func (p *ChatPool) Register(c ws.Client) {
	p.register <- c
}

func (p *ChatPool) registerClient(c ws.Client) {
	p.clients[c] = true

	body := c.Name() + " just joined!"
	msg := Message{Type: websocket.TextMessage, Sender: CHAT_INFO, Body: body}
	p.broadcastMessage(msg)
}

func (p *ChatPool) unregisterClient(c ws.Client) {
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
