package websocket

import (
	"fmt"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
	"github.com/gorilla/websocket"
)

type GamePool struct {
	Register   chan *GameClient
	Unregister chan *GameClient
	clients    map[*GameClient]bool
	Broadcast  chan GameMessage
	game       game.GameService
}

func NewGamePool() *GamePool {
	return &GamePool{
		Register:   make(chan *GameClient),
		Unregister: make(chan *GameClient),
		clients:    make(map[*GameClient]bool),
		Broadcast:  make(chan GameMessage),
		game:       game.NewGame(),
	}
}

type GameResponse struct {
	Response any `json:"response"`
}

func (p *GamePool) registerClient(c *GameClient) {
	p.clients[c] = true // TODO  måske ændres til false for dem som ikke spiller

	body := c.Name + " is ready to play!"
	msg := Message{Type: websocket.TextMessage, Sender: GAME_INFO.String(), Body: body}
	p.broadcastMessage(msg)
}

func (p *GamePool) unregisterClient(c *GameClient) {
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

func (p *GamePool) broadcastResponse(response GameResponse) error {
	for client := range p.clients {
		if err := client.Conn.WriteJSON(response); err != nil {
			return err
		}
	}
	return nil
}

// TODO fix any!
func (pool *GamePool) execute(instruction string, body int) any {
	switch instruction {
	case "Start Game":
		return pool.game.StartGame()
	case "Get Result":
		return pool.game.GetResult()
	case "Is Game Over":
		return pool.game.IsGameOver()
	case "Is Choice Valid":
		position := game.Position(body)
		return pool.game.IsChoiceValid(position)
	case "Choose Square":
		position := game.Position(body)
		return pool.game.ChooseSquare(position)
	case "Change Player In Turn":
		return pool.game.ChangePlayerInTurn()
	}

	// TODO bør nok returnerer en error i stedet for
	return nil
}

func (pool *GamePool) respond(instruction string, body int) GameResponse {
	response := pool.execute(instruction, body)
	return GameResponse{Response: response}
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
			response := pool.respond(message.Instruction, message.Content)
			if err := pool.broadcastResponse(response); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
