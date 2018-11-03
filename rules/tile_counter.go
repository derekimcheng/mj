package rules

import (
	"fmt"
	"github.com/derekimcheng/mj/domain"
)

var (
	thirteenOrphanTiles = []domain.TileBase{
		domain.NewTileBase(bamboo, 0),
		domain.NewTileBase(bamboo, 8),
		domain.NewTileBase(dots, 0),
		domain.NewTileBase(dots, 8),
		domain.NewTileBase(characters, 0),
		domain.NewTileBase(characters, 8),
		domain.NewTileBase(winds, 0),
		domain.NewTileBase(winds, 1),
		domain.NewTileBase(winds, 2),
		domain.NewTileBase(winds, 3),
		domain.NewTileBase(dragons, 0),
		domain.NewTileBase(dragons, 1),
		domain.NewTileBase(dragons, 2),
	}
)

type tileInventory = map[*domain.Suit][][]*domain.Tile

// HandTileCounter holds information on the count of each tile type + ordinal in a hand.
type HandTileCounter struct {
	// inventory maps a suit name and ordinal to a slice of tiles.
	inventory        tileInventory
	computedOutPlans *OutPlans
}

// NewHandTileCounter creates a new TileCounter with the given hand and the set of all possible
// suits.
func NewHandTileCounter(suits []*domain.Suit, h *domain.Hand) *HandTileCounter {
	inventory := make(tileInventory)
	for _, s := range suits {
		inventory[s] = make([][]*domain.Tile, s.GetSize())
	}
	for _, t := range h.GetTiles() {
		if !IsEligibleForHand(t.GetSuit()) {
			panic(fmt.Errorf("Hand should not contain ineligible tiles when counting, got %s", t))
		}
		tiles := inventory[t.GetSuit()][t.GetOrdinal()]
		inventory[t.GetSuit()][t.GetOrdinal()] = append(tiles, t)
	}
	return &HandTileCounter{inventory: inventory}
}

// ComputeOutPlans evalautes the current hand and generate possible Out plans for it.
func (c *HandTileCounter) ComputeOutPlans() OutPlans {
	if c.computedOutPlans != nil {
		return *c.computedOutPlans
	}

	c.computedOutPlans = &OutPlans{}
	numTiles := 0
	for _, suit := range c.inventory {
		for _, tiles := range suit {
			numTiles += len(tiles)
		}
	}
	computeOutPlans(numTiles, &c.inventory, c.computedOutPlans)
	return *c.computedOutPlans
}
