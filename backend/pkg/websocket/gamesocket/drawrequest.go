package gamesocket

import "errors"

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

// Common validation for DrawRequests
//
// Ensures that the client is actually playing
// and that the game is in a valid state (started, but not finished)
func IsDrawRequestCommandValid(client *GameClient) (bool, error) {
	pool := client.Pool
	game := pool.game

	clientIsPlaying := pool.xClient == client || pool.oClient == client
	if !clientIsPlaying {
		return false, errors.New("Client must be playing to accept a draw, cannot be a spectator. Indicates a larger problem with the server")
	}

	if game.IsGameOver() {
		return false, errors.New("Game is not over, and a draw cannot be requested")
	}

	if !game.IsStarted() {
		return false, errors.New("Game has not started yet, and a draw cannot be requested")
	}

	return true, nil
}

