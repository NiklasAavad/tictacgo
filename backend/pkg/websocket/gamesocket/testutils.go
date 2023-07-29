package gamesocket

import (
	"fmt"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket"
)

type ConnMock struct{}

var _ websocket.Conn = new(ConnMock)

func NewConnMock() *ConnMock {
	return &ConnMock{}
}

// Close implements websocket.Conn
func (*ConnMock) Close() error {
	panic("unimplemented")
}

// ReadJSON implements websocket.Conn
func (*ConnMock) ReadJSON(v interface{}) error {
	panic("unimplemented")
}

// ReadMessage implements websocket.Conn
func (*ConnMock) ReadMessage() (messageType int, p []byte, err error) {
	panic("unimplemented")
}

// WriteJSON implements websocket.Conn
func (*ConnMock) WriteJSON(v interface{}) error {
	return nil
}

// -----------------------------------------------------------------------------------------------

func createTestClient(pool *GamePool) *GameClient {
	client := &GameClient{
		Pool: pool,
		conn: NewConnMock(),
	}

	pool.Register(client)

	return client
}

func createSelectCharacterCommand(client *GameClient, character game.SquareCharacter) Command {
	return &SelectCharacterCommand{client, character}
}

func createStartGameCommand(client *GameClient) Command {
	return &StartGameCommand{client}
}

func createChooseSquareCommand(client *GameClient, position game.Position) Command {
	return &ChooseSquareCommand{client, position}
}

func createRequestDrawCommand(client *GameClient) Command {
	return &RequestDrawCommand{client}
}

func createRespondToDrawRequestCommand(client *GameClient, accept bool) Command {
	return &RespondToDrawRequestCommand{client, accept}
}

func createWithdrawDrawRequestCommand(client *GameClient) Command {
	return &WithdrawDrawRequestCommand{client}
}

func createBothClients(pool *GamePool) (*GameClient, *GameClient) {
	return createTestClient(pool), createTestClient(pool)
}

func selectBothCharacters(pool *GamePool, clientX, clientO *GameClient) error {
	if err := pool.Broadcast(createSelectCharacterCommand(clientX, game.X)); err != nil {
		return err
	}

	if err := pool.Broadcast(createSelectCharacterCommand(clientO, game.O)); err != nil {
		return err
	}

	return nil
}

func startGame(pool *GamePool, client *GameClient) error {
	message := createStartGameCommand(client)
	return pool.Broadcast(message)
}

func initGame(pool *GamePool) (*GameClient, *GameClient, error) {
	clientX, clientO := createBothClients(pool)

	if err := selectBothCharacters(pool, clientX, clientO); err != nil {
		return nil, nil, err
	}

	if err := startGame(pool, clientX); err != nil {
		return nil, nil, err
	}

	return clientX, clientO, nil
}

func chooseSquare(pool *GamePool, c *GameClient, position game.Position) error {
	message := createChooseSquareCommand(c, position)
	return pool.Broadcast(message)
}

func playTillGameWon(pool *GamePool, clientX, clientO *GameClient) error {
	if err := chooseSquare(pool, clientX, game.CENTER); err != nil {
		return err
	}

	if err := chooseSquare(pool, clientO, game.TOP_LEFT); err != nil {
		return err
	}

	if err := chooseSquare(pool, clientX, game.TOP_RIGHT); err != nil {
		return err
	}

	if err := chooseSquare(pool, clientO, game.BOTTOM_RIGHT); err != nil {
		return err
	}

	if err := chooseSquare(pool, clientX, game.BOTTOM_LEFT); err != nil {
		return err
	}

	if !pool.game.IsGameOver() {
		return fmt.Errorf("Game should be over")
	}

	return nil
}
