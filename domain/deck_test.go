package domain

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func Test_SliceDeckPopFrontBack(t *testing.T) {
	suit := NewSuit("Dots", SuitTypeSimple, 9, nil)
	var tiles []*Tile
	for i := 0; i < suit.GetSize(); i++ {
		tile, _ := NewTile(suit, i, 0)
		require.NotNil(t, tile)
		tiles = append(tiles, tile)
	}
	deck := NewDeck(tiles)
	assert.False(t, deck.IsEmpty())

	tile, err := deck.PopFront()
	require.NotNil(t, tile)
	require.Nil(t, err)
	assert.Equal(t, tiles[0], tile)

	tile, err = deck.PopBack()
	require.NotNil(t, tile)
	require.Nil(t, err)
	assert.Equal(t, tiles[(len(tiles)-1)], tile)
}

func Test_SliceDeckShuffle(t *testing.T) {
	rand.Seed(0)

	suit := NewSuit("Dots", SuitTypeSimple, 9, nil)
	var tiles []*Tile
	for i := 0; i < suit.GetSize(); i++ {
		tile, _ := NewTile(suit, i, 0)
		require.NotNil(t, tile)
		tiles = append(tiles, tile)
	}
	deck := NewDeck(tiles)
	assert.False(t, deck.IsEmpty())

	deck.Shuffle()
	seen := make(map[*Tile]bool)
	for !deck.IsEmpty() {
		tile, err := deck.PopFront()
		require.NotNil(t, tile)
		require.Nil(t, err)

		_, found := seen[tile]
		assert.False(t, found, "Found duplicate tile in shuffled deck")
		seen[tile] = true
	}
	assert.Len(t, seen, len(tiles), "Some tiles were lost during shuffle")
}
