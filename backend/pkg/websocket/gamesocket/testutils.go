package gamesocket
import "github.com/NiklasPrograms/tictacgo/backend/pkg/websocket"

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
