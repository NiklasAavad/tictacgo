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
	Command  string `json:"command"`
	Response any    `json:"response"`
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
func (pool *GamePool) execute(instruction string, body int) (string, any) {
	switch instruction {
	case "Start Game":
		return "Board", pool.game.StartGame()
	case "Get Result":
		fmt.Println("Getting the result")
		return "Result", pool.game.GetResult()
	case "Is Game Over":
		fmt.Println("Checking if game is over")
		return "Game Over", pool.game.IsGameOver()
	case "Choose Square":
		position := game.Position(body)
		return "Board", pool.game.ChooseSquare(position)
	case "Change Player In Turn":
		return "Player In Turn", pool.game.ChangePlayerInTurn()
	case "Get Board":
		return "Board", pool.game.Board()
	}

	// TODO bør nok returnerer en error i stedet for
	return "Error", nil
}

func (pool *GamePool) respond(instruction string, body int) GameResponse {
	command, response := pool.execute(instruction, body)
	return GameResponse{Command: command, Response: response}
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
			pool.broadcastResponse(response)
			if message.Instruction == "Choose Square" {
				fmt.Println("Inside if statement in Start method.")
				if pool.game.IsGameOver() {
					isGameOverResponse := pool.respond("Is Game Over", 0)
					resultResponse := pool.respond("Get Result", 0)
					pool.broadcastResponse(isGameOverResponse)
					pool.broadcastResponse(resultResponse)
				}
			}
		}
	}
}
