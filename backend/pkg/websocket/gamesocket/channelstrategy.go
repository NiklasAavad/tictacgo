package gamesocket

import (
	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket"
)

// ChannelStrategy defines the interface for the different channel strategies that can be used in the GamePool.
//
// Different strategies can be used for different purposes.
// For example, the SequentialChannelStrategy is used in tests, because it is easier to debug.
// The ConcurrentChannelStrategy is used in production, because it is more performant.
type ChannelStrategy interface {
	broadcast(p *GamePool, c Command) error
	register(p *GamePool, c websocket.Client)
	unregister(p *GamePool, c websocket.Client)
}

// ConcurrentChannelStrategy is a ChannelStrategy that uses a channel to broadcast messages to all clients.
// This is the most performant strategy, but it is hard to debug, and is therefore not used in tests.
type ConcurrentChannelStrategy struct{}

func NewConcurrentChannelStrategy() *ConcurrentChannelStrategy {
	return &ConcurrentChannelStrategy{}
}

// broadcast implements ChannelStrategy
func (*ConcurrentChannelStrategy) broadcast(p *GamePool, c Command) error {
	p.broadcast <- c
	return nil
}

// register implements ChannelStrategy
func (*ConcurrentChannelStrategy) register(p *GamePool, c websocket.Client) {
	p.register <- c
}

// unregister implements ChannelStrategy
func (*ConcurrentChannelStrategy) unregister(p *GamePool, c websocket.Client) {
	p.unregister <- c
}

// Assert that ConcurrentChannelStrategy implements ChannelStrategy
var _ ChannelStrategy = new(ConcurrentChannelStrategy)

// SequentialChannelStrategy is a ChannelStrategy that does not use channels to broadcast messages to all clients.
// Instead it uses the GamePool's clients map directly to send messages to all clients.
// This is not as performant as the ConcurrentChannelStrategy, but it is easier to debug, and is therefore used in tests.
type SequentialChannelStrategy struct{}

// NewSequentialChannelStrategy creates a new SequentialChannelStrategy
// This is not as performant as the ConcurrentChannelStrategy, but it is easier to debug, and is therefore used in tests.
func NewSequentialChannelStrategy() ChannelStrategy {
	return &SequentialChannelStrategy{}
}

// broadcast implements ChannelStrategy
func (*SequentialChannelStrategy) broadcast(p *GamePool, c Command) error {
	if err := p.respond(c); err != nil {
		return err
	}
	return nil
}

// register implements ChannelStrategy
func (*SequentialChannelStrategy) register(p *GamePool, c websocket.Client) {
	p.clients[c] = game.EMPTY
}

// unregister implements ChannelStrategy
func (*SequentialChannelStrategy) unregister(p *GamePool, c websocket.Client) {
	delete(p.clients, c)
}

// Assert that SequentialChannelStrategy implements ChannelStrategy
var _ ChannelStrategy = new(SequentialChannelStrategy)
