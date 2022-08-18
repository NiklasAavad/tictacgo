package websocket

type Conn interface {
	Close() error
	ReadJSON(v interface{}) error
	ReadMessage() (messageType int, p []byte, err error)
	WriteJSON(v interface{}) error
}
