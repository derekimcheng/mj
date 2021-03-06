package rules

import (
	"github.com/derekimcheng/mj/domain"
	"github.com/derekimcheng/mj/flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_PopulateHands(t *testing.T) {
	deck, err := NewDeckForGame(flags.RuleNameHK)
	require.NotNil(t, deck)
	assert.NoError(t, err)

	var hands []*domain.Hand
	numHands := 4
	for x := 0; x < numHands; x++ {
		hands = append(hands, domain.NewHand())
	}
	PopulateHands(deck, hands)

	for i, h := range hands {
		if i == 0 {
			assert.Equal(t, numTilesPerHand+1, h.NumTiles())
		} else {
			assert.Equal(t, numTilesPerHand, h.NumTiles())
		}
	}
}
