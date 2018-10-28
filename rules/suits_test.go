package rules

import (
	"testing"

	"github.com/derekimcheng/mj/domain"
)

func Test_SanityCheckCanMeld(t *testing.T) {
	meldableSuits := []*domain.Suit{dots, bamboo, characters, winds, dragons}
	nonMeldableSuits := []*domain.Suit{flowers, seasons}

	for _, s := range meldableSuits {
		if !CanMeld(s) {
			t.Errorf("Cannot meld suit %s", s.GetName())
		}
	}

	for _, s := range nonMeldableSuits {
		if CanMeld(s) {
			t.Errorf("Unnexpectedly able to meld suit %s", s.GetName())
		}
	}
}

func Test_SanityCheckCanChow(t *testing.T) {
	chowableSuits := []*domain.Suit{dots, bamboo, characters}
	nonChowableSuits := []*domain.Suit{winds, dragons, flowers, seasons}

	for _, s := range chowableSuits {
		if !CanChow(s) {
			t.Errorf("Cannot chow suit %s", s.GetName())
		}
	}

	for _, s := range nonChowableSuits {
		if CanChow(s) {
			t.Errorf("Unnexpectedly able to chow suit %s", s.GetName())
		}
	}
}

func Test_SanityIsEligibleForHand(t *testing.T) {
	handSuits := []*domain.Suit{dots, bamboo, characters, winds, dragons}
	nonHandSuits := []*domain.Suit{flowers, seasons}

	for _, s := range handSuits {
		if !IsEligibleForHand(s) {
			t.Errorf("Suit %s not eligible for hand", s.GetName())
		}
	}

	for _, s := range nonHandSuits {
		if IsEligibleForHand(s) {
			t.Errorf("Suit %s unexpecedly eligible for hand", s.GetName())
		}
	}
}
