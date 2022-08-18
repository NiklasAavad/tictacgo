package gamesocket

import (
	"testing"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket"
)

func setupTest(t *testing.T) (func(t *testing.T), websocket.Pool) {
	t.Log("Setting up testing")

	gamePool := NewGamePool()

	return func(t *testing.T) {
		t.Log("Tearing down testing")
	}, gamePool
}

func TestRegisterCharacter(t *testing.T) {
	teardown, _ := setupTest(t)
	defer teardown(t)
}
