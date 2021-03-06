package rules

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/derekimcheng/mj/domain"
)

func Test_SuitsForGame(t *testing.T) {
	suitNameSeen := make(map[string]bool)
	for _, s := range GetSuitsForGame() {
		_, found := suitNameSeen[s.GetName()]
		assert.False(t, found, "Found duplicate suit name %s", s.GetName())
		suitNameSeen[s.GetName()] = true
	}
	assert.Len(t, suitNameSeen, len(GetSuitsForGame()))
}

func Test_CheckCanPong(t *testing.T) {
	pongableSuits := []*domain.Suit{Dots, Bamboo, Characters, Winds, Dragons}
	nonPongableSuits := []*domain.Suit{Flowers, Seasons}

	for _, s := range pongableSuits {
		assert.True(t, CanPong(s), "Cannot pong suit %s", s.GetName())
	}

	for _, s := range nonPongableSuits {
		assert.False(t, CanPong(s), "Unnexpectedly able to pong suit %s", s.GetName())
	}
}

func Test_CheckCanChow(t *testing.T) {
	chowableSuits := []*domain.Suit{Dots, Bamboo, Characters}
	nonChowableSuits := []*domain.Suit{Winds, Dragons, Flowers, Seasons}

	for _, s := range chowableSuits {
		assert.True(t, CanChow(s), "Cannot chow suit %s", s.GetName())
	}

	for _, s := range nonChowableSuits {
		assert.False(t, CanChow(s), "Unnexpectedly able to chow suit %s", s.GetName())
	}
}

func Test_IsEligibleForHand(t *testing.T) {
	handSuits := []*domain.Suit{Dots, Bamboo, Characters, Winds, Dragons}
	nonHandSuits := []*domain.Suit{Flowers, Seasons}

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
			Dots,
			[]string{"1 Dots", "2 Dots", "3 Dots", "4 Dots", "5 Dots", "6 Dots", "7 Dots",
				"8 Dots", "9 Dots"},
		},
		{
			"Bamboo",
			Bamboo,
			[]string{"1 Bamboo", "2 Bamboo", "3 Bamboo", "4 Bamboo", "5 Bamboo", "6 Bamboo",
				"7 Bamboo", "8 Bamboo", "9 Bamboo"},
		},
		{
			"Characters",
			Characters,
			[]string{"1 Man", "2 Man", "3 Man", "4 Man", "5 Man", "6 Man", "7 Man", "8 Man",
				"9 Man"},
		},
		{
			"Winds",
			Winds,
			[]string{"East", "South", "West", "North"},
		},
		{
			"Dragons",
			Dragons,
			[]string{"Red", "Green", "Blue"},
		},
		{
			"Flowers",
			Flowers,
			[]string{"Flower 1", "Flower 2", "Flower 3", "Flower 4"},
		},
		{
			"Seasons",
			Seasons,
			[]string{"Season 1", "Season 2", "Season 3", "Season 4"},
		},
	}
	for _, tc := range testCases {
		require.Len(t, tc.expectedNames, tc.suit.GetSize(), tc.description)
		for i := 0; i < tc.suit.GetSize(); i++ {
			tile, _ := domain.NewTile(tc.suit, i, 0)
			require.NotNil(t, tile)
			assert.Equal(t, "["+tc.expectedNames[i]+"]", tile.String(),
				fmt.Sprintf("%s tile %d", tc.description, i))
		}
	}
}
