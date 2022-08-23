package gamesocket

import (
	"fmt"
	"net/http"
	"time"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket"
)

type GamePool struct {
	register        chan websocket.Client
	unregister      chan websocket.Client
	clients         map[websocket.Client]game.SquareCharacter
	broadcast       chan Command
	game            game.GameService
	channelStrategy ChannelStrategy
	xClient         websocket.Client
	oClient         websocket.Client
}

var _ websocket.Pool = new(GamePool)

func NewGamePool(cs ChannelStrategy) *GamePool {
	return &GamePool{
		register:        make(chan websocket.Client),
		unregister:      make(chan websocket.Client),
		clients:         make(map[websocket.Client]game.SquareCharacter),
		broadcast:       make(chan Command),
		game:            game.NewGame(),
		channelStrategy: cs,
		xClient:         nil,
		oClient:         nil,
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

func (p *GamePool) Broadcast(m GameMessage) error {
	command, err := ParseCommand(m)
	if err != nil {
		return err
	}

	p.channelStrategy.broadcast(p, command)

	return nil
}

func (p *GamePool) Register(c websocket.Client) {
	p.channelStrategy.register(p, c)
}

func (p *GamePool) Unregister(c websocket.Client) {
	p.channelStrategy.unregister(p, c)
}

func (p *GamePool) Clients() map[websocket.Client]game.SquareCharacter {
	return p.clients
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

func (pool *GamePool) respond(command Command) error {
	response, err := command.execute()

	if err != nil {
		return err
	}

	pool.broadcastResponse(response)

	if pool.game.IsGameOver() {
		pool.broadcastGameIsOver()
		time.Sleep(2500) // Sleep 2.5 seconds to await game ending on client side
	}

	return nil
}

func (pool *GamePool) respondToNewClient(client websocket.Client) error {
	gameClient, ok := client.(*GameClient)
	if !ok {
		return fmt.Errorf("Client could not be casted to a GameClient")
	}

	newClientCommand := &NewClientCommand{gameClient}

	response, err := newClientCommand.execute()
	if err != nil {
		return err
	}

	if err := client.Conn().WriteJSON(response); err != nil {
		return err
	}

	return nil
}

func (pool *GamePool) Start() {
	for {
		select {
		case client := <-pool.register:
			pool.clients[client] = game.EMPTY
			if err := pool.respondToNewClient(client); err != nil {
				fmt.Println(err)
			}
		case client := <-pool.unregister:
			delete(pool.clients, client)
		case command := <-pool.broadcast:
			if err := pool.respond(command); err != nil {
				fmt.Println(err)
			}
		}
	}
}
