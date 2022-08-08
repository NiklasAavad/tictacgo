package websocket

type InfoSender int

const (
	GAME_INFO InfoSender = iota
	CHAT_INFO
)

func (i InfoSender) String() string {
	switch i {
	case GAME_INFO:
		return "Game Info"
	case CHAT_INFO:
		return "Chat Info"
	}
	return "Unknown Sender"
}
