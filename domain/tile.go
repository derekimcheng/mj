package domain

import (
	"strings"
	"fmt"
)

// Tile contains an ID in addition to all of the fields in TileBase.
type Tile struct {
	TileBase
	// id distinguishes tiles with the same tuit and ordinal, e.g., there may be 4 "Bamboo 1"
	// tiles, but they will have IDs [0, 1, 2, 3]. This field is typically not shown to the player.
	id int
}

// TileBase contains information about the suit and the value of the tile.
type TileBase struct {
	suit *Suit
	ordinal int
}

// NewTileBase returns a new TileBase with the given parameters.
func NewTileBase(suit *Suit, ordinal int) TileBase {
	return TileBase{suit, ordinal}
}

// GetSuit ...
func (tb TileBase) GetSuit() *Suit {
	return tb.suit
}

// GetOrdinal ...
func (tb TileBase) GetOrdinal() int {
	return tb.ordinal
}

// NewTile returns a new Tile with the input parameters, or nil if the input is invalid.
func NewTile(suit *Suit, ordinal int, id int) (*Tile, error) {
	if suit == nil {
		return nil, fmt.Errorf("Suit cannot be nil")
	}
	if ordinal < 0 || ordinal >= suit.GetSize() {
		return nil, fmt.Errorf("Ordinal out of range [%d, %d): %d", 0, suit.GetSize(), ordinal)
	}
	tile := &Tile{NewTileBase(suit, ordinal), id}
	return tile, nil
}

// GetSuit ...
func (t *Tile) GetSuit() *Suit {
	return t.suit
}

// GetOrdinal ...
func (t *Tile) GetOrdinal() int {
	return t.ordinal
}

// String ...
func (t *Tile) String() string {
	if t.GetSuit().friendlyNameFunc != nil {
		return t.GetSuit().friendlyNameFunc(t)
	}
	return fmt.Sprintf("suit:%s,ord:%d,id:%d", t.GetSuit().GetName(), t.GetOrdinal(), t.id)
}

// CompareTiles is a comparison for tiles. Returns true if tile1 should come before tile2.
func CompareTiles(tile1, tile2 *Tile) bool {
	if tile1.GetSuit().GetSuitType() < tile2.GetSuit().GetSuitType() {
		return true
	}
	if tile1.GetSuit().GetSuitType() > tile2.GetSuit().GetSuitType() {
		return false
	}
	ret := strings.Compare(tile1.GetSuit().GetName(), tile2.GetSuit().GetName())
	if ret != 0 {
		return ret < 0
	}

	return tile1.GetOrdinal() < tile2.GetOrdinal()
}
