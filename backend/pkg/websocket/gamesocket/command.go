package gamesocket

import (
	"fmt"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
)

type Command interface {
	execute() (GameResponse, error)
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

func (c *StartGameCommand) execute() (GameResponse, error) {
	var response GameResponse

	bothCharactersSelected := c.client.Pool.xClient != nil && c.client.Pool.oClient != nil
	if !bothCharactersSelected { // TODO kunne overveje at lave en ErrorResponse, der bare sendes tilbage til den pågældende client og ikke broadcastes
		return response, fmt.Errorf("Both characters must be selected, before a game can start")
	}

	isClientPlaying := c.client.Pool.xClient == c.client || c.client.Pool.oClient == c.client
	if !isClientPlaying {
		return response, fmt.Errorf("Client must be playing to start the game, cannot be a spectator")
	}

	c.client.Pool.game.StartGame()

	response.ResponseType = GAME_STARTED

	return response, nil
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

func (c *ChooseSquareCommand) execute() (GameResponse, error) {
	var response GameResponse

	if !c.isClientInTurn() {
		return response, fmt.Errorf("It was not this client's turn to play")
	}

	board, err := c.client.Pool.game.ChooseSquare(c.position)
	if err != nil {
		return response, err
	}

	response.ResponseType = BOARD
	response.Body = board

	return response, nil
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

func (c *SelectCharacterCommand) execute() (GameResponse, error) {
	var response GameResponse

	hasClientAlreadySelected := c.client.Pool.xClient == c.client || c.client.Pool.oClient == c.client
	if hasClientAlreadySelected {
		return response, fmt.Errorf("Client had already selected a character")
	}

	if err := c.selectCharacter(); err != nil {
		return response, err
	}

	response.ResponseType = CHARACTER_SELECTED
	response.Body = c.character

	return response, nil
}
