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

	if pool.DrawRequestHandler.IsDrawRequested {
		return nil, fmt.Errorf("A draw has already been requested")
	}

	isValid, err := IsDrawRequestCommandValid(c.client)
	if !isValid {
		return nil, err
	}

	responseHandlers, err := CreateResponseHandlersForSelfAndOpponent(c.client, REQUEST_DRAW)
	if err != nil {
		return nil, err
	}

	pool.DrawRequestHandler.IsDrawRequested = true
	pool.DrawRequestHandler.DrawRequester = c.client

	return responseHandlers, nil
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

	if !pool.DrawRequestHandler.IsDrawRequested {
		return nil, fmt.Errorf("No draw was requested, so there is nothing to respond to")
	}

	isValid, err := IsDrawRequestCommandValid(c.client)
	if !isValid {
		return nil, err
	}

	isClientTheRequester := pool.DrawRequestHandler.DrawRequester == c.client
	if isClientTheRequester {
		return nil, fmt.Errorf("Client cannot respond to their own draw request")
	}

	if c.accept {
		game.ForceDraw()
		receivers = c.client.Pool.clients
	} else {
		receivers = []*GameClient{pool.xClient, pool.oClient} // Only opponent and self needs to know that the draw was rejected. No spectators receive the inital message either
	}

	pool.DrawRequestHandler.IsDrawRequested = false
	pool.DrawRequestHandler.DrawRequester = nil

	response.Body = c.accept
	response.ResponseType = DRAW_REQUEST_RESPONSE

	responseHandler, err := NewResponseHandler(&response, receivers)
	if err != nil {
		return nil, err
	}

	return []ResponseHandler{*responseHandler}, nil
}

// ------------------------------------------------------------------------------------------------------------

type WithdrawDrawRequestCommand struct {
	client *GameClient
}

func NewWithdrawDrawRequestCommand(client *GameClient) (*WithdrawDrawRequestCommand, error) {
	if client == nil {
		return nil, fmt.Errorf("Client cannot be nil")
	}

	return &WithdrawDrawRequestCommand{client}, nil
}

func (c *WithdrawDrawRequestCommand) execute() ([]ResponseHandler, error) {
	pool := c.client.Pool

	if !pool.DrawRequestHandler.IsDrawRequested {
		return nil, fmt.Errorf("No draw was requested, so there is nothing to withdraw")
	}

	isValid, err := IsDrawRequestCommandValid(c.client)
	if !isValid {
		return nil, err
	}

	isClientTheDrawRequester := pool.DrawRequestHandler.DrawRequester == c.client
	if !isClientTheDrawRequester {
		return nil, fmt.Errorf("Client is not the draw requester, and cannot withdraw the draw request")
	}

	responseHandlers, err := CreateResponseHandlersForSelfAndOpponent(c.client, WITHDRAW_DRAW_REQUEST)
	if err != nil {
		return nil, err
	}

	pool.DrawRequestHandler.IsDrawRequested = false
	pool.DrawRequestHandler.DrawRequester = nil

	return responseHandlers, nil
}
