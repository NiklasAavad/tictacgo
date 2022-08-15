package gamesocket

import (
	"encoding/json"
	"fmt"
	"strings"
)

type GameInstruction uint8

var _ json.Unmarshaler = new(GameInstruction) // Unmarshaller interface is implemented

const (
	START_GAME GameInstruction = iota + 1
	CHOOSE_SQUARE
	GET_BOARD
)

var (
	GameInstructionName = map[uint8]string{
		1: "start game",
		2: "choose square",
		3: "get board",
	}
	GameInstructionValue = map[string]uint8{
		"start game":    1,
		"choose square": 2,
		"get board":     3,
	}
)

func (gi GameInstruction) String() string {
	return GameInstructionName[uint8(gi)]
}

func ParseGameInstruction(s string) (GameInstruction, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := GameInstructionValue[s]
	if !ok {
		return GameInstruction(0), fmt.Errorf("%q is not a valid Game Instruction", s)
	}
	return GameInstruction(value), nil
}

func (gi *GameInstruction) UnmarshalJSON(data []byte) (err error) {
	var gameInstruction string
	if err := json.Unmarshal(data, &gameInstruction); err != nil {
		return err
	}
	if *gi, err = ParseGameInstruction(gameInstruction); err != nil {
		return err
	}
	return nil
}
