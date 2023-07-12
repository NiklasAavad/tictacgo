package gamesocket

import (
	"fmt"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket"
)

type GameMessage struct {
	InstructionParser GameInstructionParser `json:"instruction"`
	Content           any                   `json:"content"`
	Client            websocket.Client
}

func (gm *GameMessage) ToCommand() (Command, error) {
	gameClient, ok := gm.Client.(*GameClient)
	if !ok {
		return nil, fmt.Errorf("invalid client type: %T", gm.Client)
	}

	// TODO should we check if there is a gi in the parser?
	command, err := gm.InstructionParser.gi.ToCommand(gameClient, gm.Content)
	if err != nil {
		return nil, err
	}
	return command, nil
}
