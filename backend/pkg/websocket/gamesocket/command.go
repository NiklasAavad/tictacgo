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
		return &ChooseSquareCommand{gameClient, msg.Content}, nil
	case SELECT_CHARACTER:
		return &SelectCharacterCommand{gameClient, msg.Content}, nil
	}

	return nil, fmt.Errorf("invalid instruction: %s", msg.Instruction)
}

type StartGameCommand struct {
	client *GameClient
}

func (c *StartGameCommand) execute() (GameResponse, error) {
	var response GameResponse

	bothCharactersSelected := c.client.Pool.xClient != nil && c.client.Pool.oClient != nil
	if !bothCharactersSelected {
		return response, fmt.Errorf("Both characters must be selected, before a game can start")
	}

	c.client.Pool.game.StartGame()

	response.ResponseType = GAME_STARTED

	return response, nil
}

type ChooseSquareCommand struct {
	client   *GameClient
	position any
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

	position, err := game.ParsePosition(c.position)
	if err != nil {
		return response, err
	}

	board, err := c.client.Pool.game.ChooseSquare(position)
	if err != nil {
		return response, err
	}

	response.ResponseType = BOARD
	response.Body = board

	return response, nil
}

type SelectCharacterCommand struct {
	client    *GameClient
	character any
}

func (c *SelectCharacterCommand) execute() (GameResponse, error) {
	var response GameResponse

	character, err := game.ParseSquareCharacter(c.character)
	if err != nil {
		return response, err
	}

	if err := c.client.Pool.registerCharacter(c.client, character); err != nil {
		return response, err
	}

	response.ResponseType = CHARACTER_SELECTED
	response.Body = character

	return response, nil
}

type NewClientCommand struct {
	client *GameClient
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
