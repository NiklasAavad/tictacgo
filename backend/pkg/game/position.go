package game

import (
	"encoding/json"
	"fmt"
)

type Position int

var _ json.Unmarshaler = new(Position)

const (
	TOP_LEFT Position = iota
	TOP_CENTER
	TOP_RIGHT
	CENTER_LEFT
	CENTER
	CENTER_RIGHT
	BOTTOM_LEFT
	BOTTOM_CENTER
	BOTTOM_RIGHT
	NO_POSITION Position = -1
)

func castToFloat(input any) (float64, bool) {
	switch k := input.(type) {
	case float64:
		return k, true
	case int:
		return float64(k), true
	case Position:
		return float64(k), true
	default:
		return -1, false
	}
}

func ParsePosition(input any) (Position, error) {
	p, ok := castToFloat(input)

	if !ok {
		return Position(0), fmt.Errorf("invalid type for position, expected number: %v", input)
	}

	if p < 0 || 8 < p {
		return Position(0), fmt.Errorf("Position must be between 0 and 8")
	}

	return Position(int(p)), nil
}

// TODO slet eventuelt, hvis en GameMessage ikke længere indeholder en Position, men kan være hvad som helst.
func (p *Position) UnmarshalJSON(data []byte) (err error) {
	var position int
	if err = json.Unmarshal(data, &position); err != nil {
		return err
	}
	if *p, err = ParsePosition(position); err != nil {
		return err
	}
	return nil
}

func (p Position) String() string {
	switch p {
	case TOP_LEFT:
		return "TOP_LEFT"
	case TOP_CENTER:
		return "TOP_CENTER"
	case TOP_RIGHT:
		return "TOP_RIGHT"
	case CENTER_LEFT:
		return "CENTER_LEFT"
	case CENTER:
		return "CENTER"
	case CENTER_RIGHT:
		return "CENTER_RIGHT"
	case BOTTOM_LEFT:
		return "BOTTOM_LEFT"
	case BOTTOM_CENTER:
		return "BOTTOM_CENTER"
	case BOTTOM_RIGHT:
		return "BOTTOM_RIGHT"
	case NO_POSITION:
		return "NO_POSITION"
	default:
		return "unknown"
	}
}
