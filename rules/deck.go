package rules

import "github.com/derekimcheng/mj/domain"

// Deck is a collection of Tiles from which tiles can be drawn from.
type Deck interface {
	// IsEmpty returns true if the Deck is empty.
	IsEmpty() bool
	// PopFront removes a tile from the front of the deck and returns it. If the deck is empty,
	// an error will be returned.
	PopFront() (*domain.Tile, error)
	// PopFront removes a tile from the back of the deck and returns it. If the deck is empty,
	// an error will be returned.
	PopBack() (*domain.Tile, error)
}
