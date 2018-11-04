package domain

import (
	"fmt"
	"sort"
)

// Hand represents a collection of Tiles in a player's hand.
type Hand struct {
	tiles Tiles
}

// NewHand returns a new empty hand.
func NewHand() *Hand {
	return &Hand{tiles: nil}
}

// GetTiles ...
func (h *Hand) GetTiles() []*Tile {
	return h.tiles
}

// SetTiles ...
func (h *Hand) SetTiles(tiles []*Tile) {
	h.tiles = tiles
}

// AddTile adds the given tile to the hand.
func (h *Hand) AddTile(t *Tile) {
	h.tiles = append(h.tiles, t)
}

// RemoveTile removes the tile at the given index and returns it. Returns an error if the index is
// out of bounds.
func (h *Hand) RemoveTile(index int) (*Tile, error) {
	if index < 0 || index >= len(h.tiles) {
		return nil, fmt.Errorf("Index %d out of bounds [0, %d)", index, len(h.tiles))
	}
	tile := h.tiles[index]
	h.tiles = append(h.tiles[:index], h.tiles[index+1:]...)
	return tile, nil
}

// NumTiles returns the number of tiles in the hand.
func (h *Hand) NumTiles() int {
	return len(h.tiles)
}

// Sort sorts the hand using the default ordering.
func (h *Hand) Sort() {
	sort.Sort(h.tiles)
}

// StringWithoutIndices returns a string representation of the hand without the indices.
func (h *Hand) StringWithoutIndices() string {
	ret := ""
	for _, t := range h.tiles {
		ret += t.String()
	}
	return ret
}

// String ...
func (h *Hand) String() string {
	ret := ""
	for i, t := range h.tiles {
		if i > 0 {
			ret += ", "
		}
		ret += fmt.Sprintf("%d: %s", i, t.String())
	}
	return ret
}
