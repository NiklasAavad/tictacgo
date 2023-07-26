package gamesocket

import (
	"fmt"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
)

type Command interface {
	execute() ([]ResponseHandler, error)
}

// ------------------------------------------------------------------------------------------------------------

type StartGameCommand struct {
	client *GameClient
}

func NewStartGameCommand(client *GameClient) (*StartGameCommand, error) {
	if client == nil {
		return nil, fmt.Errorf("Client cannot be nil")
	}
	return &StartGameCommand{client}, nil
}

func (c *StartGameCommand) execute() ([]ResponseHandler, error) {
	var response GameResponse

	bothCharactersSelected := c.client.Pool.xClient != nil && c.client.Pool.oClient != nil
	if !bothCharactersSelected { // TODO kunne overveje at lave en ErrorResponse, der bare sendes tilbage til den pågældende client og ikke broadcastes
		return nil, fmt.Errorf("Both characters must be selected, before a game can start")
	}

	isClientPlaying := c.client.Pool.xClient == c.client || c.client.Pool.oClient == c.client
	if !isClientPlaying {
		return nil, fmt.Errorf("Client must be playing to start the game, cannot be a spectator")
	}

	c.client.Pool.game.StartGame()

	response.ResponseType = GAME_STARTED

	responseHandler, err := NewResponseHandler(&response, c.client.Pool.clients)
	if err != nil {
		return nil, err
	}

	return []ResponseHandler{*responseHandler}, nil
}

// ------------------------------------------------------------------------------------------------------------

type ChooseSquareCommand struct {
	client   *GameClient
	position game.Position
}

func NewChooseSquareCommand(client *GameClient, position game.Position) (*ChooseSquareCommand, error) {
	if client == nil {
		return nil, fmt.Errorf("Client cannot be nil")
	}

	if position == game.NO_POSITION {
		return nil, fmt.Errorf("No position was given")
	}

	return &ChooseSquareCommand{client, position}, nil
}

func (c *ChooseSquareCommand) isClientInTurn() bool {
	pool := c.client.Pool

	if pool.game.PlayerInTurn() == game.X {
		return pool.xClient == c.client
	}
	return pool.oClient == c.client
}

func (c *ChooseSquareCommand) execute() ([]ResponseHandler, error) {
	var response GameResponse

	if !c.isClientInTurn() {
		return nil, fmt.Errorf("It was not this client's turn to play")
	}

	board, err := c.client.Pool.game.ChooseSquare(c.position)
	if err != nil {
		return nil, err
	}

	response.ResponseType = BOARD
	response.Body = board

	responseHandler, err := NewResponseHandler(&response, c.client.Pool.clients)
	if err != nil {
		return nil, err
	}

	return []ResponseHandler{*responseHandler}, nil
}

// ------------------------------------------------------------------------------------------------------------

type SelectCharacterCommand struct {
	client    *GameClient
	character game.SquareCharacter
}

func NewSelectCharacterCommand(client *GameClient, character game.SquareCharacter) (*SelectCharacterCommand, error) {
	if client == nil {
		return nil, fmt.Errorf("Client cannot be nil")
	}

	if character == game.EMPTY_CHARACTER {
		return nil, fmt.Errorf("No character was given")
	}

	return &SelectCharacterCommand{client, character}, nil
}

func (c *SelectCharacterCommand) selectCharacter() error {
	pool := c.client.Pool

	if c.character == game.X && pool.xClient == nil {
		pool.xClient = c.client
		return nil
	}

	if c.character == game.O && pool.oClient == nil {
		pool.oClient = c.client
		return nil
	}

	return fmt.Errorf("Character %v is already taken", c.character)
}

func (c *SelectCharacterCommand) execute() ([]ResponseHandler, error) {
	var response GameResponse

	hasClientAlreadySelected := c.client.Pool.xClient == c.client || c.client.Pool.oClient == c.client
	if hasClientAlreadySelected {
		return nil, fmt.Errorf("Client had already selected a character")
	}

	if err := c.selectCharacter(); err != nil {
		return nil, err
	}

	response.ResponseType = CHARACTER_SELECTED
	response.Body = c.character

	responseHandler, err := NewResponseHandler(&response, c.client.Pool.clients)
	if err != nil {
		return nil, err
	}

	return []ResponseHandler{*responseHandler}, nil
}

// ------------------------------------------------------------------------------------------------------------

type RequestDrawCommand struct {
	client *GameClient
}

func NewRequestDrawCommand(client *GameClient) (*RequestDrawCommand, error) {
	if client == nil {
		return nil, fmt.Errorf("Client cannot be nil")
	}
	return &RequestDrawCommand{client}, nil
}

func (c *RequestDrawCommand) execute() ([]ResponseHandler, error) {
	pool := c.client.Pool
	game := pool.game

	if !game.IsStarted() {
		return nil, fmt.Errorf("Game has not started yet, and a draw cannot be requested")
	}

	if game.IsGameOver() {
		return nil, fmt.Errorf("Game is already over, and a draw cannot be requested")
	}

	if pool.IsDrawRequested {
		return nil, fmt.Errorf("A draw has already been requested")
	}

	clientIsPlaying := pool.xClient == c.client || pool.oClient == c.client
	if !clientIsPlaying {
		return nil, fmt.Errorf("Client must be playing to request a draw, cannot be a spectator")
	}

	isOpponent := true
	responseToOpponent := GameResponse{
		ResponseType: REQUEST_DRAW,
		Body:         isOpponent,
	}

	isOpponent = false
	responseToSelf := GameResponse{
		ResponseType: REQUEST_DRAW,
		Body:         isOpponent,
	}

	var opponent *GameClient
	if pool.xClient == c.client {
		opponent = pool.oClient
	} else {
		opponent = pool.xClient
	}

	opponentResponseHandler, err := NewResponseHandler(&responseToOpponent, []*GameClient{opponent})
	if err != nil {
		return nil, err
	}

	selfResponseHandler, err := NewResponseHandler(&responseToSelf, []*GameClient{c.client})
	if err != nil {
		return nil, err
	}

	c.client.Pool.IsDrawRequested = true

	return []ResponseHandler{*opponentResponseHandler, *selfResponseHandler}, nil
}

// ------------------------------------------------------------------------------------------------------------

type RespondToDrawRequestCommand struct {
	client *GameClient
	accept bool
}

func NewRespondToDrawRequestCommand(client *GameClient, accept bool) (*RespondToDrawRequestCommand, error) {
	if client == nil {
		return nil, fmt.Errorf("Client cannot be nil")
	}

	return &RespondToDrawRequestCommand{client, accept}, nil
}

// TODO should add a time limit to respond to a draw request, probably through an attached timestamp?
func (c *RespondToDrawRequestCommand) execute() ([]ResponseHandler, error) {
	pool := c.client.Pool
	game := pool.game

	var response GameResponse
	var receivers []*GameClient

	if !pool.IsDrawRequested {
		return nil, fmt.Errorf("No draw was requested, so there is nothing to respond to")
	}

	if !game.IsStarted() {
		pool.IsDrawRequested = false
		return nil, fmt.Errorf("Game has not started yet, and a draw cannot be requested")
	}

	if game.IsGameOver() {
		pool.IsDrawRequested = false
		return nil, fmt.Errorf("Game is already over, and a draw cannot be requested")
	}

	clientIsPlaying := pool.xClient == c.client || pool.oClient == c.client
	if !clientIsPlaying {
		return nil, fmt.Errorf("Client must be playing to accept a draw, cannot be a spectator. Indicates a larger problem with the server")
	}

	if c.accept {
		game.ForceDraw()
		receivers = c.client.Pool.clients
	} else {
		receivers = []*GameClient{pool.xClient, pool.oClient} // Only opponent and self needs to know that the draw was rejected. No spectators receive the inital message either
	}

	pool.IsDrawRequested = false

	response.Body = c.accept
	response.ResponseType = DRAW_REQUEST_RESPONSE

	responseHandler, err := NewResponseHandler(&response, receivers)
 	if err != nil {
		return nil, err
	}

	return []ResponseHandler{*responseHandler}, nil
}
