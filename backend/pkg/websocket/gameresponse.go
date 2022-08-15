package websocket

import "encoding/json"

type GameResponse struct {
	Command Command `json:"command"`
	Body    any     `json:"body"`
}

type Command uint8

var _ json.Marshaler = new(Command) // Command implements Marshaler interface

const (
	BOARD Command = iota + 1
	GAME_OVER
	RESULT
)

var (
	CommandName = map[uint8]string{
		1: "board",
		2: "game over",
		3: "result",
	}
	CommandValue = map[string]uint8{
		"board":     1,
		"game over": 2,
		"result":    3,
	}
)

func (c Command) String() string {
	return CommandName[uint8(c)]
}

func (c Command) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}
