package gamesocket

import "errors"

type GameMessage struct {
	InstructionParser GameInstructionParser `json:"instruction"`
	Content           any                   `json:"content"`
	Client            *GameClient
}

func (msg *GameMessage) ToCommand() (Command, error) {
	if msg.Client == nil {
		return nil, errors.New("No GameClient found in GameMessage")
	}

	instruction := msg.InstructionParser.GameInstruction
	if instruction == nil {
		return nil, errors.New("No GameInstruction found in GameMessage")
	}

	if err := instruction.ParseContent(msg.Content); err != nil {
		return nil, err
	}

	return instruction.ToCommand(msg.Client)
}
