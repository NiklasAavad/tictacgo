package gamesocket

import (
	"fmt"
	"net/http"
	"time"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket"
)

// GamePool is a websocket pool that handles the game logic
type GamePool struct {
	register           chan *GameClient
	unregister         chan *GameClient
	clients            []*GameClient
	broadcast          chan Command
	game               game.GameService
	channelStrategy    ChannelStrategy
	xClient            *GameClient
	oClient            *GameClient
	DrawRequestHandler *DrawRequestHandler
}

// NewGamePool creates a new GamePool
func NewGamePool(cs ChannelStrategy) *GamePool {
	return &GamePool{
		register:           make(chan *GameClient),
		unregister:         make(chan *GameClient),
		clients:            []*GameClient{},
		broadcast:          make(chan Command),
		game:               game.NewGame(),
		channelStrategy:    cs,
		xClient:            nil, // TODO: overvej at tilføje en 'noClient'
		oClient:            nil, // TODO: overvej at tilføje en 'noClient'
		DrawRequestHandler: NewDrawRequestHandler(),
	}
}

// NewClient creates a new GameClient
func (p *GamePool) NewClient(w http.ResponseWriter, r *http.Request) *GameClient {
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
func (p *GamePool) Register(c *GameClient) {
	p.channelStrategy.register(p, c)
}

// Unregister unregisters a client from the pool
// returns an error for testing purposes. Will always return nil using the concurrent channel strategy
func (p *GamePool) Unregister(c *GameClient) error {
	return p.channelStrategy.unregister(p, c)
}

// Clients returns all clients in the pool
func (p *GamePool) Clients() []*GameClient {
	return p.clients
}

// broadcastResponse broadcasts a GameResponse to all clients in the pool
// It returns an error if the broadcast fails
func (p *GamePool) broadcastResponse(response GameResponse) error {
	for _, client := range p.clients {
		if err := client.Conn().WriteJSON(response); err != nil {
			return err
		}
	}
	return nil
}

// broadcastToSelected broadcasts a GameResponse to selected clients in the pool
// The argument responseHandlers contains the response and the clients to which the response should be sent
// It returns an error if the broadcast fails
func (p *GamePool) broadcastToSelected(responseHandlers []ResponseHandler) error {
	for _, responseHandler := range responseHandlers {
		for _, client := range responseHandler.Receivers {
			if err := client.Conn().WriteJSON(responseHandler.Response); err != nil {
				// TODO consider if we should return early or continue (and thereby collect all errors)
				return err
			}
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
	responseHandlers, err := command.execute()

	if err != nil {
		return err
	}

	if err := pool.broadcastToSelected(responseHandlers); err != nil {
		return err
	}

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
			pool.clients = append(pool.clients, client)
			if err := pool.respondToNewClient(client); err != nil {
				fmt.Println(err)
			}
		case client := <-pool.unregister:
			updatedClients, err := RemoveElement(pool.clients, client)
			if err != nil {
				fmt.Println(err)
			} else {
				pool.clients = updatedClients
			}
		case command := <-pool.broadcast:
			if err := pool.respond(command); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (pool *GamePool) ServeWs(w http.ResponseWriter, r *http.Request) {
	client := pool.NewClient(w, r)
	pool.Register(client)
	client.Read()
}
