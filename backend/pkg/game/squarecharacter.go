package game

import (
	"encoding/json"
	"fmt"
)

type SquareCharacter uint8

var _ json.Marshaler = new(SquareCharacter)

const (
	X SquareCharacter = iota + 1
	O
	EMPTY_CHARACTER
)

func (s SquareCharacter) String() string {
	switch s {
	case X:
		return "X"
	case O:
		return "O"
	case EMPTY_CHARACTER:
		return ""
	}
	return "unknown"
}

func (c SquareCharacter) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// TODO test
func ParseSquareCharacter(c any) (SquareCharacter, error) {
	c, ok := c.(string)
	if !ok {
		return SquareCharacter(0), fmt.Errorf("invalid type for squarecharacter, expected string: %v", c)
	}

	switch c {
	case "X":
		return X, nil
	case "O":
		return O, nil
	case "":
		return EMPTY_CHARACTER, nil
	default:
		return SquareCharacter(42), fmt.Errorf("%v is not a valid SquareCharacter", c)
	}
}
