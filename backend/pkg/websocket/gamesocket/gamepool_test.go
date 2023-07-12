package gamesocket

import (
	"testing"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket"
	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket/testutils"
)

func setupTest(t *testing.T) (func(t *testing.T), *GamePool) {
	t.Log("Setting up testing")

	pool := NewGamePool(NewSequentialChannelStrategy())

	return func(t *testing.T) {
		t.Log("Tearing down testing")
	}, pool
}

func createTestClient(pool *GamePool) *GameClient {
	client := &GameClient{
		Pool: pool,
		conn: testutils.NewConnMock(),
	}

	pool.Register(client)

	return client
}

func createSelectCharacterMessage(client websocket.Client, character game.SquareCharacter) GameMessage {
	giParser := GameInstructionParser{
		gi: &SelectCharacterInstruction{},
	}
	return GameMessage{giParser, character.String(), client}
}

func createStartGameMessage(client websocket.Client) GameMessage {
	giParser := GameInstructionParser{
		gi: &StartGameInstruction{},
	}
	return GameMessage{giParser, 0, client}
}

func createChooseSquareMessage(client websocket.Client, position game.Position) GameMessage {
	giParser := GameInstructionParser{
		gi: &ChooseSquareInstruction{},
	}
	return GameMessage{giParser, position, client}
}

func startGame(pool *GamePool) (websocket.Client, websocket.Client, error) {
	clientX, clientO := createTestClient(pool), createTestClient(pool)

	messageX := createSelectCharacterMessage(clientX, game.X)
	messageO := createSelectCharacterMessage(clientO, game.O)

	if err := pool.Broadcast(messageX); err != nil {
		return nil, nil, err
	}

	if err := pool.Broadcast(messageO); err != nil {
		return nil, nil, err
	}

	message := createStartGameMessage(clientX)
	if err := pool.Broadcast(message); err != nil {
		return nil, nil, err
	}

	return clientX, clientO, nil
}

func chooseSquare(pool *GamePool, c websocket.Client, position game.Position) error {
	message := createChooseSquareMessage(c, position)
	return pool.Broadcast(message)
}

func TestRegisterClient(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	clientsInPool := len(pool.Clients())
	if clientsInPool != 0 {
		t.Fatalf("Expected no clients in pool, got %d", clientsInPool)
	}

	createTestClient(pool)
	clientsInPool = len(pool.Clients())
	if clientsInPool != 1 {
		t.Fatalf("Expected 1 client in pool, got %d", clientsInPool)
	}

	createTestClient(pool)
	clientsInPool = len(pool.Clients())
	if clientsInPool != 2 {
		t.Log(pool.Clients())
		t.Fatalf("Expected 2 clients in pool, got %d", clientsInPool)
	}
}

func TestUnregisterClient(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	createTestClient(pool)
	createTestClient(pool)

	clientToUnregister := createTestClient(pool)
	pool.Unregister(clientToUnregister)

	clientsInPool := len(pool.Clients())
	if clientsInPool != 2 {
		t.Fatalf("Expected 2 clients in pool, got %d", clientsInPool)
	}
}

func TestShouldBeCharacterX(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	client := createTestClient(pool)

	message := createSelectCharacterMessage(client, game.X)

	if err := pool.Broadcast(message); err != nil {
		t.Fatal(err)
	}

	if pool.xClient != client {
		t.Errorf("Expected client to be xClient, but got %v", pool.xClient)
	}
}

func TestShouldBeCharacterO(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	client := createTestClient(pool)

	message := createSelectCharacterMessage(client, game.O)

	if err := pool.Broadcast(message); err != nil {
		t.Fatal(err)
	}

	if pool.oClient != client {
		t.Errorf("Expected client to be oClient, but got %v", pool.oClient)
	}
}

func TestShouldNotChangeCharacterIfCharacterIsAlreadyTaken(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	client1 := createTestClient(pool)
	message1 := createSelectCharacterMessage(client1, game.X)
	if err := pool.Broadcast(message1); err != nil {
		t.Fatal(err)
	}

	client2 := createTestClient(pool)
	message2 := createSelectCharacterMessage(client2, game.X)
	if err := pool.Broadcast(message2); err != nil {
		t.Fatal(err)
	}

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

	client := createTestClient(pool)

	message := createStartGameMessage(client)
	if err := pool.Broadcast(message); err != nil {
		t.Fatal(err)
	}

	if pool.game.IsStarted() {
		t.Errorf("Game should still not have started, despite the Start Game message, since both characters must've been selected")
	}
}

func TestGameShouldBeAbleToStartWhenBothCharactersAreSelected(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	client1 := createTestClient(pool)
	message1 := createSelectCharacterMessage(client1, game.X)
	if err := pool.Broadcast(message1); err != nil {
		t.Fatal(err)
	}

	client2 := createTestClient(pool)
	message2 := createSelectCharacterMessage(client2, game.O)
	if err := pool.Broadcast(message2); err != nil {
		t.Fatal(err)
	}

	message := createStartGameMessage(client1)
	if err := pool.Broadcast(message); err != nil {
		t.Fatal(err)
	}

	if !pool.game.IsStarted() {
		t.Errorf("Game should be started")
	}
}

func TestOCannotChooseSquareWhenItIsX(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	_, clientO, err := startGame(pool)
	if err != nil {
		t.Fatal(err)
	}

	originalBoard := pool.game.Board()

	if err := chooseSquare(pool, clientO, game.CENTER); err != nil {
		t.Fatal(err)
	}

	want := originalBoard
	got := pool.game.Board()

	if want != got {
		t.Errorf("Wanted %v, got %v", want, got)
	}
}

func TestXCannotChooseSquareWhenItIsO(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	clientX, _, err := startGame(pool)
	if err != nil {
		t.Fatal(err)
	}

	if err := chooseSquare(pool, clientX, game.CENTER); err != nil {
		t.Fatal(err)
	}

	boardAfterFirstPlay := pool.game.Board()

	if err := chooseSquare(pool, clientX, game.BOTTOM_CENTER); err != nil {
		t.Fatal(err)
	}

	want := boardAfterFirstPlay
	got := pool.game.Board()

	if want != got {
		t.Errorf("Wanted %v, got %v", want, got)
	}
}

func TestSpectatorCannotStartGame(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	clientX := createTestClient(pool)
	clientO := createTestClient(pool)
	clientSpectator := createTestClient(pool)

	messageX := createSelectCharacterMessage(clientX, game.X)
	messageO := createSelectCharacterMessage(clientO, game.O)

	if err := pool.Broadcast(messageX); err != nil {
		t.Fatal(err)
	}

	if err := pool.Broadcast(messageO); err != nil {
		t.Fatal(err)
	}

	startGameMessage := createStartGameMessage(clientSpectator)
	if err := pool.Broadcast(startGameMessage); err != nil {
		t.Fatal(err)
	}

	if pool.game.IsStarted() {
		t.Errorf("Game should not have started, since it was started by the spectator")
	}
}

func TestClientCannotChooseBothCharacters(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	client := createTestClient(pool)

	messageX := createSelectCharacterMessage(client, game.X)
	if err := pool.Broadcast(messageX); err != nil {
		t.Fatal(err)
	}

	if pool.xClient != client {
		t.Errorf("Client should have selected X")
	}

	messageO := createSelectCharacterMessage(client, game.O)
	if err := pool.Broadcast(messageO); err != nil {
		t.Fatal(err)
	}

	if pool.oClient == client {
		t.Errorf("Client should not be able to select O, when they already selected X")
	}

	if pool.xClient != client {
		t.Errorf("Client should still have selected X")
	}
}
