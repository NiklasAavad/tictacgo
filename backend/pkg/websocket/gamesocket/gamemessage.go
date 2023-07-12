package gamesocket

type GameMessage struct {
	InstructionParser GameInstructionParser `json:"instruction"`
	Content           any                   `json:"content"`
	Client            *GameClient
}

func (gm *GameMessage) ToCommand() (Command, error) {
	// TODO should we check if there is a gi in the parser?
	command, err := gm.InstructionParser.gi.ToCommand(gm.Client, gm.Content)
	if err != nil {
		return nil, err
	}
	return command, nil
}
