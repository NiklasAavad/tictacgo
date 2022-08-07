package websocket

type Board = [9]int
type WinningCombination = [3]int

type SquareCharacter int

const (
	X SquareCharacter = iota
	O
)

func (s SquareCharacter) String() string {
	switch s {
	case X:
		return "X"
	case O:
		return "O"
	}
	return "unknown"
}

type Result struct {
	WinningCombination WinningCombination
	WinningCharacter   SquareCharacter
}

type GameService interface {
	StartGame() Board
	GetResult() Result
	IsGameOver() bool
	IsChoiceValid(position int) bool
	ChooseSquare(position int)
	ChangePlayerInTurn()
}

// TODO bør nok have info om hvilke clients der er hvem
type Game struct {
	Board        Board
	PlayerInTurn SquareCharacter
}

var _ GameService = new(Game) // check that Game implements GameService

func (g Game) StartGame() Board {
	return [9]int{}
}

func (g Game) GetResult() Result {
	return Result{}
}

func (g Game) IsGameOver() bool {
	return false
}

func (g Game) IsChoiceValid(position int) bool {
	return false
}

func (g Game) ChooseSquare(position int) {
	// not implemented
}

func (g Game) ChangePlayerInTurn() {
	// not implemented
}
