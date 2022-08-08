package game

// TODO bør nok have info om hvilke clients der er hvem
type Game struct {
	board        Board
	playerInTurn SquareCharacter
}

func NewGame() *Game {
	return &Game{
		board:        newBoard(),
		playerInTurn: X,
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
	g.board = newBoard()
	g.playerInTurn = X
	return g.board
}

func (g *Game) GetResult() Result {
	for _, winningCombination := range WinningCombinations {
		if g.isWinningRow(winningCombination) {
			winningCharacter := g.board[winningCombination[0]]
			return Result{winningCombination, winningCharacter, true}
		}
	}

	return Result{HasWinner: false}
}

func (g *Game) occupiedSquares() []SquareCharacter {
	var slice []SquareCharacter
	for _, square := range g.board {
		if square != EMPTY {
			slice = append(slice, square)
		}
	}
	return slice
}

func (g *Game) notEnoughInputs() bool {
	return len(g.occupiedSquares()) < 5
}

func (g *Game) isBoardFull() bool {
	return len(g.occupiedSquares()) == 9
}

func (g *Game) isWinningRow(p [3]Position) bool {
	if g.board[p[0]] == EMPTY {
		return false
	}

	isSameCharacter := g.board[p[0]] == g.board[p[1]] && g.board[p[1]] == g.board[p[2]]
	return isSameCharacter
}

func (g *Game) hasWinner() bool {
	for _, winningCombination := range WinningCombinations {
		if g.isWinningRow(winningCombination) {
			return true
		}
	}
	return false
}

func (g *Game) IsGameOver() bool {
	if g.notEnoughInputs() {
		return false
	}

	if g.isBoardFull() {
		return true
	}

	return g.hasWinner()
}

func (g *Game) IsChoiceValid(p Position) bool {
	if g.IsGameOver() {
		return false
	}

	return g.board[p] == EMPTY
}

func (g *Game) ChooseSquare(p Position) Board {
	if g.IsChoiceValid(p) {
		g.board[p] = g.playerInTurn
		g.ChangePlayerInTurn()
	}

	return g.board
}

func (g *Game) ChangePlayerInTurn() SquareCharacter {
	if g.playerInTurn == X {
		g.playerInTurn = O
	} else if g.playerInTurn == O {
		g.playerInTurn = X
	}

	return g.playerInTurn
}

func (g *Game) Board() Board {
	return g.board
}
