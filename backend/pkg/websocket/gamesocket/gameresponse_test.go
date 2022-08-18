package gamesocket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandString(t *testing.T) {
	type stringTestCases struct {
		input          Command
		expectedResult string
	}

	for _, testCase := range []stringTestCases{
		{
			input:          Command(1),
			expectedResult: "board",
		},
		{
			input:          Command(2),
			expectedResult: "game over",
		},
		{
			input:          Command(3),
			expectedResult: "result",
		},
	} {
		t.Run("Testing String method", func(t *testing.T) {
			actualResult := testCase.input.String()
			assert.Equal(t, testCase.expectedResult, actualResult)
		})
	}
}

func TestCommandMarshal(t *testing.T) {
	if _, err := BOARD.MarshalJSON(); err != nil {
		t.Fatal(err)
	}
}
