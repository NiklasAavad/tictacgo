package game

import "encoding/json"

type SquareCharacter uint8

var _ json.Marshaler = new(SquareCharacter)

const (
	X SquareCharacter = iota + 1
	O
	EMPTY
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
