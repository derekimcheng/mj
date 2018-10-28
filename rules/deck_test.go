package rules

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_NewDeckForGame(t *testing.T) {
	deck := NewDeckForGame()

	// Build a map from friendly name to count
	friendlyNameCounts := make(map[string]int)
	numTiles := 0
	for !deck.IsEmpty() {
		tile, _ := deck.PopFront()
		require.NotNil(t, tile)
		friendlyNameCounts[tile.String()]++
		numTiles++
	}

	expectedNumTiles := 144
	expectedNumTileValues := 42

	assert.Equal(t, numTiles, expectedNumTiles)
	assert.Len(t, friendlyNameCounts, expectedNumTileValues)
}
