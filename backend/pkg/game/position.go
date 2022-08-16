package game

import (
	"encoding/json"
	"fmt"
)

type Position uint8

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

func ParsePosition(p int8) (Position, error) {
	if p < 0 || 8 < p {
		return Position(0), fmt.Errorf("Position must be between 0 and 8")
	}
	return Position(p), nil
}

func (p *Position) UnmarshalJSON(data []byte) (err error) {
	var position int8
	if err = json.Unmarshal(data, &position); err != nil {
		return err
	}
	if *p, err = ParsePosition(position); err != nil {
		return err
	}
	return nil
}
