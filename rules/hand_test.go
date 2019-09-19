package rules

import (
	"github.com/stretchr/testify/require"
	"github.com/derekimcheng/mj/flags"
	"github.com/stretchr/testify/assert"
	"github.com/derekimcheng/mj/domain"
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

	for _, h := range hands {
		assert.Equal(t, numTilesPerHand, h.NumTiles())
	}
}