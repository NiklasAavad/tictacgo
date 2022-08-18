package testutils

import "github.com/NiklasPrograms/tictacgo/backend/pkg/websocket"

type ClientMock struct {
	conn websocket.Conn
	name string
}

var _ websocket.Client = new(ClientMock)

func NewClientMock() *ClientMock {
	return &ClientMock{
		conn: NewConnMock(),
	}
}

// Conn implements websocket.Client
func (c *ClientMock) Conn() websocket.Conn {
	return c.conn
}

// Name implements websocket.Client
func (c *ClientMock) Name() string {
	return c.name
}

// Read implements websocket.Client
func (*ClientMock) Read() {
	panic("unimplemented")
}
