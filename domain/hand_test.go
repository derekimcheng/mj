package domain

import (
	"sort"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func createTile(t *testing.T, s *Suit, ordinal int) *Tile {
	tile, _ := NewTile(s, ordinal, 0)
	require.NotNil(t, tile, "Failed to create file for suit=%s ord=%d", s.GetName(), ordinal)
	return tile
}
func Test_AddRemoveTile(t *testing.T) {
	dots := NewSuit("Dots", SuitTypeSimple, 9, nil)
	bamboo := NewSuit("Bamboo", SuitTypeSimple, 9, nil)
	characters := NewSuit("Characters", SuitTypeSimple, 9, nil)
	// winds   := NewSuit("Winds", SuitTypeHonor, 4, nil)
	// dragons := NewSuit("Dragons", SuitTypeHonor, 3, nil)
	// flowers := NewSuit("Flowers", SuitTypeBonus, 4, nil)
	// seasons := NewSuit("Seasons", SuitTypeBonus, 4, nil)

	hand := NewHand()
	expectedTiles := []*Tile{
		createTile(t, dots, 0),
		createTile(t, dots, 8),
		createTile(t, bamboo, 0),
		createTile(t, bamboo, 8),
		createTile(t, characters, 0),
		createTile(t, characters, 8),
	}
	for _, tile := range expectedTiles {
		hand.AddTile(tile)
	}
	assert.Equal(t, expectedTiles, hand.GetTiles())

	tile, err := hand.RemoveTile(20)
	assert.Nil(t, tile)
	assert.NotNil(t, err)

	expectedTiles = expectedTiles[1:]
	tile, err = hand.RemoveTile(0)
	assert.NotNil(t, tile)
	assert.Nil(t, err)
	assert.Equal(t, expectedTiles, hand.GetTiles())
}

func Test_HandSort(t *testing.T) {
	dots := NewSuit("Dots", SuitTypeSimple, 9, nil)
	bamboo := NewSuit("Bamboo", SuitTypeSimple, 9, nil)
	characters := NewSuit("Characters", SuitTypeSimple, 9, nil)
	winds := NewSuit("Winds", SuitTypeHonor, 4, nil)
	dragons := NewSuit("Dragons", SuitTypeHonor, 3, nil)

	dots1 := createTile(t, dots, 0)
	dots9 := createTile(t, dots, 8)
	bamboo1 := createTile(t, bamboo, 0)
	bamboo9 := createTile(t, bamboo, 8)
	char1 := createTile(t, characters, 0)
	char9 := createTile(t, characters, 8)
	wind1 := createTile(t, winds, 0)
	wind2 := createTile(t, winds, 1)
	wind3 := createTile(t, winds, 2)
	wind4 := createTile(t, winds, 3)
	dragon1 := createTile(t, dragons, 0)
	dragon2 := createTile(t, dragons, 1)
	dragon3 := createTile(t, dragons, 2)

	unsortedTiles := []*Tile{
		char1, dragon2, wind4, wind3, dots1, bamboo9, dots9, bamboo1, wind2, char9, wind1, dragon3,
		dragon1,
	}
	sortedTiles := []*Tile{
		bamboo1, bamboo9, char1, char9, dots1, dots9, dragon1, dragon2, dragon3, wind1, wind2,
		wind3, wind4,
	}
	hand := NewHand()
	for _, tile := range unsortedTiles {
		hand.AddTile(tile)
	}

	sort.Sort(hand)
	assert.Equal(t, sortedTiles, hand.GetTiles())
}
