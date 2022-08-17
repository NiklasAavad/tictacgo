package game

import (
	"testing"
)

func TestPositionShouldBeParsedForZero(t *testing.T) {
	p, err := ParsePosition(0)
	if err != nil {
		t.Fatal(err)
	}
	if p != TOP_LEFT {
		t.Errorf("expected position to be TOP_LEFT, got %v", p)
	}
}

func TestPositionShouldBeParsedForEight(t *testing.T) {
	p, err := ParsePosition(8)
	if err != nil {
		t.Fatal(err)
	}
	if p != BOTTOM_RIGHT {
		t.Errorf("expected position to be BOTTOM_RIGHT, got %v", p)
	}
}

func TestShouldFailIfInputLowerThanZero(t *testing.T) {
	if _, err := ParsePosition(-1); err == nil {
		t.Errorf("Expected error with input position %d", -1)
	}
}

func TestShouldFailIfInputHeigherThanEight(t *testing.T) {
	if _, err := ParsePosition(9); err == nil {
		t.Errorf("Expected error with input position %d", 9)
	}
}

func TestUnmarshalSucces(t *testing.T) {
	var p Position
	input := []byte{0x30}

	if err := p.UnmarshalJSON(input); err != nil {
		t.Fatal(err)
	}

	if p != TOP_LEFT {
		t.Errorf("expected position to be TOP_LEFT, got %v", p)
	}
}

func TestUnmarshalFailure(t *testing.T) {
	var p Position
	input := []byte{0x90}

	err := p.UnmarshalJSON(input)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestUnmarshalFailureParsing(t *testing.T) {
	var p Position
	input := []byte{0x39}

	err := p.UnmarshalJSON(input)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
