package rules

import (
	"fmt"

	"github.com/derekimcheng/mj/domain"
)

var (
	// Dots is one of the three simple types.
	Dots = domain.NewSuit("Dots", domain.SuitTypeSimple, 9, ordinalPlusOneAndSuit("Dots"))
	// Bamboo is one of the three simple types.
	Bamboo = domain.NewSuit("Bamboo", domain.SuitTypeSimple, 9,
		ordinalPlusOneAndSuit("Bamboo"))
	// Characters is one of the three simple types.
	Characters = domain.NewSuit("Characters", domain.SuitTypeSimple, 9,
		ordinalPlusOneAndSuit("Man"))
	// Winds is one the two honor types.
	Winds = domain.NewSuit("Winds", domain.SuitTypeHonor, 4, fixedNameFromOrdinal(windNames))
	// Dragons is one the two honor types.
	Dragons = domain.NewSuit("Dragons", domain.SuitTypeHonor, 3,
		fixedNameFromOrdinal(dragonNames))
	// Flowers is one the two bonus types.
	Flowers = domain.NewSuit("Flowers", domain.SuitTypeBonus, 4, suitAndOrdinalPlusOne("Flower"))
	// Seasons is one the two bonus types.
	Seasons = domain.NewSuit("Seasons", domain.SuitTypeBonus, 4, suitAndOrdinalPlusOne("Season"))
)

func ordinalPlusOneAndSuit(suffix string) domain.TileFriendlyNameFunc {
	return func(t *domain.Tile) string {
		return fmt.Sprintf("%d %s", 1+t.GetOrdinal(), suffix)
	}
}

var windNames = []string{"East", "South", "West", "North"}
var dragonNames = []string{"Red", "Green", "Blue"}

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
	Bamboo, Characters, Dots, // Simples
	Dragons, Winds, // Honors
	Flowers, Seasons, // Bonus
}

// GetSuitsForGame returns the set of all suits used in the game.
// TODO: remove this and replace with one that takes in rule name as input.
func GetSuitsForGame() []*domain.Suit {
	return suits
}

// CanPong returns whether the given suit can be considered for pong, i.e. tiles of same value and
// suit. Note that this also applies to kong.
func CanPong(s *domain.Suit) bool {
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

// IsWindSuit returns whether the suit is the Wind suit. This can be used to distinguish between
// Dragon and Wind suits.
func IsWindSuit(s *domain.Suit) bool {
	return s == Winds
}
