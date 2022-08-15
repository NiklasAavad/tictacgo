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

// TODO This response is not received by all clients..
func (pool *GamePool) broadcastGameIsOver() {
	gameOverResponse := GameResponse{GAME_OVER, true}

	result := pool.game.GetResult()
	resultResponse := GameResponse{RESULT, result}

	pool.broadcastResponse(gameOverResponse)
	pool.broadcastResponse(resultResponse)
}

func (pool *GamePool) executeInstruction(message GameMessage) (game.Board, error) {
	switch message.Instruction {
	case START_GAME:
		return pool.game.StartGame(), nil
	case CHOOSE_SQUARE:
		position := game.Position(message.Content)
		return pool.game.ChooseSquare(position), nil
	case GET_BOARD:
		return pool.game.Board(), nil
	}

	return game.Board{}, fmt.Errorf("GameInstruction could not be found: %v", message.Instruction)
}

func (pool *GamePool) respond(message GameMessage) error {
	command := BOARD
	body, err := pool.executeInstruction(message)

	if err != nil {
		return err
	}

	response := GameResponse{command, body}
	pool.broadcastResponse(response)

	if pool.game.IsGameOver() {
		pool.broadcastGameIsOver()
	}

	return nil
}

func (pool *GamePool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.registerClient(client)
		case client := <-pool.Unregister:
			pool.unregisterClient(client)
		case message := <-pool.Broadcast:
			if err := pool.respond(message); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
