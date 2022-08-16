package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartGame(t *testing.T) {
	g := NewGame()

	board := g.StartGame()
	expectedBoard := [9]SquareCharacter{
		EMPTY, EMPTY, EMPTY,
		EMPTY, EMPTY, EMPTY,
		EMPTY, EMPTY, EMPTY,
	}

	expectedPlayerInTurn := X

	assert.Equal(t, expectedBoard, board)
	assert.Equal(t, expectedBoard, g.board)
	assert.Equal(t, expectedPlayerInTurn, g.playerInTurn)
}

func TestShouldBeValidChoiceIfFreePosition(t *testing.T) {
	g := NewGame()

	assert.True(t, g.isChoiceValid(TOP_LEFT))
}

func TestShouldBeInvalidChoiceIfOccupiedPosition(t *testing.T) {
	g := NewGame()

	g.ChooseSquare(TOP_LEFT)

	assert.False(t, g.isChoiceValid(TOP_LEFT))
}

func TestShouldBeInvalidChoiceIfGameIsOver(t *testing.T) {
	g := NewGame()

	g.board = mockWinningBoard()

	assert.False(t, g.isChoiceValid(BOTTOM_RIGHT))
}

func TestShouldChooseTopLeft(t *testing.T) {
	g := NewGame()

	g.ChooseSquare(TOP_LEFT)

	expectedBoard := [9]SquareCharacter{
		X, EMPTY, EMPTY,
		EMPTY, EMPTY, EMPTY,
		EMPTY, EMPTY, EMPTY,
	}

	assert.Equal(t, expectedBoard, g.board)
}

func TestChangePlayerInTurn(t *testing.T) {
	g := NewGame()

	g.changePlayerInTurn()
	assert.Equal(t, O, g.playerInTurn)

	g.changePlayerInTurn()
	assert.Equal(t, X, g.playerInTurn)
}

func mockWinningBoard() Board {
	return [9]SquareCharacter{
		X, X, X,
		O, O, EMPTY,
		EMPTY, EMPTY, EMPTY,
	}
}

func mockAlmostFullBoard() Board {
	return [9]SquareCharacter{
		X, X, O,
		O, O, X,
		EMPTY, EMPTY, EMPTY,
	}
}

func TestGameOver(t *testing.T) {
	g := NewGame()
	assert.False(t, g.IsGameOver())

	g.board = mockAlmostFullBoard()
	assert.False(t, g.IsGameOver())

	g.board = mockWinningBoard()
	assert.True(t, g.IsGameOver())
}

func mockFullBoard() Board {
	return [9]SquareCharacter{
		X, O, O,
		O, X, X,
		X, X, O,
	}
}

func TestShouldBeGameOverIfAllPositionsOccupied(t *testing.T) {
	g := NewGame()

	assert.False(t, g.IsGameOver())

	g.board = mockFullBoard()
	assert.True(t, g.IsGameOver())
}

func TestShouldBeWinningResult(t *testing.T) {
	g := NewGame()

	g.board = mockWinningBoard()

	winningCombination := [3]Position{TOP_LEFT, TOP_CENTER, TOP_RIGHT}
	expectedResult := Result{winningCombination, X, true}

	assert.Equal(t, expectedResult, g.GetResult())
}

func TestShouldBeEmptyResult(t *testing.T) {
	g := NewGame()

	g.board = mockFullBoard()

	expectedResult := Result{}
	assert.Equal(t, expectedResult, g.GetResult())
}

func TestShouldGetBoard(t *testing.T) {
	g := NewGame()

	g.board = mockFullBoard()

	assert.Equal(t, mockFullBoard(), g.Board())
}
