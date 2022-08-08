package game

type Position int

const (
	TOP_LEFT Position = iota
	TOP_CENTER
	TOP_RIGHT
	CENTER_LEFT
	CENTER
	CENTER_RIGHT
	BOTTOM_LEFT
	BOTTOM_CENTER
	BOTTOM_RIGHT
)