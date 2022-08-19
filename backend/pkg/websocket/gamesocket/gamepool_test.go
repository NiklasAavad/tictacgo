package gamesocket

import (
	"testing"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket"
	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket/testutils"
)

func setupTest(t *testing.T) (func(t *testing.T), websocket.Pool) {
	t.Log("Setting up testing")

	pool := NewGamePool(NewSequentialChannelStrategy())

	return func(t *testing.T) {
		t.Log("Tearing down testing")
	}, pool
}

func createNewClient(p websocket.Pool) websocket.Client {
	client := testutils.NewClientMock()
	p.Register(client)
	return client
}

func TestRegisterClient(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	clientsInPool := len(pool.Clients())
	if clientsInPool != 0 {
		t.Fatalf("Expected no clients in pool, got %d", clientsInPool)
	}

	createNewClient(pool)
	clientsInPool = len(pool.Clients())
	if clientsInPool != 1 {
		t.Fatalf("Expected 1 client in pool, got %d", clientsInPool)
	}

	createNewClient(pool)
	clientsInPool = len(pool.Clients())
	if clientsInPool != 2 {
		t.Log(pool.Clients())
		t.Fatalf("Expected 2 clients in pool, got %d", clientsInPool)
	}
}

func TestUnregisterClient(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	createNewClient(pool)
	createNewClient(pool)

	clientToUnregister := createNewClient(pool)
	pool.Unregister(clientToUnregister)

	clientsInPool := len(pool.Clients())
	if clientsInPool != 2 {
		t.Fatalf("Expected 2 clients in pool, got %d", clientsInPool)
	}
}
