package rules

import (
	"fmt"
	"github.com/derekimcheng/mj/domain"
)

const (
	numTilesPerHand        = 13
	numTilesPerInitialDraw = 4
)

// PopulateHands populates the given hands with the given deck. Returns nil on success, or error
// if an error occurred.
func PopulateHands(d domain.Deck, hands []*domain.Hand) error {
	numRequiredTiles := numTilesPerHand * numTilesPerInitialDraw
	if d.NumRemainingTiles() < numRequiredTiles {
		return fmt.Errorf("Not enough remaining tiles in deck, have %d, need %d",
			d.NumRemainingTiles(), numRequiredTiles)
	}

	// Full draws, each player takes turn drawing numTilesPerInitialDraw tiles.
	numFullRounds := numTilesPerHand / numTilesPerInitialDraw
	for round := 0; round < numFullRounds; round++ {
		for _, h := range hands {
			drawTilesFromDeck(d, h, numTilesPerInitialDraw)
		}
	}

	// Partial draws, each player takes turn drawing 1 tile at a time.
	numPartialRounds := numTilesPerHand % numTilesPerInitialDraw
	for round := 0; round < numPartialRounds; round++ {
		for _, h := range hands {
			drawTilesFromDeck(d, h, 1)
		}
	}

	return nil
}

func drawTilesFromDeck(d domain.Deck, h *domain.Hand, numTiles int) {
	for x := 0; x < numTiles; x++ {
		tile, err := d.PopFront()
		if err != nil {
			panic(err)
		}
		h.AddTile(tile)
	}
}

// RemoveIneligibleTilesFromHand removes ineligible (e.g. bonus) tiles from the given hand and
// returns the removed tiles.
func RemoveIneligibleTilesFromHand(h *domain.Hand) []*domain.Tile {
	var eligibleTiles, ineligibleTiles []*domain.Tile
	for _, t := range h.GetTiles() {
		if IsEligibleForHand(t.GetSuit()) {
			eligibleTiles = append(eligibleTiles, t)
		} else {
			ineligibleTiles = append(ineligibleTiles, t)
		}
	}
	h.SetTiles(eligibleTiles)
	return ineligibleTiles
}

// IsOutHand returns whether the given hand represents an Out hand.
// An Out hand consists of a pair + triplets of pongs/kongs and melds, with a couple of exceptions:
// (1) Seven pairs
// (2) 13 Orphans
func IsOutHand(h *domain.Hand) bool {
	counter := NewHandTileCounter(GetSuitsForGame(), h)
	// TODO: expand other rules.
	return counter.IsSevenPairs() || counter.IsThirteenOrphans()
}