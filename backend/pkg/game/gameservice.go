package game

type Board = [9]SquareCharacter
type WinningCombination = [3]Position

type Result struct {
	WinningCombination WinningCombination
	WinningCharacter   SquareCharacter
}

type GameService interface {
	StartGame() Board
	GetResult() Result
	IsGameOver() bool
	IsChoiceValid(p Position) bool
	ChooseSquare(p Position)
	ChangePlayerInTurn()
}
