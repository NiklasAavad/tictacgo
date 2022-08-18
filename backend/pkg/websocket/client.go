package websocket

type Client interface {
	Conn() Conn
	Name() string
	Read()
}
