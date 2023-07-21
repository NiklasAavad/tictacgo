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
			input:          new(StartGameInstruction),
			expectedResult: "start game",
		},
		{
			input:          new(ChooseSquareInstruction),
			expectedResult: "choose square",
		},
		{
			input:          new(SelectCharacterInstruction),
			expectedResult: "select character",
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
	if gameInstruction.String() != "start game" {
		t.Errorf("expected START_GAME, got %v", gameInstruction)
	}
}

func TestParseChooseSquare(t *testing.T) {
	gameInstruction, err := ParseGameInstruction("choose square")
	if err != nil {
		t.Fatal(err)
	}
	if gameInstruction.String() != "choose square" {
		t.Errorf("expected CHOOSE_SQUARE, got %v", gameInstruction)
	}
}

func TestParseIsCaseInsensitive(t *testing.T) {
	gameInstruction, err := ParseGameInstruction("sTaRt gAME")
	if err != nil {
		t.Fatal(err)
	}
	if gameInstruction.String() != "start game" {
		t.Errorf("expected START_GAME, got %v", gameInstruction)
	}
}

func TestParseTrimsBothSides(t *testing.T) {
	gameInstruction, err := ParseGameInstruction("    start game        ")
	if err != nil {
		t.Fatal(err)
	}
	if gameInstruction.String() != "start game" {
		t.Errorf("expected START_GAME, got %v", gameInstruction)
	}
}

func TestParseShouldThrowError(t *testing.T) {
	_, err := ParseGameInstruction("should throw error")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestUnmarshalSucces(t *testing.T) {
	var giParser GameInstructionParser

	input := []byte("\"start game\"")

	if err := giParser.UnmarshalJSON(input); err != nil {
		t.Fatal(err)
	}

	if giParser.GameInstruction.String() != "start game" {
		t.Errorf("expected START_GAME, got %v", giParser.GameInstruction)
	}
}

func TestUnmarshalFailure(t *testing.T) {
	var giParser GameInstructionParser

	input := []byte("Cannot unmarshal this")

	err := giParser.UnmarshalJSON(input)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestUnmarshalFailureParsing(t *testing.T) {
	var giParser GameInstructionParser

	input := []byte("\"wrong string\"")

	err := giParser.UnmarshalJSON(input)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
