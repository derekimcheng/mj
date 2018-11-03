package rules

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/derekimcheng/mj/domain"
)

func Test_SuitsForGame(t *testing.T) {
	var suitNameSeen map[string]bool
	for _, s := range GetSuitsForGame() {
		_, found := suitNameSeen[s.GetName()]
		assert.False(t, found, "Found duplicate suit name %s", s.GetName())
		suitNameSeen[s.GetName()] = true
	}
	assert.Len(t, suitNameSeen, len(GetSuitsForGame()))
}

func Test_CheckCanPong(t *testing.T) {
	pongableSuits := []*domain.Suit{dots, bamboo, characters, winds, dragons}
	nonPongableSuits := []*domain.Suit{flowers, seasons}

	for _, s := range pongableSuits {
		assert.True(t, CanPong(s), "Cannot pong suit %s", s.GetName())
	}

	for _, s := range nonPongableSuits {
		assert.False(t, CanPong(s), "Unnexpectedly able to pong suit %s", s.GetName())
	}
}

func Test_CheckCanChow(t *testing.T) {
	chowableSuits := []*domain.Suit{dots, bamboo, characters}
	nonChowableSuits := []*domain.Suit{winds, dragons, flowers, seasons}

	for _, s := range chowableSuits {
		assert.True(t, CanChow(s), "Cannot chow suit %s", s.GetName())
	}

	for _, s := range nonChowableSuits {
		assert.False(t, CanChow(s), "Unnexpectedly able to chow suit %s", s.GetName())
	}
}

func Test_IsEligibleForHand(t *testing.T) {
	handSuits := []*domain.Suit{dots, bamboo, characters, winds, dragons}
	nonHandSuits := []*domain.Suit{flowers, seasons}

	for _, s := range handSuits {
		assert.True(t, IsEligibleForHand(s), "Suit %s not eligible for hand", s.GetName())
	}

	for _, s := range nonHandSuits {
		assert.False(t, IsEligibleForHand(s), "Suit %s unexpecedly eligible for hand", s.GetName())
	}
}

func Test_TileFriendlyName(t *testing.T) {
	testCases := []struct {
		description   string
		suit          *domain.Suit
		expectedNames []string
	}{
		{
			"Dots",
			dots,
			[]string{"1 Dots", "2 Dots", "3 Dots", "4 Dots", "5 Dots", "6 Dots", "7 Dots",
				"8 Dots", "9 Dots"},
		},
		{
			"Bamboo",
			bamboo,
			[]string{"1 Bamboo", "2 Bamboo", "3 Bamboo", "4 Bamboo", "5 Bamboo", "6 Bamboo",
				"7 Bamboo", "8 Bamboo", "9 Bamboo"},
		},
		{
			"Characters",
			characters,
			[]string{"1 Wan", "2 Wan", "3 Wan", "4 Wan", "5 Wan", "6 Wan", "7 Wan", "8 Wan",
				"9 Wan"},
		},
		{
			"Winds",
			winds,
			[]string{"East", "South", "West", "North"},
		},
		{
			"Dragons",
			dragons,
			[]string{"White", "Red", "Blue"},
		},
		{
			"Flowers",
			flowers,
			[]string{"Flower 1", "Flower 2", "Flower 3", "Flower 4"},
		},
		{
			"Seasons",
			seasons,
			[]string{"Season 1", "Season 2", "Season 3", "Season 4"},
		},
	}
	for _, tc := range testCases {
		require.Len(t, tc.expectedNames, tc.suit.GetSize(), tc.description)
		for i := 0; i < tc.suit.GetSize(); i++ {
			tile, _ := domain.NewTile(tc.suit, i, 0)
			require.NotNil(t, tile)
			assert.Equal(t, tc.expectedNames[i], tile.String(),
				fmt.Sprintf("%s tile %d", tc.description, i))
		}
	}
}
