package gamesocket

import (
	"testing"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
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
			input:          &ChooseSquareInstruction{position: game.TOP_RIGHT},
			expectedResult: "choose square: TOP_RIGHT",
		},
		{
			input:          &SelectCharacterInstruction{character: game.O},
			expectedResult: "select character: O",
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
	if gameInstruction.String() != "choose square: NO_POSITION" {
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
	type marshallTestCases struct {
		input       []byte
		instruction string
	}

	for _, testCase := range []marshallTestCases{
		{
			input:       []byte("{\"instruction\":\"start game\"}"),
			instruction: "start game",
		},
		{
			input:       []byte("{\"instruction\":\"select character\",\"content\":\"X\"}"),
			instruction: "select character: X",
		},
		{
			input:       []byte("{\"instruction\":\"choose square\",\"content\":4}"),
			instruction: "choose square: CENTER",
		},
	} {
		t.Run("Testing unmarshal of game messages", func(t *testing.T) {
			var msg GameMessage
			err := msg.UnmarshalJSON(testCase.input)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, testCase.instruction, msg.Instruction.String())
		})
	}
}

func TestUnmarshalFailure(t *testing.T) {
	var msg GameMessage

	input := []byte("Cannot unmarshal this")

	err := msg.UnmarshalJSON(input)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestUnmarshalFailureParsing(t *testing.T) {
	var msg GameMessage

	input := []byte("\"wrong string\"")

	err := msg.UnmarshalJSON(input)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestChooseSquareInstructionShouldHaveNoPositionDefault(t *testing.T) {
	chooseSquareInstruction := NewChooseSquareInstruction()

	want := game.NO_POSITION
	got := chooseSquareInstruction.position

	if want != got {
		t.Errorf("Expected %v, got %v", want, got)
	}
}

func TestSelectCharacterInstructionShouldHaveEmptyCharacterDefault(t *testing.T) {
	selectCharacterInstruction := NewSelectCharacterInstruction()

	want := game.EMPTY_CHARACTER
	got := selectCharacterInstruction.character

	if want != got {
		t.Errorf("Expected %v, got %v", want, got)
	}
}

func getSelectCharacterInstructionFromParseInstruction() *SelectCharacterInstruction {
	gameInstruction, _ := ParseGameInstruction("select character")
	return gameInstruction.(*SelectCharacterInstruction)
}

func TestNewGameInstructionsAreCreatedAtParser(t *testing.T) {
	firstSelectCharacterInstruction := getSelectCharacterInstructionFromParseInstruction()
	firstSelectCharacterInstruction.character = game.X

	secondSelectCharacterInstruction := getSelectCharacterInstructionFromParseInstruction()

	want := game.EMPTY_CHARACTER
	got := secondSelectCharacterInstruction.character

	if want != got {
		t.Errorf("Every instance must be new instance. Expected %v, got %v", want, got)
	}
}
