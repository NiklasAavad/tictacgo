package testutils

import "github.com/NiklasPrograms/tictacgo/backend/pkg/websocket"

type ClientMock struct {
	conn websocket.Conn
}

var _ websocket.Client = new(ClientMock)

func NewClientMock(pool websocket.Pool) *ClientMock {
	client := &ClientMock{
		conn: NewConnMock(),
	}

	pool.Register(client)

	return client
}

// Conn implements websocket.Client
func (c *ClientMock) Conn() websocket.Conn {
	return c.conn
}

// Name implements websocket.Client
func (c *ClientMock) Name() string {
	panic("unimplemented")
}

// Read implements websocket.Client
func (*ClientMock) Read() {
	panic("unimplemented")
}
