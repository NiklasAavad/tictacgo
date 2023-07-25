package gamesocket

import (
	"fmt"
	"net/http"
	"time"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket"
)

// GamePool is a websocket pool that handles the game logic
// It implements the websocket.Pool interface
type GamePool struct {
	register        chan websocket.Client
	unregister      chan websocket.Client
	clients         map[websocket.Client]game.SquareCharacter // TODO: evt lav om til liste af clients?
	broadcast       chan Command
	game            game.GameService
	channelStrategy ChannelStrategy
	xClient         websocket.Client
	oClient         websocket.Client
}

// Assert that GamePool implements the websocket.Pool interface
var _ websocket.Pool = new(GamePool)

// NewGamePool creates a new GamePool
func NewGamePool(cs ChannelStrategy) *GamePool {
	return &GamePool{
		register:        make(chan websocket.Client),
		unregister:      make(chan websocket.Client),
		clients:         make(map[websocket.Client]game.SquareCharacter),
		broadcast:       make(chan Command),
		game:            game.NewGame(),
		channelStrategy: cs,
		xClient:         nil, // TODO: overvej at tilføje en 'noClient'
		oClient:         nil, // TODO: overvej at tilføje en 'noClient'
	}
}

// NewClient creates a new GameClient
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

// Broadcast broadcasts a GameMessage to all clients in the pool
func (p *GamePool) Broadcast(m GameMessage) error {
	command, err := m.ToCommand()
	if err != nil {
		return err
	}

	return p.channelStrategy.broadcast(p, command)
}

// Register registers a client in the pool
func (p *GamePool) Register(c websocket.Client) {
	p.channelStrategy.register(p, c)
}

// Unregister unregisters a client from the pool
func (p *GamePool) Unregister(c websocket.Client) {
	p.channelStrategy.unregister(p, c)
}

// Clients returns a map of all clients in the pool
// The key is the client and the value is the client's character in the game
func (p *GamePool) Clients() map[websocket.Client]game.SquareCharacter {
	return p.clients
}

// broadcastResponse broadcasts a GameResponse to all clients in the pool
// It returns an error if the broadcast fails
func (p *GamePool) broadcastResponse(response GameResponse) error {
	for client := range p.clients {
		if err := client.Conn().WriteJSON(response); err != nil {
			return err
		}
	}
	return nil
}

// broadcastGameIsOver broadcasts a GAME_OVER response to all clients in the pool
// It also broadcasts a RESULT response to all clients in the pool
// The RESULT response contains the result of the game
func (pool *GamePool) broadcastGameIsOver() {
	gameOverResponse := GameResponse{GAME_OVER, true}

	result := pool.game.GetResult()
	resultResponse := GameResponse{RESULT, result}

	pool.broadcastResponse(gameOverResponse)
	pool.broadcastResponse(resultResponse)
}

// respond responds to a command
// It executes the command and broadcasts the response to all clients in the pool
// It returns an error if the command execution fails
func (pool *GamePool) respond(command Command) error {
	response, err := command.execute()

	if err != nil {
		return err
	}

	pool.broadcastResponse(response)

	if pool.game.IsGameOver() {
		pool.broadcastGameIsOver()
		time.Sleep(2500 * time.Millisecond) // Sleep 2.5 seconds to await game ending on client side
	}

	return nil
}

// respondToNewClient responds to a new client
// It sends a WELCOME response to the client
// The WELCOME response contains information about the game
// It returns an error if the broadcast fails
func (pool *GamePool) respondToNewClient(client websocket.Client) error {
	xClientName := ""
	if pool.xClient != nil {
		xClientName = pool.xClient.Name()
	}

	oClientName := ""
	if pool.oClient != nil {
		oClientName = pool.oClient.Name()
	}

	response := GameResponse{
		ResponseType: WELCOME,
		Body: WelcomeResponse{
			IsGameStarted: pool.game.IsStarted(),
			XClient:       xClientName,
			OClient:       oClientName,
			Board:         pool.game.Board(),
		},
	}

	if err := client.Conn().WriteJSON(response); err != nil {
		return err
	}

	return nil
}

// Start starts the GamePool
// It listens for new clients, unregistered clients and commands
// It responds to new clients and commands
func (pool *GamePool) Start() {
	for {
		select {
		case client := <-pool.register:
			pool.clients[client] = game.EMPTY_CHARACTER
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
