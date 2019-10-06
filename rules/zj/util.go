package zj

import (
	"github.com/derekimcheng/mj/rules"
	"github.com/derekimcheng/mj/domain"
)

// getTilesToCheckForSevenPairs returns for the given Seven Pairs tile group, list of tiles to
// check using the helper functions provided below.
func getTilesToCheckForSevenPairs(group *rules.TileGroup) domain.Tiles {
	// We only need to check every other tile.
	var tiles domain.Tiles
	groupTiles := group.GetTiles()
	for i := 0; i < len(groupTiles); i += 2 {
		tiles = append(tiles, groupTiles[i])
	}
	return tiles
}

// terminalTileHelper returns three bools given a set of tiles:
// The first bool indicates where there is at least one terminal tile.
// The second bool indicates whether there is at least one non-terminal tile.
// The third bool indicates whether there is at least one honor tile.
func terminalTileHelper(tiles domain.Tiles) (bool, bool, bool) {
	hasTerminal := false
	hasNonTerminal := false
	hasHonor := false
	for _, tile := range tiles {
		suitType := tile.GetSuit().GetSuitType()
		if suitType == domain.SuitTypeHonor {
			hasHonor = true
		} else if suitType == domain.SuitTypeSimple {
			if tile.IsTerminal() {
				hasTerminal = true
			} else {
				hasNonTerminal = true
			}
		}
	}
	return hasTerminal, hasNonTerminal, hasHonor
}

type simpleSuitCount int

const (
	noSimpleSuits simpleSuitCount = iota
	oneSimpleSuit
	moreThanOneSimpleSuits
)

// oneSuitHonorHelper returns:
// a simple suit name, if there is exactly one simple suit amongst all tiles. Empty otherwise.
// a bool indicating whether there is at least one honor tile AND there's at most one simple suit.
func oneSuitHonorHelper(tiles domain.Tiles) (simpleSuitCount, bool) {
	hasHonorTiles := false
	var suitName string
	for _, tile := range tiles {
		switch tile.GetSuit().GetSuitType() {
		case domain.SuitTypeSimple:
			if len(suitName) == 0 {
				suitName = tile.GetSuit().GetName()
			} else if suitName != tile.GetSuit().GetName() {
				// Encountered more than one simple suit.
				return moreThanOneSimpleSuits, false
			}
		case domain.SuitTypeHonor:
			hasHonorTiles = true
		}
	}
	if len(suitName) > 0 {
		return oneSimpleSuit, hasHonorTiles
	}
	return noSimpleSuits, hasHonorTiles
}
