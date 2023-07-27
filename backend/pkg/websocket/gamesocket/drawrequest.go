package gamesocket

type DrawRequestHandler struct {
	IsDrawRequested bool
	DrawRequester  *GameClient
}

func NewDrawRequestHandler() *DrawRequestHandler {
	return &DrawRequestHandler{
		IsDrawRequested: false,
		DrawRequester: nil,
	}
}
