package websocket

import (
	"net/http"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
)

type Pool interface {
	Register(c Client)
	Unregister(c Client)
	NewClient(w http.ResponseWriter, r *http.Request) Client
	Clients() map[Client]game.SquareCharacter
}
