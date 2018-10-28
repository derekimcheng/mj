package rules

import "github.com/derekimcheng/mj/domain"

var (
	dots       = domain.NewSuit("Dots", domain.SuitTypeSimple, 10)
	bamboo     = domain.NewSuit("Bamboo", domain.SuitTypeSimple, 10)
	characters = domain.NewSuit("Characters", domain.SuitTypeSimple, 10)
	winds      = domain.NewSuit("Winds", domain.SuitTypeHonor, 4)
	dragons    = domain.NewSuit("Dragons", domain.SuitTypeHonor, 4)
	flowers    = domain.NewSuit("Flowers", domain.SuitTypeBonus, 4)
	seasons    = domain.NewSuit("Seasons", domain.SuitTypeBonus, 4)
)
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
