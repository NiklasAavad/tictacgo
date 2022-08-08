package game

type Board = [9]SquareCharacter
type WinningCombination = [3]SquareCharacter

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
