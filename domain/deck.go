package domain

import (
	"fmt"
	"math/rand"
)

// Deck is a collection of Tiles from which tiles can be drawn from.
type Deck interface {
	// Size returns the number of remaining tiles in the deck.
	NumRemainingTiles() int
	// IsEmpty returns true if the Deck is empty.
	IsEmpty() bool
	// Shuttle randomly shuffles the tiles in the deck. It is assumed that the random source has
	// been properly seeded.
	Shuffle()
	// PopFront removes a tile from the front of the deck and returns it. If the deck is empty,
	// an error will be returned.
	PopFront() (*Tile, error)
	// PopFront removes a tile from the back of the deck and returns it. If the deck is empty,
	// an error will be returned.
	PopBack() (*Tile, error)
}

// SliceDeck is an implementation of Deck using slices.
type SliceDeck struct {
	tiles []*Tile
}

// NumRemainingTiles .. (Deck implementation)
func (d *SliceDeck) NumRemainingTiles() int {
	return len(d.tiles)
}

// IsEmpty ... (Deck implementation)
func (d *SliceDeck) IsEmpty() bool {
	return len(d.tiles) == 0
}

// Shuffle ... (Deck implementation)
func (d *SliceDeck) Shuffle() {
	rand.Shuffle(len(d.tiles), func(i, j int) {
		d.tiles[i], d.tiles[j] = d.tiles[j], d.tiles[i]
	})
}

// PopFront ... (Deck implementation)
func (d *SliceDeck) PopFront() (*Tile, error) {
	if d.IsEmpty() {
		return nil, fmt.Errorf("Deck is empty")
	}
	tile := d.tiles[0]
	d.tiles = d.tiles[1:]
	return tile, nil
}

// PopBack ... (Deck implementation)
func (d *SliceDeck) PopBack() (*Tile, error) {
	if d.IsEmpty() {
		return nil, fmt.Errorf("Deck is empty")
	}
	tile := d.tiles[len(d.tiles)-1]
	d.tiles = d.tiles[:len(d.tiles)-1]
	return tile, nil
}

func newSliceDeck(tiles []*Tile) *SliceDeck {
	return &SliceDeck{tiles: tiles}
}

// NewDeck creates and returns a new Deck given input tiles.
func NewDeck(tiles []*Tile) Deck {
	return newSliceDeck(tiles)
}
