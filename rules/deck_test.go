package rules

import (
	"github.com/derekimcheng/mj/flags"
	"github.com/derekimcheng/mj/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_NewDeckForGame(t *testing.T) {
	deck, err := NewDeckForGame(flags.RuleNameHK)
	require.NotNil(t, deck)
	assert.NoError(t, err)

	// Build a map from friendly name to count
	nonBonusCounts := make(map[string]int)
	bonusCounts := make(map[string]int)
	numTiles := 0
	for !deck.IsEmpty() {
		tile, _ := deck.PopFront()
		require.NotNil(t, tile)
		if tile.GetSuit().GetSuitType() == domain.SuitTypeBonus {
			bonusCounts[tile.String()]++
		} else {
			nonBonusCounts[tile.String()]++
		}
		numTiles++
	}

	expectedNumTiles := 144
	expectedNumTileValues := 42
	expectedNumTilesPerValueNonBonus := 4
	expectedNumTilesPerValueBonus := 1
	assert.Equal(t, numTiles, expectedNumTiles)
	assert.Equal(t, expectedNumTileValues, len(nonBonusCounts)+len(bonusCounts))
	for friendlyName, count := range nonBonusCounts {
		assert.Equal(t, expectedNumTilesPerValueNonBonus, count,
			"Unexpected number of tiles for %s", friendlyName)
	}
	for friendlyName, count := range bonusCounts {
		assert.Equal(t, expectedNumTilesPerValueBonus, count,
			"Unexpected number of tiles for %s", friendlyName)
	}
}

func Test_NewDeckForGame_UnknownRule(t *testing.T) {
	deck, err := NewDeckForGame("unknownrule")
	assert.Nil(t, deck)
	assert.Error(t, err)
}
