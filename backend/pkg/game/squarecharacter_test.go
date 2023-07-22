package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSquareCharacterString(t *testing.T) {
	type stringTestCases struct {
		input          SquareCharacter
		expectedResult string
	}

	for _, testCase := range []stringTestCases{
		{
			input:          SquareCharacter(1),
			expectedResult: "X",
		},
		{
			input:          SquareCharacter(2),
			expectedResult: "O",
		},
		{
			input:          SquareCharacter(3),
			expectedResult: "",
		},
		{
			input:          SquareCharacter(4),
			expectedResult: "unknown",
		},
	} {
		t.Run("Testing String method", func(t *testing.T) {
			actualResult := testCase.input.String()
			assert.Equal(t, testCase.expectedResult, actualResult)
		})
	}
}

func TestMarshalX(t *testing.T) {
	_, err := X.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
}

func TestMarshalO(t *testing.T) {
	_, err := O.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
}

func TestMarshalEMPTY(t *testing.T) {
	_, err := EMPTY_CHARACTER.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
}
