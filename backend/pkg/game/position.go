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
)

func ParsePosition(p any) (Position, error) {
	if p, ok := p.(int); !ok {
		return Position(0), fmt.Errorf("invalid position: %v", p)
	} else {
		if p < 0 || 8 < p {
			return Position(0), fmt.Errorf("Position must be between 0 and 8")
		}
		return Position(p), nil
	}
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
