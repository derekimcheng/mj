package domain

import "fmt"

// Tile represents a tile in MJ. It contains information about the suit and the value of the tile.
type Tile struct {
	suit    *Suit
	ordinal int
	// id distinguishes tiles with the same tuit and ordinal, e.g., there may be 4 "Bamboo 1"
	// tiles, but they will have IDs [0, 1, 2, 3]. This field is typically not shown to the player.
	// TODO: move this?
	// id int
}

// NewTile returns a new Tile with the input parameters, or nil if the input is invalid.
func NewTile(suit *Suit, ordinal int, id int) (*Tile, error) {
	if suit == nil {
		return nil, fmt.Errorf("Suit cannot be nil")
	}
	if ordinal < 0 || ordinal >= suit.GetSize() {
		return nil, fmt.Errorf("Ordinal out of range [%d, %d): %d", 0, suit.GetSize(), ordinal)
	}
	return &Tile{suit: suit, ordinal: ordinal}, nil
}

// GetSuit ...
func (t *Tile) GetSuit() *Suit {
	return t.suit
}

// GetOrdinal ...
func (t *Tile) GetOrdinal() int {
	return t.ordinal
}
