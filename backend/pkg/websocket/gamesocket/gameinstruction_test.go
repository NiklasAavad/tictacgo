package gamesocket

import (
	"testing"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
)

func TestParseContentSucces(t *testing.T) {
	type testCases struct {
		input   GameInstruction
		content any
	}

	for _, testCase := range []testCases{
		{
			input:   NewStartGameInstruction(),
			content: nil,
		},
		{
			input:   NewChooseSquareInstruction(),
			content: game.TOP_RIGHT,
		},
		{
			input:   NewSelectCharacterInstruction(),
			content: game.O,
		},
		{
			input:   NewRequestDrawInstruction(),
			content: nil,
		},
		{
			input:   NewRespondToDrawInstruction(),
			content: true,
		},
	} {
		t.Run("Testing ParseContent method", func(t *testing.T) {
			err := testCase.input.ParseContent(testCase.content)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestRespondToDrawInstructionParseContentError(t *testing.T) {
	respondToDrawInstruction := NewRespondToDrawInstruction()
	err := respondToDrawInstruction.ParseContent("not a bool")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
