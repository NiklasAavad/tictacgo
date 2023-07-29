package gamesocket

import (
	"testing"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
)

func TestParseContentSucces(t *testing.T) {
	type testCases struct {
		input   Command
		content any
	}

	for _, testCase := range []testCases{
		{
			input:   &StartGameCommand{},
			content: nil,
		},
		{
			input:   &ChooseSquareCommand{},
			content: game.TOP_RIGHT,
		},
		{
			input:   &SelectCharacterCommand{},
			content: "X",
		},
		{
			input:   &RequestDrawCommand{},
			content: nil,
		},
		{
			input:   &RespondToDrawRequestCommand{},
			content: true,
		},
		{
			input:   &WithdrawDrawRequestCommand{},
			content: nil,
		},
	} {
		t.Run("Testing ParseContent method", func(t *testing.T) {
			err := testCase.input.parseContent(testCase.content)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestRespondToDrawInstructionParseContentError(t *testing.T) {
	respondToDrawInstruction := &RespondToDrawRequestCommand{}
	err := respondToDrawInstruction.parseContent("not a bool")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
