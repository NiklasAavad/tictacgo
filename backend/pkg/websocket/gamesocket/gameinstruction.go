package gamesocket

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
)

type GameInstruction interface {
	String() string
	ToCommand(*GameClient, any) (Command, error)
}

// ---------------------------------------------------------------------------------------------------

type GameInstructionParser struct {
	gi GameInstruction
}

var _ json.Unmarshaler = new(GameInstructionParser)

var (
	StartGame       = new(StartGameInstruction)
	ChooseSquare    = new(ChooseSquareInstruction)
	SelectCharacter = new(SelectCharacterInstruction)
)

// map from strings to GameInstruction
var GameInstructionValue = map[string]GameInstruction{
	"start game":       StartGame,
	"choose square":    ChooseSquare,
	"select character": SelectCharacter,
}

func ParseGameInstruction(s string) (GameInstruction, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	gameInstruction, ok := GameInstructionValue[s]
	if !ok {
		return nil, fmt.Errorf("%q is not a valid Game Instruction", s)
	}
	return gameInstruction, nil
}

func (parser *GameInstructionParser) UnmarshalJSON(data []byte) (err error) {
	var gameInstruction string

	if err := json.Unmarshal(data, &gameInstruction); err != nil {
		return err
	}

	gi, err := ParseGameInstruction(gameInstruction)
	if err != nil {
		return err
	}

	parser.gi = gi

	return nil
}

// ---------------------------------------------------------------------------------------------------

type StartGameInstruction struct{}

// String implements GameInstruction
func (*StartGameInstruction) String() string {
	return "start game"
}

// ToCommand implements GameInstruction
// TODO consider if we should check that content is nil / 0 / whatever
func (*StartGameInstruction) ToCommand(gc *GameClient, content any) (Command, error) {
	return &StartGameCommand{gc}, nil
}

var _ GameInstruction = new(StartGameInstruction)

// ---------------------------------------------------------------------------------------------------

type ChooseSquareInstruction struct{}

// String implements GameInstruction
func (*ChooseSquareInstruction) String() string {
	return "choose square"
}

// ToCommand implements GameInstruction
func (*ChooseSquareInstruction) ToCommand(gc *GameClient, content any) (Command, error) {
	position, err := game.ParsePosition(content)
	if err != nil {
		return nil, err
	}

	return &ChooseSquareCommand{gc, position}, nil
}

var _ GameInstruction = new(ChooseSquareInstruction)

// ---------------------------------------------------------------------------------------------------

type SelectCharacterInstruction struct{}

// String implements GameInstruction
func (*SelectCharacterInstruction) String() string {
	return "select character"
}

// ToCommand implements GameInstruction
func (*SelectCharacterInstruction) ToCommand(gc *GameClient, content any) (Command, error) {
	character, err := game.ParseSquareCharacter(content)
	if err != nil {
		return nil, err
	}

	return &SelectCharacterCommand{gc, character}, nil
}

var _ GameInstruction = new(SelectCharacterInstruction)
