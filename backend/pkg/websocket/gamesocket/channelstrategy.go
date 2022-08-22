package gamesocket

import (
	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket"
)

type ChannelStrategy interface {
	broadcast(p *GamePool, c Command)
	register(p *GamePool, c websocket.Client)
	unregister(p *GamePool, c websocket.Client)
}

type ConcurrentChannelStrategy struct{}

func NewConcurrentChannelStrategy() *ConcurrentChannelStrategy {
	return &ConcurrentChannelStrategy{}
}

// broadcast implements ChannelStrategy
func (*ConcurrentChannelStrategy) broadcast(p *GamePool, c Command) {
	p.broadcast <- c
}

// register implements ChannelStrategy
func (*ConcurrentChannelStrategy) register(p *GamePool, c websocket.Client) {
	p.register <- c
}

// unregister implements ChannelStrategy
func (*ConcurrentChannelStrategy) unregister(p *GamePool, c websocket.Client) {
	p.unregister <- c
}

var _ ChannelStrategy = new(ConcurrentChannelStrategy)

type SequentialChannelStrategy struct{}

func NewSequentialChannelStrategy() ChannelStrategy {
	return &SequentialChannelStrategy{}
}

// broadcast implements ChannelStrategy
func (*SequentialChannelStrategy) broadcast(p *GamePool, c Command) {
	p.respond(c)
}

// register implements ChannelStrategy
func (*SequentialChannelStrategy) register(p *GamePool, c websocket.Client) {
	p.clients[c] = game.EMPTY
}

// unregister implements ChannelStrategy
func (*SequentialChannelStrategy) unregister(p *GamePool, c websocket.Client) {
	delete(p.clients, c)
}

var _ ChannelStrategy = new(SequentialChannelStrategy)
