package websocket

import (
	"net/http"
)

type Pool interface {
	Register(c Client)
	NewClient(w http.ResponseWriter, r *http.Request) Client
}
