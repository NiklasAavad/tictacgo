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

func createSelectCharacterMessage(client *GameClient, character game.SquareCharacter) GameMessage {
	instruction := &SelectCharacterInstruction{
		character: character,
	}
	return GameMessage{instruction, client}
}

func createStartGameMessage(client *GameClient) GameMessage {
	instruction := &StartGameInstruction{}
	return GameMessage{instruction, client}
}

func createChooseSquareMessage(client *GameClient, position game.Position) GameMessage {
	instruction := &ChooseSquareInstruction{
		position: position,
	}
	return GameMessage{instruction, client}
}

func createRequestDrawMessage(client *GameClient) GameMessage {
	instruction := &RequestDrawInstruction{}
	return GameMessage{instruction, client}
}

func createBothClients(pool *GamePool) (*GameClient, *GameClient) {
	return createTestClient(pool), createTestClient(pool)
}

func selectBothCharacters(pool *GamePool, clientX, clientO *GameClient) error {
	if err := pool.Broadcast(createSelectCharacterMessage(clientX, game.X)); err != nil {
		return err
	}

	if err := pool.Broadcast(createSelectCharacterMessage(clientO, game.O)); err != nil {
		return err
	}

	return nil
}

func startGame(pool *GamePool, client *GameClient) error {
	message := createStartGameMessage(client)
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
	message := createChooseSquareMessage(c, position)
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
