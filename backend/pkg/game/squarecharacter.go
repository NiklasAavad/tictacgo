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
	EMPTY
)

var (
	SquareCharacterName = map[uint8]string{
		1: "X",
		2: "O",
		3: "",
	}
	SquareCharacterValue = map[string]uint8{
		"X": 1,
		"O": 2,
		"":  3,
	}
)

func (s SquareCharacter) String() string {
	switch s {
	case X:
		return "X"
	case O:
		return "O"
	case EMPTY:
		return ""
	}
	return "unknown"
}

func (c SquareCharacter) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// TODO test
func ParseSquareCharacter(c any) (SquareCharacter, error) {
	if c, ok := c.(string); !ok {
		return SquareCharacter(0), fmt.Errorf("invalid type for squarecharacter, expected string: %v", c)
	} else {
		value, ok := SquareCharacterValue[c]
		if !ok {
			return SquareCharacter(0), fmt.Errorf("%v is not a valid SquareCharacter", c)
		}
		return SquareCharacter(value), nil
	}
}
