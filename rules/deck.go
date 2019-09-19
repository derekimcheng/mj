package rules

import (
	"github.com/derekimcheng/mj/flags"
	"fmt"
	"github.com/derekimcheng/mj/domain"
	"github.com/golang/glog"
)

// tileCountRule specifies which suits are available and the count of tiles in each suit.
type tileCountRule struct {
	suit  *domain.Suit
	count int
}

// tileCountRules is a list of TileCountRule for a game (HK, ZJ, etc.).
type tileCountRules []tileCountRule

// tileCountRulesHK is the set of rules for Hong Kong MJ.
var tileCountRulesHK = tileCountRules{
	{dots, 4}, {bamboo, 4}, {characters, 4}, {winds, 4}, {dragons, 4}, {flowers, 1}, {seasons, 1},
}

// tileCountRulesZJ is the set of rules for Zung Jung MJ.
var tileCountRulesZJ = tileCountRules{
	{dots, 4}, {bamboo, 4}, {characters, 4}, {winds, 4}, {dragons, 4},
}

// tileCountRulesMap is a map from the string abbreviation of a MJ rule name to its set of rules.
var tileCountRulesMap = map[flags.RuleName]tileCountRules{
	flags.RuleNameHK: tileCountRulesHK,
	flags.RuleNameZJ: tileCountRulesZJ,
}

// NewDeckForGame creates an unshuffled Deck with tiles according for the given rule, or an error if
// the given rule does not exist.
func NewDeckForGame(ruleName flags.RuleName) (domain.Deck, error) {
	rules, found := tileCountRulesMap[ruleName]
	if !found {
		glog.V(2).Infof("Rule %s not found\n", ruleName)
		return nil, fmt.Errorf("Rule %s not found", ruleName)
	}
	var tiles []*domain.Tile
	for _, rule := range rules {
		tiles = addTilesForSuit(rule, tiles)
	}
	return domain.NewDeck(tiles), nil
}

func addTilesForSuit(rule tileCountRule, tiles []*domain.Tile) []*domain.Tile {
	for ordinal := 0; ordinal < rule.suit.GetSize(); ordinal++ {
		for id := 0; id < rule.count; id++ {
			tile, err := domain.NewTile(rule.suit, ordinal, id)
			if err != nil {
				panic(fmt.Errorf("Unable to create tile: suit=%s ordinal=%d id=%d",
					rule.suit.GetName(), ordinal, id))
			}
			tiles = append(tiles, tile)
		}
	}
	return tiles
}
