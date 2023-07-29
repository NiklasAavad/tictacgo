package gamesocket

import (
	"testing"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
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
	err := pool.Unregister(clientToUnregister)
	if err != nil {
		t.Fatal(err)
	}

	clientsInPool := len(pool.Clients())
	if clientsInPool != 2 {
		t.Fatalf("Expected 2 clients in pool, got %d", clientsInPool)
	}
}

func TestShouldBeCharacterX(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	client := createTestClient(pool)

	command := createSelectCharacterCommand(client, game.X)

	if err := pool.Broadcast(command); err != nil {
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

	command := createSelectCharacterCommand(client, game.O)

	if err := pool.Broadcast(command); err != nil {
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
	command1 := createSelectCharacterCommand(client1, game.X)
	if err := pool.Broadcast(command1); err != nil {
		t.Fatal(err)
	}

	client2 := createTestClient(pool)
	command2 := createSelectCharacterCommand(client2, game.X)
	if err := pool.Broadcast(command2); err == nil { // should return error!
		t.Fatal(err)
	}

	want := client1
	got := pool.xClient

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

	command := createStartGameCommand(client)
	if err := pool.Broadcast(command); err == nil { // should return error!
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
	command1 := createSelectCharacterCommand(client1, game.X)
	if err := pool.Broadcast(command1); err != nil {
		t.Fatal(err)
	}

	client2 := createTestClient(pool)
	command2 := createSelectCharacterCommand(client2, game.O)
	if err := pool.Broadcast(command2); err != nil {
		t.Fatal(err)
	}

	command := createStartGameCommand(client1)
	if err := pool.Broadcast(command); err != nil {
		t.Fatal(err)
	}

	if !pool.game.IsStarted() {
		t.Errorf("Game should be started")
	}
}

func TestOCannotChooseSquareWhenItIsX(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	_, clientO, err := initGame(pool)
	if err != nil {
		t.Fatal(err)
	}

	originalBoard := pool.game.Board()

	if err := chooseSquare(pool, clientO, game.CENTER); err == nil { // should return error!
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

	clientX, _, err := initGame(pool)
	if err != nil {
		t.Fatal(err)
	}

	if err := chooseSquare(pool, clientX, game.CENTER); err != nil {
		t.Fatal(err)
	}

	boardAfterFirstPlay := pool.game.Board()

	if err := chooseSquare(pool, clientX, game.BOTTOM_CENTER); err == nil { // should return error!
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

	messageX := createSelectCharacterCommand(clientX, game.X)
	messageO := createSelectCharacterCommand(clientO, game.O)

	if err := pool.Broadcast(messageX); err != nil {
		t.Fatal(err)
	}

	if err := pool.Broadcast(messageO); err != nil {
		t.Fatal(err)
	}

	startGameCommand := createStartGameCommand(clientSpectator)
	if err := pool.Broadcast(startGameCommand); err == nil { // should return error!
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

	messageX := createSelectCharacterCommand(client, game.X)
	if err := pool.Broadcast(messageX); err != nil {
		t.Fatal(err)
	}

	if pool.xClient != client {
		t.Errorf("Client should have selected X")
	}

	messageO := createSelectCharacterCommand(client, game.O)
	if err := pool.Broadcast(messageO); err == nil { // should return error!
		t.Fatal(err)
	}

	if pool.oClient == client {
		t.Errorf("Client should not be able to select O, when they already selected X")
	}

	if pool.xClient != client {
		t.Errorf("Client should still have selected X")
	}
}

func TestClientCanRequestDraw(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	clientX, _, err := initGame(pool)
	if err != nil {
		t.Fatal(err)
	}

	if pool.DrawRequestHandler.IsDrawRequested {
		t.Errorf("Game should not have a draw requested")
	}

	requestDrawCommand := createRequestDrawCommand(clientX)
	if err := pool.Broadcast(requestDrawCommand); err != nil {
		t.Fatal(err)
	}

	if !pool.DrawRequestHandler.IsDrawRequested {
		t.Errorf("Game should have a draw requested")
	}
}

func TestClientCannotRequestDrawWhenGameIsNotStarted(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	clientX, clientO := createBothClients(pool)
	if err := selectBothCharacters(pool, clientX, clientO); err != nil {
		t.Fatal(err)
	}

	requestDrawCommand := createRequestDrawCommand(clientX)
	if err := pool.Broadcast(requestDrawCommand); err == nil { // should return error!
		t.Fatalf("Broadcast should fail, since game is not started")
	}

	if pool.DrawRequestHandler.IsDrawRequested {
		t.Errorf("Game should not have a draw requested, since it was not started")
	}
}

func TestClientCannotRequestDrawWhenGameIsOver(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	clientX, clientO, err := initGame(pool)
	if err != nil {
		t.Fatal(err)
	}

	if err := playTillGameWon(pool, clientX, clientO); err != nil {
		t.Fatal(err)
	}

	requestDrawCommand := createRequestDrawCommand(clientO)
	if err := pool.Broadcast(requestDrawCommand); err == nil { // should return error!
		t.Errorf("Broadcast should fail, since the game is over")
	}

	if pool.DrawRequestHandler.IsDrawRequested {
		t.Errorf("Game should not have a draw requested, since the game is over")
	}
}

func TestClientCannotRequestDrawWhenDrawIsAlreadyRequested(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	clientX, clientO, err := initGame(pool)
	if err != nil {
		t.Fatal(err)
	}

	firstRequestDrawCommand := createRequestDrawCommand(clientX)
	if err := pool.Broadcast(firstRequestDrawCommand); err != nil {
		t.Fatal(err)
	}

	secondRequestDrawCommand := createRequestDrawCommand(clientO)
	if err := pool.Broadcast(secondRequestDrawCommand); err == nil { // should return error!
		t.Errorf("Broadcast should fail, since a draw is already requested")
	}
}

func TestSpectatorCannotRequestDraw(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	_, _, err := initGame(pool)
	if err != nil {
		t.Fatal(err)
	}

	spectator := createTestClient(pool)

	requestDrawCommand := createRequestDrawCommand(spectator)
	if err := pool.Broadcast(requestDrawCommand); err == nil { // should return error!
		t.Errorf("Broadcast should fail, since the spectator is not a player")
	}

	if pool.DrawRequestHandler.IsDrawRequested {
		t.Errorf("Game should not have a draw requested, since the spectator is not a player")
	}
}

func TestClientCanAcceptDraw(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	clientX, clientO, err := initGame(pool)
	if err != nil {
		t.Fatal(err)
	}

	requestDrawCommand := createRequestDrawCommand(clientX)
	if err := pool.Broadcast(requestDrawCommand); err != nil {
		t.Fatal(err)
	}

	acceptDrawCommand := createRespondToDrawRequestCommand(clientO, true)
	if err := pool.Broadcast(acceptDrawCommand); err != nil {
		t.Fatal(err)
	}

	if pool.DrawRequestHandler.IsDrawRequested {
		t.Errorf("Game should not have a draw requested")
	}

	if !pool.game.IsGameOver() {
		t.Errorf("Game should be over, since a draw was accepted")
	}

}

func TestClientCanDeclineDraw(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	clientX, clientO, err := initGame(pool)
	if err != nil {
		t.Fatal(err)
	}

	requestDrawCommand := createRequestDrawCommand(clientX)
	if err := pool.Broadcast(requestDrawCommand); err != nil {
		t.Fatal(err)
	}

	acceptDrawCommand := createRespondToDrawRequestCommand(clientO, false)
	if err := pool.Broadcast(acceptDrawCommand); err != nil {
		t.Fatal(err)
	}

	if pool.DrawRequestHandler.IsDrawRequested {
		t.Errorf("Game should not have a draw requested")
	}

	if pool.game.IsGameOver() {
		t.Errorf("Game should not be over, since a draw was declined")
	}
}

func TestClientCannotRespondToDrawRequestIfNoRequestIsActive(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	_, clientO, err := initGame(pool)
	if err != nil {
		t.Fatal(err)
	}

	acceptDrawCommand := createRespondToDrawRequestCommand(clientO, true)
	if err := pool.Broadcast(acceptDrawCommand); err == nil { // should return error!
		t.Errorf("Broadcast should fail, since no draw request is active")
	}
}

func TestClientCannotRespondToDrawRequestIfNotPlayer(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	clientX, _, err := initGame(pool)
	if err != nil {
		t.Fatal(err)
	}

	requestDrawCommand := createRequestDrawCommand(clientX)
	if err := pool.Broadcast(requestDrawCommand); err != nil {
		t.Fatal(err)
	}

	spectator := createTestClient(pool)

	acceptDrawCommand := createRespondToDrawRequestCommand(spectator, true)
	if err := pool.Broadcast(acceptDrawCommand); err == nil { // should return error!
		t.Errorf("Broadcast should fail, since spectator is not a player")
	}
}

func TestClientCannotRespondToDrawRequestIfGameIsOver(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	clientX, clientO, err := initGame(pool)
	if err != nil {
		t.Fatal(err)
	}

	requestDrawCommand := createRequestDrawCommand(clientO)
	if err := pool.Broadcast(requestDrawCommand); err != nil {
		t.Fatal(err)
	}

	if err := playTillGameWon(pool, clientX, clientO); err != nil {
		t.Fatal(err)
	}

	acceptDrawCommand := createRespondToDrawRequestCommand(clientX, true)
	if err := pool.Broadcast(acceptDrawCommand); err == nil { // should return error!
		t.Errorf("Broadcast should fail, since game is over")
	}
}

func TestClientCannotRespondToOwnDrawRequest(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	clientX, _, err := initGame(pool)
	if err != nil {
		t.Fatal(err)
	}

	requestDrawCommand := createRequestDrawCommand(clientX)
	if err := pool.Broadcast(requestDrawCommand); err != nil {
		t.Fatal(err)
	}

	acceptDrawCommand := createRespondToDrawRequestCommand(clientX, true)
	if err := pool.Broadcast(acceptDrawCommand); err == nil { // should return error!
		t.Errorf("Broadcast should fail, since clientX is the one who requested the draw")
	}

	if !pool.DrawRequestHandler.IsDrawRequested {
		t.Errorf("Game should not have a draw requested")
	}
}

func TestClientCanWithdrawDrawRequest(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	clientX, _, err := initGame(pool)
	if err != nil {
		t.Fatal(err)
	}

	requestDrawCommand := createRequestDrawCommand(clientX)
	if err := pool.Broadcast(requestDrawCommand); err != nil {
		t.Fatal(err)
	}

	if !pool.DrawRequestHandler.IsDrawRequested {
		t.Errorf("Game should have a draw requested")
	}

	withdrawDrawCommand := createWithdrawDrawRequestCommand(clientX)
	if err := pool.Broadcast(withdrawDrawCommand); err != nil {
		t.Fatal(err)
	}

	if pool.DrawRequestHandler.IsDrawRequested {
		t.Errorf("Game should NOT have a draw requested")
	}
}

func TestClientCannotWithdrawDrawRequestIfNoRequestIsActive(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	clientX, _, err := initGame(pool)
	if err != nil {
		t.Fatal(err)
	}

	withdrawDrawCommand := createWithdrawDrawRequestCommand(clientX)
	if err := pool.Broadcast(withdrawDrawCommand); err == nil { // should return error!
		t.Errorf("Broadcast should fail, since no draw request is active")
	}
}

func TestClientCannotWithdrawDrawRequestIfNotPlayer(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	clientX, _, err := initGame(pool)
	if err != nil {
		t.Fatal(err)
	}

	requestDrawCommand := createRequestDrawCommand(clientX)
	if err := pool.Broadcast(requestDrawCommand); err != nil {
		t.Fatal(err)
	}

	spectator := createTestClient(pool)

	withdrawDrawCommand := createWithdrawDrawRequestCommand(spectator)
	if err := pool.Broadcast(withdrawDrawCommand); err == nil { // should return error!
		t.Errorf("Broadcast should fail, since spectator is not a player")
	}
}

func TestClientCannotWithdrawDrawRequestIfGameIsOver(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	clientX, clientO, err := initGame(pool)
	if err != nil {
		t.Fatal(err)
	}

	requestDrawCommand := createRequestDrawCommand(clientX)
	if err := pool.Broadcast(requestDrawCommand); err != nil {
		t.Fatal(err)
	}

	if err := playTillGameWon(pool, clientX, clientO); err != nil {
		t.Fatal(err)
	}

	withdrawDrawCommand := createWithdrawDrawRequestCommand(clientX)
	if err := pool.Broadcast(withdrawDrawCommand); err == nil { // should return error!
		t.Errorf("Broadcast should fail, since game is over")
	}
}

func TestClientCannotWithdrawDrawRequestIfGameIsNotStarted(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	clientX, clientO := createBothClients(pool)
	err := selectBothCharacters(pool, clientX, clientO)
	if err != nil {
		t.Fatal(err)
	}

	pool.DrawRequestHandler.IsDrawRequested = true // simulate draw request, as a draw request should NOT be possible either before game starts

	withdrawDrawCommand := createWithdrawDrawRequestCommand(clientX)
	if err := pool.Broadcast(withdrawDrawCommand); err == nil { // should return error!
		t.Errorf("Broadcast should fail, since game is not started")
	}
}

func TestOnlyClientWhoSendDrawRequestCanWithdrawIt(t *testing.T) {
	teardown, pool := setupTest(t)
	defer teardown(t)

	clientX, clientO, err := initGame(pool)
	if err != nil {
		t.Fatal(err)
	}

	requestDrawCommand := createRequestDrawCommand(clientX)
	if err := pool.Broadcast(requestDrawCommand); err != nil {
		t.Fatal(err)
	}

	withdrawDrawCommand := createWithdrawDrawRequestCommand(clientO)
	if err := pool.Broadcast(withdrawDrawCommand); err == nil { // should return error!
		t.Errorf("Broadcast should fail, since only client who sent draw request can withdraw it")
	}
}
