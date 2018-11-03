package domain

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// CreateTileForTest creates a Tile with the given parameters and asserts that the Tile can be
// created.
func CreateTileForTest(t *testing.T, s *Suit, ordinal int) *Tile {
	tile, _ := NewTile(s, ordinal, 0)
	require.NotNil(t, tile, "Failed to create file for suit=%s ord=%d", s.GetName(), ordinal)
	return tile
}
