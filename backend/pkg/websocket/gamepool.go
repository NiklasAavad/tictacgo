package websocket

import (
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
	Command string `json:"command"`
	Body    any    `json:"body"`
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
func (pool *GamePool) execute(instruction GameInstruction, content int) (string, any) {
	switch instruction {
	case START_GAME:
		return "Board", pool.game.StartGame()
	case CHOOSE_SQUARE:
		position := game.Position(content)
		return "Board", pool.game.ChooseSquare(position)
	case GET_BOARD:
		return "Board", pool.game.Board()
	}

	// TODO bør nok returnerer en error i stedet for
	return "Error", nil
}

// TODO med det her sender vi ikke til klienten at spille er slut og at der er fundet en vinder.
func (pool *GamePool) respond(instruction GameInstruction, content int) GameResponse {
	command, body := pool.execute(instruction, content)
	return GameResponse{Command: command, Body: body}
}

func (pool *GamePool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.registerClient(client)
		case client := <-pool.Unregister:
			pool.unregisterClient(client)
		case message := <-pool.Broadcast:
			response := pool.respond(message.Instruction, message.Content)
			pool.broadcastResponse(response)
		}
	}
}
