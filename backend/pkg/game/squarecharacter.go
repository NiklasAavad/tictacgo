package game

type SquareCharacter int

const (
	X SquareCharacter = iota
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
