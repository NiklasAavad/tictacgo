package game

type Board = [9]SquareCharacter
type WinningCombination = [3]Position

type Result struct {
	WinningCombination WinningCombination
	WinningCharacter   SquareCharacter
	HasWinner          bool
}

type GameService interface {
	StartGame()
	GetResult() Result
	IsGameOver() bool
	isChoiceValid(p Position) bool
	ChooseSquare(p Position) (Board, error)
	changePlayerInTurn() SquareCharacter
	Board() Board
	IsStarted() bool
	PlayerInTurn() SquareCharacter
	ForceDraw()
}
