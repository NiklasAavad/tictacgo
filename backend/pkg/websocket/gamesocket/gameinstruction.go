package gamesocket

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
)

type GameInstruction interface {
	String() string
	ParseContent(any) error
	ToCommand(*GameClient) (Command, error)
}

// ---------------------------------------------------------------------------------------------------

type GameInstructionParser struct {
	GameInstruction GameInstruction
}

var _ json.Unmarshaler = new(GameInstructionParser)

// map from strings to GameInstruction
var GameInstructionValue = map[string]GameInstruction{
	"start game":       new(StartGameInstruction),
	"choose square":    new(ChooseSquareInstruction),
	"select character": new(SelectCharacterInstruction),
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

	parser.GameInstruction = gi

	return nil
}

// ---------------------------------------------------------------------------------------------------

type StartGameInstruction struct{}

func NewStartGameInstruction() *StartGameInstruction {
	return &StartGameInstruction{}
}

// String implements GameInstruction
func (*StartGameInstruction) String() string {
	return "start game"
}

// ParseContent implements GameInstruction
func (*StartGameInstruction) ParseContent(any) error {
	// TODO consider if we should check that content is nil / 0 / whatever
	return nil // StartGameInstruction has no content
}

// ToCommand implements GameInstruction
func (*StartGameInstruction) ToCommand(gc *GameClient) (Command, error) {
	return NewStartGameCommand(gc)
}

var _ GameInstruction = new(StartGameInstruction)

// ---------------------------------------------------------------------------------------------------

type ChooseSquareInstruction struct {
	position game.Position
}

func NewChooseSquareInstruction() *ChooseSquareInstruction {
	return &ChooseSquareInstruction{
		position: game.NO_POSITION,
	}
}

// String implements GameInstruction
func (*ChooseSquareInstruction) String() string {
	return "choose square"
}

// ParseContent implements GameInstruction
func (instruction *ChooseSquareInstruction) ParseContent(content any) error {
	position, err := game.ParsePosition(content)
	if err != nil {
		return err
	}

	instruction.position = position
	return nil
}

// ToCommand implements GameInstruction
func (instruciton *ChooseSquareInstruction) ToCommand(gc *GameClient) (Command, error) {
	return NewChooseSquareCommand(gc, instruciton.position)
}

var _ GameInstruction = new(ChooseSquareInstruction)

// ---------------------------------------------------------------------------------------------------

type SelectCharacterInstruction struct {
	character game.SquareCharacter
}

func NewSelectCharacterInstruction() *SelectCharacterInstruction {
	return &SelectCharacterInstruction{
		character: game.EMPTY_CHARACTER,
	}
}

// String implements GameInstruction
func (*SelectCharacterInstruction) String() string {
	return "select character"
}

// ParseContent implements GameInstruction
func (instruction *SelectCharacterInstruction) ParseContent(content any) error {
	character, err := game.ParseSquareCharacter(content)
	if err != nil {
		return err
	}

	instruction.character = character
	return nil
}

// ToCommand implements GameInstruction
func (instruction *SelectCharacterInstruction) ToCommand(gc *GameClient) (Command, error) {
	return NewSelectCharacterCommand(gc, instruction.character)
}

var _ GameInstruction = new(SelectCharacterInstruction)
