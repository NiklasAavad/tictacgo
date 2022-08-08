package game

type SquareCharacter int

type Board = [9]SquareCharacter
type WinningCombination = [3]SquareCharacter

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

func NewGame() *Game {
	return &Game{
		Board:        newBoard(),
		PlayerInTurn: X,
	}
}

var _ GameService = new(Game) // check that Game implements GameService

func newBoard() Board {
	return [9]SquareCharacter{
		EMPTY, EMPTY, EMPTY,
		EMPTY, EMPTY, EMPTY,
		EMPTY, EMPTY, EMPTY,
	}
}

func (g *Game) StartGame() Board {
	g.Board = newBoard()
	g.PlayerInTurn = X
	return g.Board
}

func (g *Game) GetResult() Result {
	return Result{}
}

func (g *Game) IsGameOver() bool {
	return false
}

func (g *Game) IsChoiceValid(position int) bool {
	return false
}

func (g *Game) ChooseSquare(position int) {
	// not implemented
}

func (g *Game) ChangePlayerInTurn() {
	// not implemented
}
