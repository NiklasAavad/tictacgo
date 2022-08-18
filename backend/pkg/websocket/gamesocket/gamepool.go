package gamesocket

import (
	"fmt"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
)

type GamePool struct {
	Register   chan *GameClient
	Unregister chan *GameClient
	clients    map[*GameClient]game.SquareCharacter
	Broadcast  chan GameMessage
	game       game.GameService
}

func NewGamePool() *GamePool {
	return &GamePool{
		Register:   make(chan *GameClient),
		Unregister: make(chan *GameClient),
		clients:    make(map[*GameClient]game.SquareCharacter),
		Broadcast:  make(chan GameMessage),
		game:       game.NewGame(),
	}
}

// TODO lav funktion for at tjekke om value er optaget

func (g *GamePool) registerCharacter(client *GameClient, character game.SquareCharacter) error {
	// TODO når test miljø er sat op, så lav tdd på denne. Start med bare at ændre værdien til client key i clients mappet.
	// efterfølgende, så lav et tjek på, om værdien er optaget.
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

func (pool *GamePool) broadcastGameIsOver() {
	gameOverResponse := GameResponse{GAME_OVER, true}

	result := pool.game.GetResult()
	resultResponse := GameResponse{RESULT, result}

	pool.broadcastResponse(gameOverResponse)
	pool.broadcastResponse(resultResponse)
}

func (pool *GamePool) executeMessage(message GameMessage) (game.Board, error) { // TODO implementer Marshaler i Board og brug derefter Marshaler som returtype
	switch message.Instruction {
	case START_GAME:
		return pool.game.StartGame(), nil
	case CHOOSE_SQUARE:
		return pool.game.ChooseSquare(message.Content), nil
	case GET_BOARD:
		return pool.game.Board(), nil
	}

	return game.Board{}, fmt.Errorf("GameInstruction could not be found: %v", message.Instruction)
}

func (pool *GamePool) respond(message GameMessage) error {
	command := BOARD
	body, err := pool.executeMessage(message)

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
			pool.clients[client] = game.EMPTY
		case client := <-pool.Unregister:
			delete(pool.clients, client)
		case message := <-pool.Broadcast:
			if err := pool.respond(message); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
