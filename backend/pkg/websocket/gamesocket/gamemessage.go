package gamesocket

import (
	"encoding/json"
	"errors"
)

type GameMessage struct {
	Instruction GameInstruction
	Client      *GameClient
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *GameMessage) UnmarshalJSON(bytes []byte) error {
	var data struct {
		Instruction string `json:"instruction"`
		Content     any    `json:"content"`
	}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	parsedGameInstruction, err := ParseGameInstruction(data.Instruction)
	if err != nil {
		return err
	}
	msg.Instruction = parsedGameInstruction

	if err := msg.Instruction.ParseContent(data.Content); err != nil {
		return err
	}

	return nil
}

func (msg *GameMessage) ToCommand() (Command, error) {
	if msg.Client == nil {
		return nil, errors.New("No GameClient found in GameMessage")
	}

	if msg.Instruction == nil {
		return nil, errors.New("No GameInstruction found in GameMessage")
	}

	return msg.Instruction.ToCommand(msg.Client)
}

// assert that GameMessage implements unmarshaller interface
var _ json.Unmarshaler = new(GameMessage)
