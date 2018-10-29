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