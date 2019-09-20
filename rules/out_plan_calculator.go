package rules

import (
	"github.com/golang/glog"
	"sort"
	"fmt"
	"github.com/derekimcheng/mj/domain"
)

type specialPlanMatcherFunc func(numRemainingTiles int, inventory *tileInventory) OutTileGroups

var (
	specialPlanMatchers = []specialPlanMatcherFunc{
		matchSevenPairs,
		matchThirteenOrphans,
	}
)

// OutPlanCalculator calculators Out plans for a given state.
type OutPlanCalculator struct {
	handInventory    tileInventory
	meldedGroups     OutTileGroups
	computedOutPlans *OutPlans
}

// NewOutPlanCalculator creates a new OutPlanCalculator with the given state.
// TODO: discardTile should be generalized into a source: self-drawn, kong, discard, replacement.
func NewOutPlanCalculator(suits []*domain.Suit, player *PlayerGameState, discardTile *domain.Tile) *OutPlanCalculator {
	inventory := make(tileInventory)
	for _, s := range suits {
		inventory[s] = make([][]*domain.Tile, s.GetSize())
	}
	allTiles := player.GetHand().GetTiles()
	if discardTile != nil {
		allTiles = append(allTiles, discardTile)
	}
	for _, t := range allTiles {
		if !IsEligibleForHand(t.GetSuit()) {
			panic(fmt.Errorf("Hand should not contain ineligible tiles when counting, got %s", t))
		}
		tiles := inventory[t.GetSuit()][t.GetOrdinal()]
		inventory[t.GetSuit()][t.GetOrdinal()] = append(tiles, t)
	}

	meldedGroupsCopy := append(player.GetMeldGroups())
	sort.Sort(meldedGroupsCopy)

	return &OutPlanCalculator{handInventory: inventory, meldedGroups: meldedGroupsCopy, computedOutPlans: nil}
}

// Calculate generates possible Out plans for the given hand / melded groups.
func (c *OutPlanCalculator) Calculate() OutPlans {
	if c.computedOutPlans != nil {
		return *c.computedOutPlans
	}

	c.computedOutPlans = &OutPlans{}
	numTiles := 0
	for _, suit := range c.handInventory {
		for _, tiles := range suit {
			numTiles += len(tiles)
		}
	}
	c.computeOutPlans(numTiles, &c.handInventory, c.computedOutPlans)
	return *c.computedOutPlans
}

// computeOutPlans is the entry point for generating out plans.
func (c *OutPlanCalculator) computeOutPlans(
	numRemainingTiles int,
	inventory *tileInventory,
	outPlansSoFar *OutPlans) {
	c.computeSpecialPlans(numRemainingTiles, inventory, outPlansSoFar)
	c.computeOutPlansHelper(numRemainingTiles, inventory, nil, nil, outPlansSoFar)
}

func (c *OutPlanCalculator) computeSpecialPlans(
	numRemainingTiles int,
	inventory *tileInventory,
	outPlansSoFar *OutPlans) {
	for _, matcher := range specialPlanMatchers {
		groups := matcher(numRemainingTiles, inventory)
		if groups != nil {
			*outPlansSoFar = append(*outPlansSoFar, c.generateNewOutPlan(groups))
		}
	}
}

// IsSevenPairs returns an OutPlan if the hand represents the "Seven Pairs" Out hand. Note that
// four of a kind is considered as two pairs.
func matchSevenPairs(numRemainingTiles int, inventory *tileInventory) OutTileGroups {
	if numRemainingTiles != 14 {
		return nil
	}

	numPairs := 0
	var outTiles domain.Tiles
	for _, suit := range *inventory {
		for _, tiles := range suit {
			if len(tiles)%2 == 0 {
				numPairs += len(tiles) / 2
				outTiles = append(outTiles, tiles...)
			}
		}
	}
	if numPairs != 7 {
		return nil
	}

	return OutTileGroups{NewOutTileGroup(outTiles, OutTileGroupTypeSevenPairs)}
}

// matchThirteenOrphans returns an OutPlan if the hand represents the "Thirteen Orphans" Out hand.
func matchThirteenOrphans(numRemainingTiles int, inventory *tileInventory) OutTileGroups {
	if numRemainingTiles != 14 {
		return nil
	}

	seenPair := false
	var outTiles domain.Tiles
	for _, tile := range thirteenOrphanTiles {
		tiles := (*inventory)[tile.GetSuit()][tile.GetOrdinal()]
		numTiles := len(tiles)
		if numTiles != 1 && numTiles != 2 {
			return nil
		}
		if numTiles == 2 {
			if seenPair {
				return nil
			}
			seenPair = true
		}
		outTiles = append(outTiles, tiles...)
	}

	return OutTileGroups{NewOutTileGroup(outTiles, OutTileGroupTypeThirteenOrphans)}
}

func (c *OutPlanCalculator) computeOutPlansHelper(
	numRemainingTiles int,
	inventory *tileInventory,
	outGroupsSoFar OutTileGroups,
	pairTileGroup *OutTileGroup,
	outPlansSoFar *OutPlans) {
	if numRemainingTiles == 0 {
		if pairTileGroup == nil {
			glog.V(2).Infof("Base case: not a valid plan because there is no pair group")
			return
		}
		outGroupsSoFar := append(outGroupsSoFar, pairTileGroup)
		outPlan := c.generateNewOutPlan(outGroupsSoFar)
		*outPlansSoFar = append(*outPlansSoFar, outPlan)
		glog.V(2).Infof("Base case: added valid plan to set\n")
		return
	}

	for _, suit := range *inventory {
		for i, tiles := range suit {
			if len(tiles) == 0 {
				continue
			}

			used := false
			// Use this tile as a chow or pong group (or pair group, if there is not one assigned
			// yet). If it cannot be used in any group, this path cannot possibly yield a valid
			// plan.
			if pairTileGroup == nil && len(tiles) >= 2 {
				used = true
				oldTiles := suit[i]
				outTiles := tiles[:2]
				glog.V(2).Infof("Using %s as Pair out group\n", outTiles)
				newPairOutGroup := NewOutTileGroup(outTiles, OutTileGroupTypePair)
				// Take out 2 tiles from inventory for pair group.
				suit[i] = suit[i][2:]
				c.computeOutPlansHelper(numRemainingTiles-2, inventory, outGroupsSoFar,
					newPairOutGroup, outPlansSoFar)
				// Restore previous state.
				suit[i] = oldTiles
			}
			// Note that it is possible for a kind of tile to be used as a pong (or a pair) and
			// chow at the same time.
			if len(tiles) >= 3 {
				used = true
				oldTiles := suit[i]
				outTiles := tiles[:3]
				glog.V(2).Infof("Using %s as Pong out group\n", outTiles)
				newPongOutGroup := NewOutTileGroup(outTiles, OutTileGroupTypePong)
				outGroupsSoFar = append(outGroupsSoFar, newPongOutGroup)
				// Take out 3 tiles from inventory for pong group.
				suit[i] = suit[i][3:]
				c.computeOutPlansHelper(numRemainingTiles-3, inventory, outGroupsSoFar,
					pairTileGroup, outPlansSoFar)
				// Restore previous state.
				suit[i] = oldTiles
				outGroupsSoFar = outGroupsSoFar[:len(outGroupsSoFar)-1]
			}
			if i < len(suit)-2 {
				numChows := 0
				oldTiles1 := suit[i]
				oldTiles2 := suit[i+1]
				oldTiles3 := suit[i+2]
				for len(suit[i]) >= 1 && len(suit[i+1]) >= 1 && len(suit[i+2]) >= 1 {
					numChows++
					outTiles := domain.Tiles{suit[i][0], suit[i+1][0], suit[i+2][0]}
					glog.V(2).Infof("Using %s as Chow out group (%d)\n", outTiles, numChows)
					newChowOutGroup := NewOutTileGroup(outTiles, OutTileGroupTypeChow)
					outGroupsSoFar = append(outGroupsSoFar, newChowOutGroup)
					// Take out 3 tiles from inventory for pong group.
					suit[i] = suit[i][1:]
					suit[i+1] = suit[i+1][1:]
					suit[i+2] = suit[i+2][1:]
				}
				if numChows > 0 {
					// There is only use recursing if we have used up all tiles in the current
					// suit+ordinal. Otherwise, it is a dead tile. Note that the pair -> chow
					// recursion path is taken care of above. Because of this check, there won't be
					// a chow -> pair path.
					if len(suit[i]) == 0 {
						used = true
						c.computeOutPlansHelper(numRemainingTiles-3*numChows, inventory,
							outGroupsSoFar, pairTileGroup, outPlansSoFar)
					}
					// Restore previous state.
					suit[i] = oldTiles1
					suit[i+1] = oldTiles2
					suit[i+2] = oldTiles3
					outGroupsSoFar = outGroupsSoFar[:len(outGroupsSoFar)-numChows]
				}
			}

			if !used {
				glog.V(2).Infof("%s is a dead tile -- no further match can be made\n", tiles[0])
			}

			// To prevent duplicates, we will let the next level of recursion handle tiles in the
			// next ordinal / suit. This is regardless whether the current suit + ordinal has been
			// used in a out group.
			return
		}
	}
}

func (c *OutPlanCalculator) generateNewOutPlan(handGroups OutTileGroups) OutPlan {
	return NewOutPlan(handGroups, c.meldedGroups)
}