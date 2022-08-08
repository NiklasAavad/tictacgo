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
	assert.Equal(t, expectedPlayerInTurn, g.PlayerInTurn)
}
