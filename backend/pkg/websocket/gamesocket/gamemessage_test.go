package gamesocket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameInstructionString(t *testing.T) {
	type stringTestCases struct {
		input          GameInstruction
		expectedResult string
	}

	for _, testCase := range []stringTestCases{
		{
			input:          GameInstruction(1),
			expectedResult: "start game",
		},
		{
			input:          GameInstruction(2),
			expectedResult: "choose square",
		},
		{
			input:          GameInstruction(3),
			expectedResult: "get board",
		},
	} {
		t.Run("Testing String method", func(t *testing.T) {
			actualResult := testCase.input.String()
			assert.Equal(t, testCase.expectedResult, actualResult)
		})
	}
}

func TestParseStartGame(t *testing.T) {
	gameInstruction, err := ParseGameInstruction("start game")
	if err != nil {
		t.Fatal(err)
	}
	if gameInstruction != START_GAME {
		t.Errorf("expected START_GAME, got %v", gameInstruction)
	}
}

func TestParseChooseSquare(t *testing.T) {
	gameInstruction, err := ParseGameInstruction("choose square")
	if err != nil {
		t.Fatal(err)
	}
	if gameInstruction != CHOOSE_SQUARE {
		t.Errorf("expected CHOOSE_SQUARE, got %v", gameInstruction)
	}
}

func TestParseGetBoard(t *testing.T) {
	gameInstruction, err := ParseGameInstruction("get board")
	if err != nil {
		t.Fatal(err)
	}
	if gameInstruction != GET_BOARD {
		t.Errorf("expected GET_BOARD, got %v", gameInstruction)
	}
}

func TestParseIsCaseInsensitive(t *testing.T) {
	gameInstruction, err := ParseGameInstruction("sTaRt gAME")
	if err != nil {
		t.Fatal(err)
	}
	if gameInstruction != START_GAME {
		t.Errorf("expected START_GAME, got %v", gameInstruction)
	}
}

func TestParseTrimsBothSides(t *testing.T) {
	gameInstruction, err := ParseGameInstruction("    start game        ")
	if err != nil {
		t.Fatal(err)
	}
	if gameInstruction != START_GAME {
		t.Errorf("expected START_GAME, got %v", gameInstruction)
	}
}

func TestParseShouldThrowError(t *testing.T) {
	_, err := ParseGameInstruction("should throw error")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
