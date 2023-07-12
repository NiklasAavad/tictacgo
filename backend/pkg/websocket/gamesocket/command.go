package gamesocket

import (
	"fmt"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket"
)

type Command interface {
	execute() (GameResponse, error)
}

func ParseCommand(msg GameMessage) (Command, error) {
	gameClient, ok := msg.Client.(*GameClient)
	if !ok {
		return nil, fmt.Errorf("invalid client type: %T", msg.Client)
	}

	switch msg.Instruction {
	case START_GAME:
		return &StartGameCommand{gameClient}, nil
	case CHOOSE_SQUARE:
		position, err := game.ParsePosition(msg.Content)
		if err != nil {
			return nil, err
		}
		return &ChooseSquareCommand{gameClient, position}, nil
	case SELECT_CHARACTER:
		character, err := game.ParseSquareCharacter(msg.Content)
		if err != nil {
			return nil, err
		}
		return &SelectCharacterCommand{gameClient, character}, nil
	}

	return nil, fmt.Errorf("invalid instruction: %s", msg.Instruction)
}

type StartGameCommand struct {
	client *GameClient
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

type ChooseSquareCommand struct {
	client   *GameClient
	position game.Position
}

func (c *ChooseSquareCommand) isClientInTurn() bool {
	if c.client.Pool.game.PlayerInTurn() == game.X {
		return c.client.Pool.xClient == c.client
	}
	return c.client.Pool.oClient == c.client
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

type SelectCharacterCommand struct {
	client    *GameClient
	character game.SquareCharacter
}

func (c *SelectCharacterCommand) selectCharacter(character game.SquareCharacter) error {
	if character == game.X && c.client.Pool.xClient == nil {
		c.client.Pool.xClient = c.client
		return nil
	}

	if character == game.O && c.client.Pool.oClient == nil {
		c.client.Pool.oClient = c.client
		return nil
	}

	return fmt.Errorf("Character %v is already taken", c.character)
}

func (c *SelectCharacterCommand) hasClientAlreadySelected() bool {
	return c.client.Pool.xClient == c.client || c.client.Pool.oClient == c.client
}

func (c *SelectCharacterCommand) execute() (GameResponse, error) {
	var response GameResponse

	if c.hasClientAlreadySelected() {
		return response, fmt.Errorf("Client had already selected a character")
	}

	if err := c.selectCharacter(c.character); err != nil {
		return response, err
	}

	response.ResponseType = CHARACTER_SELECTED
	response.Body = c.character

	return response, nil
}

type NewClientCommand struct {
	client *GameClient
}

func MakeNewClientCommand(client *GameClient) Command {
	return &NewClientCommand{client}
}

func (c *NewClientCommand) getClientName(client websocket.Client) string {
	if client == nil {
		return ""
	}

	return client.Name()
}

func (c *NewClientCommand) execute() (GameResponse, error) {
	var response GameResponse

	isGameStarted := c.client.Pool.game.IsStarted()

	xClientName := c.getClientName(c.client.Pool.xClient)
	oClientName := c.getClientName(c.client.Pool.oClient)

	board := c.client.Pool.game.Board()

	response.ResponseType = WELCOME
	response.Body = WelcomeResponse{
		IsGameStarted: isGameStarted,
		XClient:       xClientName,
		OClient:       oClientName,
		Board:         board,
	}

	return response, nil
}
