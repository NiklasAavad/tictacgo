package gamesocket

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket"
)

type GameMessage struct {
	Instruction GameInstruction `json:"instruction"`
	Content     any             `json:"content"`
	Client      websocket.Client
}

type GameInstruction uint8

var _ json.Unmarshaler = new(GameInstruction) // Unmarshaller interface is implemented

const (
	START_GAME GameInstruction = iota + 1
	CHOOSE_SQUARE
	GET_BOARD
	SELECT_CHARACTER
)

var (
	GameInstructionName = map[uint8]string{
		1: "start game",
		2: "choose square",
		3: "get board",
		4: "select character",
	}
	GameInstructionValue = map[string]uint8{
		"start game":       1,
		"choose square":    2,
		"get board":        3,
		"select character": 4,
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
