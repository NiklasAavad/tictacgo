package gamesocket

import (
	"testing"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket/testutils"
)

func setupTest(t *testing.T) (func(t *testing.T), *GamePool) {
	t.Log("Setting up testing")

	pool := NewGamePool(NewSequentialChannelStrategy())

	return func(t *testing.T) {
		t.Log("Tearing down testing")
	}, pool
}

func TestRegisterClient(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	clientsInPool := len(pool.Clients())
	if clientsInPool != 0 {
		t.Fatalf("Expected no clients in pool, got %d", clientsInPool)
	}

	testutils.NewClientMock(pool)
	clientsInPool = len(pool.Clients())
	if clientsInPool != 1 {
		t.Fatalf("Expected 1 client in pool, got %d", clientsInPool)
	}

	testutils.NewClientMock(pool)
	clientsInPool = len(pool.Clients())
	if clientsInPool != 2 {
		t.Log(pool.Clients())
		t.Fatalf("Expected 2 clients in pool, got %d", clientsInPool)
	}
}

func TestUnregisterClient(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	testutils.NewClientMock(pool)
	testutils.NewClientMock(pool)

	clientToUnregister := testutils.NewClientMock(pool)
	pool.Unregister(clientToUnregister)

	clientsInPool := len(pool.Clients())
	if clientsInPool != 2 {
		t.Fatalf("Expected 2 clients in pool, got %d", clientsInPool)
	}
}

func TestShouldBeCharacterX(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	client := testutils.NewClientMock(pool)

	message := GameMessage{SELECT_CHARACTER, game.X.String(), client}

	pool.Broadcast(message)

	want := game.X
	got := pool.Clients()[client]

	if want != got {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestShouldBeCharacterO(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	client := testutils.NewClientMock(pool)

	message := GameMessage{SELECT_CHARACTER, game.O.String(), client}

	pool.Broadcast(message)

	want := game.O
	got := pool.Clients()[client]

	if want != got {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestShouldNotChangeCharacterIfCharacterIsAlreadyTaken(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	client1 := testutils.NewClientMock(pool)
	message1 := GameMessage{SELECT_CHARACTER, game.X.String(), client1}
	pool.Broadcast(message1)

	client2 := testutils.NewClientMock(pool)
	message2 := GameMessage{SELECT_CHARACTER, game.X.String(), client2}
	pool.Broadcast(message2)

	want := game.EMPTY
	got := pool.Clients()[client2]

	if want != got {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGameShouldNotStartWhenNoCharactersSelected(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	if pool.game.IsStarted() {
		t.Fatal("Game should not have started yet")
	}

	client := testutils.NewClientMock(pool)

	message := GameMessage{START_GAME, 0, client}
	pool.Broadcast(message)

	if pool.game.IsStarted() {
		t.Errorf("Game should still not have started, despite the Start Game message, since both characters must've been selected")
	}
}

func TestGameShouldBeAbleToStartWhenBothCharactersAreSelected(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	client1 := testutils.NewClientMock(pool)
	message1 := GameMessage{SELECT_CHARACTER, game.X.String(), client1}
	pool.Broadcast(message1)

	client2 := testutils.NewClientMock(pool)
	message2 := GameMessage{SELECT_CHARACTER, game.O.String(), client2}
	pool.Broadcast(message2)

	message := GameMessage{START_GAME, 0, client1}
	pool.Broadcast(message)

	if !pool.game.IsStarted() {
		t.Errorf("Game should be started")
	}
}
