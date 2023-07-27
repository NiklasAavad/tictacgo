package gamesocket

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
)

type GameInstruction interface {
	String() string

	// ParseContent parses the content of the instruction, and returns an error if the content is invalid. The content is then stored in the instruction as a field.
	ParseContent(any) error

	ToCommand(*GameClient) (Command, error)
}

func ParseGameInstruction(s string) (GameInstruction, error) {
	s = strings.TrimSpace(strings.ToLower(s))

	switch s {
	case "start game":
		return NewStartGameInstruction(), nil
	case "choose square":
		return NewChooseSquareInstruction(), nil
	case "select character":
		return NewSelectCharacterInstruction(), nil
	case "request draw":
		return NewRequestDrawInstruction(), nil
	case "respond to draw":
		return NewRespondToDrawInstruction(), nil
	case "withdraw draw request":
		return NewWithdrawDrawRequestInstruction(), nil
	default:
		return nil, fmt.Errorf("invalid game instruction: %s", s)
	}
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
func (ins *ChooseSquareInstruction) String() string {
	return fmt.Sprintf("choose square: %v", ins.position.String())
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
func (ins *SelectCharacterInstruction) String() string {
	return fmt.Sprintf("select character: %v", ins.character)
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

// ---------------------------------------------------------------------------------------------------

type RequestDrawInstruction struct{}

func NewRequestDrawInstruction() *RequestDrawInstruction {
	return &RequestDrawInstruction{}
}

// String implements GameInstruction
func (*RequestDrawInstruction) String() string {
	return "request draw"
}

// ParseContent implements GameInstruction
func (*RequestDrawInstruction) ParseContent(any) error {
	return nil // RequestDrawInstruction has no content
}

// ToCommand implements GameInstruction
func (*RequestDrawInstruction) ToCommand(gc *GameClient) (Command, error) {
	return NewRequestDrawCommand(gc)
}

// ---------------------------------------------------------------------------------------------------

type RespondToDrawRequestInstruction struct {
	accept bool
}

func NewRespondToDrawInstruction() *RespondToDrawRequestInstruction {
	return &RespondToDrawRequestInstruction{
		accept: false,
	}
}

// String implements GameInstruction
func (ins *RespondToDrawRequestInstruction) String() string {
	return fmt.Sprintf("respond to draw: %v", ins.accept)
}

// ParseContent implements GameInstruction
func (instruction *RespondToDrawRequestInstruction) ParseContent(content any) error {
	accept, ok := content.(bool)
	if !ok {
		return fmt.Errorf("invalid content type for RespondToDrawInstruction: %v", reflect.TypeOf(content))
	}

	instruction.accept = accept
	return nil
}

// ToCommand implements GameInstruction
func (instruction *RespondToDrawRequestInstruction) ToCommand(gc *GameClient) (Command, error) {
	return NewRespondToDrawRequestCommand(gc, instruction.accept)
}

// ---------------------------------------------------------------------------------------------------

type WithdrawDrawRequestInstruction struct{}

func NewWithdrawDrawRequestInstruction() *WithdrawDrawRequestInstruction {
	return &WithdrawDrawRequestInstruction{}
}

// String implements GameInstruction
func (*WithdrawDrawRequestInstruction) String() string {
	return "withdraw draw request"
}

// ParseContent implements GameInstruction
func (*WithdrawDrawRequestInstruction) ParseContent(any) error {
	return nil // WithdrawDrawRequestInstruction has no content
}

// ToCommand implements GameInstruction
func (*WithdrawDrawRequestInstruction) ToCommand(gc *GameClient) (Command, error) {
	return NewWithdrawDrawRequestCommand(gc)
}
