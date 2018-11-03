package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)


func Test_AddRemoveTile(t *testing.T) {
	dots := NewSuit("Dots", SuitTypeSimple, 9, nil)
	bamboo := NewSuit("Bamboo", SuitTypeSimple, 9, nil)
	characters := NewSuit("Characters", SuitTypeSimple, 9, nil)

	hand := NewHand()
	expectedTiles := []*Tile{
		CreateTileForTest(t, dots, 0),
		CreateTileForTest(t, dots, 8),
		CreateTileForTest(t, bamboo, 0),
		CreateTileForTest(t, bamboo, 8),
		CreateTileForTest(t, characters, 0),
		CreateTileForTest(t, characters, 8),
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

	dots1 := CreateTileForTest(t, dots, 0)
	dots9 := CreateTileForTest(t, dots, 8)
	bamboo1 := CreateTileForTest(t, bamboo, 0)
	bamboo9 := CreateTileForTest(t, bamboo, 8)
	char1 := CreateTileForTest(t, characters, 0)
	char9 := CreateTileForTest(t, characters, 8)
	wind1 := CreateTileForTest(t, winds, 0)
	wind2 := CreateTileForTest(t, winds, 1)
	wind3 := CreateTileForTest(t, winds, 2)
	wind4 := CreateTileForTest(t, winds, 3)
	dragon1 := CreateTileForTest(t, dragons, 0)
	dragon2 := CreateTileForTest(t, dragons, 1)
	dragon3 := CreateTileForTest(t, dragons, 2)

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

	hand.Sort()
	assert.Equal(t, sortedTiles, hand.GetTiles())
}
