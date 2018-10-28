package rules

import "github.com/derekimcheng/mj/domain"

var suits = []*domain.Suit{
	domain.NewSuit("Dots", domain.SuitTypeSimple, 10),
	domain.NewSuit("Bamboo", domain.SuitTypeSimple, 10),
	domain.NewSuit("Characters", domain.SuitTypeSimple, 10),
	domain.NewSuit("Wind", domain.SuitTypeHonor, 4),
	domain.NewSuit("Dragons", domain.SuitTypeHonor, 4),
	domain.NewSuit("Flowers", domain.SuitTypeBonus, 4),
	domain.NewSuit("Seasons", domain.SuitTypeBonus, 4),
}

// GetSuits returns the set of all suits used in the game.
func GetSuits() []*domain.Suit {
	return suits
}

// CanMeld returns whether the given suit can be considered for melds.
func CanMeld(s *domain.Suit) bool {
	return s.GetSuitType() == domain.SuitTypeSimple
}

// CanChow returns whether the given suit can be considered for chows.
func CanChow(s *domain.Suit) bool {
	return s.GetSuitType() == domain.SuitTypeSimple || s.GetSuitType() == domain.SuitTypeHonor
}

// IsEligibleForHand returns whether a tile of a given suit can make up part of a hand. If not, the
// tile must be discarded and replaced with another tile when it is drawn.
func IsEligibleForHand(s *domain.Suit) bool {
	return s.GetSuitType() != domain.SuitTypeBonus
}
