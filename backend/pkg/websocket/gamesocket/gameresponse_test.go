package gamesocket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseTypeString(t *testing.T) {
	type stringTestCases struct {
		input          ResponseType
		expectedResult string
	}

	for _, testCase := range []stringTestCases{
		{
			input:          ResponseType(1),
			expectedResult: "board",
		},
		{
			input:          ResponseType(2),
			expectedResult: "game over",
		},
		{
			input:          ResponseType(3),
			expectedResult: "result",
		},
		{
			input:          ResponseType(4),
			expectedResult: "new message",
		},
		{
			input:          ResponseType(5),
			expectedResult: "character selected",
		},
		{
			input:          ResponseType(6),
			expectedResult: "game started",
		},
		{
			input:          ResponseType(7),
			expectedResult: "welcome",
		},
		{
			input:          ResponseType(8),
			expectedResult: "request draw",
		},
		{
			input:          ResponseType(9),
			expectedResult: "draw request response",
		},
		{
			input:          ResponseType(10),
			expectedResult: "withdraw draw request",
		},
	} {
		t.Run("Testing String method", func(t *testing.T) {
			actualResult := testCase.input.String()
			assert.Equal(t, testCase.expectedResult, actualResult)
		})
	}
}

func TestResponseTypeMarshall(t *testing.T) {
	if _, err := BOARD.MarshalJSON(); err != nil {
		t.Fatal(err)
	}
}
