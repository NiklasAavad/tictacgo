package gamesocket

import (
	"fmt"
	"net/http"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket"
)

type GamePool struct {
	register   chan websocket.Client
	Unregister chan websocket.Client
	clients    map[websocket.Client]game.SquareCharacter
	Broadcast  chan GameMessage
	game       game.GameService
}

var _ websocket.Pool = new(GamePool)

func NewGamePool() *GamePool {
	return &GamePool{
		register:   make(chan websocket.Client),
		Unregister: make(chan websocket.Client),
		clients:    make(map[websocket.Client]game.SquareCharacter),
		Broadcast:  make(chan GameMessage),
		game:       game.NewGame(),
	}
}

func (p *GamePool) NewClient(w http.ResponseWriter, r *http.Request) websocket.Client {
	clientName := r.URL.Query().Get("name")
	if clientName == "" {
		clientName = "Unknown"
	}

	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+V\n", err)
	}

	client := &GameClient{
		name: clientName,
		conn: conn,
		Pool: p,
	}

	return client
}

func (p *GamePool) Register(c websocket.Client) {
	p.register <- c
}

// TODO lav funktion for at tjekke om value er optaget

func (g *GamePool) registerCharacter(client *GameClient, character game.SquareCharacter) error {
	// TODO når test miljø er sat op, så lav tdd på denne. Start med bare at ændre værdien til client key i clients mappet.
	// efterfølgende, så lav et tjek på, om værdien er optaget.
	return nil
}

func (p *GamePool) broadcastResponse(response GameResponse) error {
	for client := range p.clients {
		if err := client.Conn().WriteJSON(response); err != nil {
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
		case client := <-pool.register:
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
