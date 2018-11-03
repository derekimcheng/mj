package domain

import (
	"sort"
	"fmt"
)

// Hand represents a collection of Tiles in a player's hand.
type Hand struct {
	tiles []*Tile
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
func (h* Hand) SetTiles(tiles []*Tile) {
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

// Sort sorts the hand using the default ordering.
func (h *Hand) Sort() {
	sort.Sort(h)
}

// Len ... (implements sort.Interface)
func (h *Hand) Len() int {
	return len(h.tiles)
}

// Swap ... (implements sort.Interface)
func (h *Hand) Swap(i, j int) {
	h.tiles[i], h.tiles[j] = h.tiles[j], h.tiles[i]
}

// Less ... (implements sort.Interface)
func (h *Hand) Less(i, j int) bool {
	tile1 := h.tiles[i]
	tile2 := h.tiles[j]
	return CompareTiles(tile1, tile2)
}

// String ...
func (h *Hand) String() string {
	ret := ""
	for _, t := range h.tiles {
		ret += "[" + t.String() + "]"
	}
	return ret
}
