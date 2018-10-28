package rules

import (
	"fmt"

	"github.com/derekimcheng/mj/domain"
)

var (
	dots       = domain.NewSuit("Dots", domain.SuitTypeSimple, 9, ordinalPlusOneAndSuit("Dots"))
	bamboo     = domain.NewSuit("Bamboo", domain.SuitTypeSimple, 9, ordinalPlusOneAndSuit("Bamboo"))
	characters = domain.NewSuit("Characters", domain.SuitTypeSimple, 9, ordinalPlusOneAndSuit("Wan"))
	winds      = domain.NewSuit("Winds", domain.SuitTypeHonor, 4, fixedNameFromOrdinal(windNames))
	dragons    = domain.NewSuit("Dragons", domain.SuitTypeHonor, 3, fixedNameFromOrdinal(dragonNames))
	flowers    = domain.NewSuit("Flowers", domain.SuitTypeBonus, 4, suitAndOrdinalPlusOne("Flower"))
	seasons    = domain.NewSuit("Seasons", domain.SuitTypeBonus, 4, suitAndOrdinalPlusOne("Season"))
)

func ordinalPlusOneAndSuit(suffix string) domain.TileFriendlyNameFunc {
	return func(t *domain.Tile) string {
		return fmt.Sprintf("%d %s", 1+t.GetOrdinal(), suffix)
	}
}

var windNames = []string{"East", "South", "West", "North"}
var dragonNames = []string{"White", "Red", "Blue"}

func fixedNameFromOrdinal(fixedNames []string) domain.TileFriendlyNameFunc {
	return func(t *domain.Tile) string {
		return fixedNames[t.GetOrdinal()]
	}
}

func suitAndOrdinalPlusOne(prefix string) domain.TileFriendlyNameFunc {
	return func(t *domain.Tile) string {
		return fmt.Sprintf("%s %d", prefix, 1+t.GetOrdinal())
	}
}

var suits = []*domain.Suit{
	dots, bamboo, characters, winds, dragons, flowers, seasons,
}

// GetSuits returns the set of all suits used in the game.
func GetSuits() []*domain.Suit {
	return suits
}

// CanMeld returns whether the given suit can be considered for melds, i.e. tiles of same value and
// suit.
func CanMeld(s *domain.Suit) bool {
	return s.GetSuitType() == domain.SuitTypeSimple || s.GetSuitType() == domain.SuitTypeHonor
}

// CanChow returns whether the given suit can be considered for chows, i.e. tiles of consecutive
// values of the same suit.
func CanChow(s *domain.Suit) bool {
	return s.GetSuitType() == domain.SuitTypeSimple
}

// IsEligibleForHand returns whether a tile of a given suit can make up part of a hand. If not, the
// tile must be discarded and replaced with another tile when it is drawn.
func IsEligibleForHand(s *domain.Suit) bool {
	return s.GetSuitType() != domain.SuitTypeBonus
}
