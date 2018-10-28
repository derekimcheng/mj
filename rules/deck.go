package rules

import (
	"fmt"
	"github.com/derekimcheng/mj/domain"
)

type tileCountRule struct {
	suit  *domain.Suit
	count int
}

var tileCountRules = []tileCountRule{
	{dots, 4}, {bamboo, 4}, {characters, 4}, {winds, 4}, {dragons, 4}, {flowers, 1}, {seasons, 1},
}

// NewDeckForGame creates an unshuffled Deck with tiles according to the rules.
func NewDeckForGame() domain.Deck {
	var tiles []*domain.Tile

	for _, rule := range tileCountRules {
		tiles = addTilesForSuit(rule, tiles)
	}

	return domain.NewDeck(tiles)
}

func addTilesForSuit(rule tileCountRule, tiles []*domain.Tile) []*domain.Tile {
	for ordinal := 0; ordinal < rule.suit.GetSize(); ordinal++ {
		for id := 0; id < rule.count; id++ {
			tile, err := domain.NewTile(rule.suit, ordinal, id)
			fmt.Printf("Adding tile %s (id=%d) to deck\n", tile, id)
			if err != nil {
				panic(fmt.Errorf("Unable to create tile: suit=%s ordinal=%d id=%d",
					rule.suit.GetName(), ordinal, id))
			}
			tiles = append(tiles, tile)
		}
	}
	return tiles
}
