package rules

import (
	"fmt"
	"github.com/derekimcheng/mj/domain"
	"github.com/golang/glog"
	"sort"
	"strings"
)

// OutPlans is a slice of OutPlan.
// TODO: OutPlans should be sorted according to score / type /etc.
type OutPlans []OutPlan

// String ...
func (plans OutPlans) String() string {
	if len(plans) == 0 {
		return ""
	}
	ret := ""
	for i, plan := range plans {
		if i > 0 {
			ret += ", "
		}
		ret += fmt.Sprintf("Plan %d: %s", i+1, plan)
	}
	return ret
}

// OutPlan consists of a strategy of declaring an Out using a combination of tile groups in the
// hand and meld area.
type OutPlan struct {
	handGroups   OutTileGroups
	meldedGroups OutTileGroups
}

// GetHandGroups ...
func (p OutPlan) GetHandGroups() OutTileGroups {
	return p.handGroups
}

// String ...
func (p OutPlan) String() string {
	handGroupStrs := []string{}
	for _, handGroup := range p.handGroups {
		handGroupStrs = append(handGroupStrs, handGroup.String())
	}
	hand := fmt.Sprintf("Hand: %s", strings.Join(handGroupStrs, ", "))

	meldedGroupStrs := []string{}
	for _, meldedGroup := range p.meldedGroups {
		meldedGroupStrs = append(meldedGroupStrs, meldedGroup.String())
	}
	melded := fmt.Sprintf("Melded: %s", strings.Join(meldedGroupStrs, ", "))

	return hand + " " + melded
}

// NewOutPlan creates a new OutPlan with the given parameters. The input groups is not
// copied, and will be modified by sorting.
func NewOutPlan(handGroups OutTileGroups, meldedGroups OutTileGroups) OutPlan {
	sort.Sort(handGroups)
	return OutPlan{handGroups: handGroups, meldedGroups: meldedGroups}
}

// OutTileGroupType represents a type of Out tile group
type OutTileGroupType int

const (
	// OutTileGroupTypePair represents a pair of same suit+ordinal tiles. With the exception of
	// special hands, each out plan must contain exactly one Pair group.
	OutTileGroupTypePair OutTileGroupType = iota
	// OutTileGroupTypePong represents a triplet of the same suit+ordinal tiles.
	OutTileGroupTypePong
	// OutTileGroupTypeChow represents a triplet of consecutive tiles.
	OutTileGroupTypeChow
	// OutTileGroupTypeKong represents a quadruplet of the same suit+ordinal tiles. This is only
	// available in the meld area.
	OutTileGroupTypeKong
	// OutTileGroupTypeConcealedKong represents a concealed quadruplet of the same suit+ordinal
	// tiles. This is only available in the meld area.
	OutTileGroupTypeConcealedKong
	// OutTileGroupTypeSevenPairs is a special designation for "Seven Pairs". All of the tiles will
	// be represented as a single group.
	OutTileGroupTypeSevenPairs
	// OutTileGroupTypeThirteenOrphans is a special designation for "Thirteen Orphans". All of
	// the tiles will be represented as a single group.
	OutTileGroupTypeThirteenOrphans
)

func (t OutTileGroupType) String() string {
	switch t {
	case OutTileGroupTypePair:
		return "Pair"
	case OutTileGroupTypeChow:
		return "Chow"
	case OutTileGroupTypePong:
		return "Pong"
	case OutTileGroupTypeKong:
		return "Kong"
	case OutTileGroupTypeConcealedKong:
		return "ConcealedKong"
	case OutTileGroupTypeSevenPairs:
		return "SevenPairs"
	case OutTileGroupTypeThirteenOrphans:
		return "ThirteenOrphans"
	}
	glog.Errorf("Unhandled OutTileGroupType %d\n", t)
	return "?"
}

// OutTileGroup represents a group of tiles that make up an Out plan.
type OutTileGroup struct {
	// tiles is the sorted list of tiles that make up the group.
	tiles     domain.Tiles
	groupType OutTileGroupType
}

// UpgradeToKong upgrades the group from Pong to Kong as a result of additional Kong. It is an error
// to call this if the group cannot be upgraded with the given tile.
func (g *OutTileGroup) UpgradeToKong(tile *domain.Tile) {
	if g.groupType != OutTileGroupTypePong || domain.CompareTiles(g.tiles[0], tile) != 0 {
		panic(fmt.Errorf("UpgradeToKong failed for group %s, tile %s", g, tile))
	}
	g.tiles = append(g.tiles, tile)
	g.groupType = OutTileGroupTypeKong
}

// GetTiles ...
func (g *OutTileGroup) GetTiles() domain.Tiles {
	return g.tiles
}

// GetGroupType ...
func (g *OutTileGroup) GetGroupType() OutTileGroupType {
	return g.groupType
}

// String ...
func (g *OutTileGroup) String() string {
	return fmt.Sprintf("%s: %v", g.groupType, g.tiles)
}

// NewOutTileGroup creates a new OutTileGroup with the given parameters. The input tiles is not
// copied, and will be modified by sorting.
func NewOutTileGroup(tiles domain.Tiles, groupType OutTileGroupType) *OutTileGroup {
	sort.Sort(tiles)
	return &OutTileGroup{tiles: tiles, groupType: groupType}
}

// OutTileGroups is a slice of OutTileGroup.
type OutTileGroups []*OutTileGroup

// Len ... (implements sort.Interface)
func (groups OutTileGroups) Len() int {
	return len(groups)
}

// Swap ... (implements sort.Interface)
func (groups OutTileGroups) Swap(i, j int) {
	groups[i], groups[j] = groups[j], groups[i]
}

// Less ... (implements sort.Interface)
func (groups OutTileGroups) Less(i, j int) bool {
	res := domain.CompareTiles(groups[i].GetTiles()[0], groups[j].GetTiles()[0])
	if res != 0 {
		return res < 0
	}
	return groups[i].GetGroupType() < groups[j].GetGroupType()
}

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
func NewOutPlanCalculator(suits []*domain.Suit, h *domain.Hand, discardTile *domain.Tile, meldedGroups OutTileGroups) *OutPlanCalculator {
	inventory := make(tileInventory)
	for _, s := range suits {
		inventory[s] = make([][]*domain.Tile, s.GetSize())
	}
	allTiles := h.GetTiles()
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

	meldedGroupsCopy := append(meldedGroups)
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
