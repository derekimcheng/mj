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

// HandTileCounter holds information on the count of each tile type + ordinal in a hand.
type HandTileCounter struct {
	// counter maps a suit name to a slice of counters. For each suit, the corresponding
	// slice has length equal to the suit's size.
	counter map[*domain.Suit][]int
}

// NewHandTileCounter creates a new TileCounter with the given hand and the set of all possible
// suits.
func NewHandTileCounter(suits []*domain.Suit, h *domain.Hand) *HandTileCounter {
	counter := make(map[*domain.Suit][]int)
	for _, s := range suits {
		counter[s] = make([]int, s.GetSize())
	}
	for _, t := range h.GetTiles() {
		if !IsEligibleForHand(t.GetSuit()) {
			panic(fmt.Errorf("Hand should not contain ineligible tiles when counting, got %s", t))
		}
		counter[t.GetSuit()][t.GetOrdinal()]++
	}
	return &HandTileCounter{counter: counter}
}

// IsSevenPairs returns true if the hand represents the "Seven Pairs" Out hand. Note that four of
// a kind is considered as two pairs.
func (c *HandTileCounter) IsSevenPairs() bool {
	numPairs := 0
	for _, suit := range c.counter {
		for _, count := range suit {
			if count % 2 == 0 {
				numPairs += count / 2
			}
		}
	}
	return numPairs == 7
}

// IsThirteenOrphans returns true if the hand represents the "Thirteen Orphans" Out hand.
func (c *HandTileCounter) IsThirteenOrphans() bool {
	numTiles := 0
	seenPair := false
	for _, tile := range thirteenOrphanTiles {
		count := c.counter[tile.GetSuit()][tile.GetOrdinal()]
		if count != 1 && count != 2 {
			return false
		}
		if count == 2 {
			if seenPair {
				return false
			}
			seenPair = true
		}
		numTiles += count
	}
	return numTiles == 14
}
