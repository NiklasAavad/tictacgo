package websocket

// Unused interface, which should be implemented in ChatPool and GamePool to reduce dry code.
type Pool interface {
	Register() chan *Client
	Unregister() chan *Client
	Broadcast() chan Message
}
