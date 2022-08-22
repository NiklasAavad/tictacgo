package gamesocket

import (
	"fmt"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
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
	case GET_BOARD:
		return &GetBoardCommand{gameClient}, nil
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

	bothCharactersSelected := c.client.Pool.isCharacterTaken(game.X) && c.client.Pool.isCharacterTaken(game.O)
	if !bothCharactersSelected {
		return response, fmt.Errorf("Both characters must be selected, before a game can start")
	}

	board := c.client.Pool.game.StartGame()

	response.ResponseType = BOARD
	response.Body = board

	return response, nil
}

type ChooseSquareCommand struct {
	client   *GameClient
	position any
}

func (c *ChooseSquareCommand) execute() (GameResponse, error) {
	var response GameResponse

	isClientInTurn := c.client.Pool.clients[c.client] == c.client.Pool.game.PlayerInTurn()
	if !isClientInTurn {
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

type GetBoardCommand struct {
	client *GameClient
}

func (c *GetBoardCommand) execute() (GameResponse, error) {
	var response GameResponse

	if !c.client.Pool.game.IsStarted() {
		return response, fmt.Errorf("game is not started yet")
	}

	board := c.client.Pool.game.Board()

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
