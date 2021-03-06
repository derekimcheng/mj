package domain

// SuitType represents the fundamental type of a Suit. The suit type determines most of the basic
// rules in MJ.
type SuitType int

const (
	// SuitTypeSimple - typically, this has 9 values from 1 to 9 inclusive. This suit type can be
	// considered both for melds and chows. Examples: Bamboo, Dots, Characters.
	SuitTypeSimple SuitType = iota
	// SuitTypeHonor - typically, this has 4 values which are enumerable, but are typically considered for
	// melds only. Examples: Winds, Dragons.
	SuitTypeHonor
	// SuitTypeBonus - typically, this has 4 values which are enumerable, but do not form part of a hand.
	// Examples: Flowers, Seasons.
	SuitTypeBonus
)

// TileFriendlyNameFunc is a type of function that returns a friendly name for a given tile.
type TileFriendlyNameFunc func(*Tile) string

// Suit represents the configuration on a suit of a Tile.
type Suit struct {
	name     string
	suitType SuitType
	// The number of possible values in this suit.
	size             int
	friendlyNameFunc TileFriendlyNameFunc
}

// NewSuit returns a Suit object with the specified parameters.
func NewSuit(name string, suitType SuitType, size int, friendlyNameFunc TileFriendlyNameFunc) *Suit {
	return &Suit{name: name, suitType: suitType, size: size, friendlyNameFunc: friendlyNameFunc}
}

// GetName ...
func (s *Suit) GetName() string {
	return s.name
}

// GetSuitType ...
func (s *Suit) GetSuitType() SuitType {
	return s.suitType
}

// GetSize ...
func (s *Suit) GetSize() int {
	return s.size
}
