package game

type Board = [9]SquareCharacter
type WinningCombination = [3]Position

type Result struct {
	WinningCombination WinningCombination
	WinningCharacter   SquareCharacter
	HasWinner          bool
}

type GameService interface {
	StartGame() Board
	GetResult() Result
	IsGameOver() bool
	IsChoiceValid(p Position) bool
	ChooseSquare(p Position) Board
	ChangePlayerInTurn() SquareCharacter
}
