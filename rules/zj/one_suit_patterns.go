package zj

import (
	"errors"
	"github.com/golang/glog"
	"github.com/derekimcheng/mj/domain"
	"github.com/derekimcheng/mj/rules"
)

// 2.0 One-Suit patterns

// 2.1.1 Mixed One-Suit (混一色) : 40
// 2.1.2 Pure One-Suit (清一色) : 80
// Also for optimization:
// 3.4 All Honors (字一色) : 320
func oneSuit(plan rules.OutPlan, context *rules.OutPlanScoringContext) []*rules.Pattern {
	allGroups := append(plan.GetHandGroups(), plan.GetMeldedGroups()...)
	var tilesToCheck domain.Tiles
	for _, group := range allGroups {
		if group.GetGroupType() == rules.TileGroupTypeThirteenOrphans {
			return nil
		}
		if group.GetGroupType() == rules.TileGroupTypeSevenPairs {
			tilesToCheck = append(tilesToCheck, getTilesToCheckForSevenPairs(group)...)
		} else {
			tilesToCheck = append(tilesToCheck, group.GetTiles()[0])
		}
	}
	suitCount, hasHonorTiles := oneSuitHonorHelper(tilesToCheck)
	switch suitCount {
	case noSimpleSuits:
		if hasHonorTiles {
			return []*rules.Pattern{rules.NewPattern("字一色", 320)}
		}
		glog.V(2).Infof("Detected invalid plan having neither simple nor honor tiles: %s\n", plan)
		return nil
	case oneSimpleSuit:
		if hasHonorTiles {
			return []*rules.Pattern{rules.NewPattern("混一色", 40)}
		}
		return []*rules.Pattern{rules.NewPattern("清一色", 80)}
	case moreThanOneSimpleSuits:
		return nil
	}
	panic(errors.New("Unreachable code"))
}

// 2.2 Nine Gates (九蓮寶燈) : 480
func nineGates(plan rules.OutPlan, context *rules.OutPlanScoringContext) []*rules.Pattern {
	if len(plan.GetMeldedGroups()) > 0 {
		return nil
	}
	handGroups := plan.GetHandGroups()
	var tilesToCheck domain.Tiles
	for _, group := range handGroups {
		if group.GetGroupType() == rules.TileGroupTypeThirteenOrphans ||
			group.GetGroupType() == rules.TileGroupTypeSevenPairs {
			return nil
		}
		tilesToCheck = append(tilesToCheck, group.GetTiles()[0])
	}
	suitCount, hasHonorTiles := oneSuitHonorHelper(tilesToCheck)
	if hasHonorTiles || suitCount != oneSimpleSuit {
		return nil
	}
	// Nine Gates must wait from a 1112345678999 hand. To compute what was in the hand, subtract
	// the outTileSource from the tiles in the plan.
	suitSize := handGroups[0].GetTiles()[0].GetSuit().GetSize()
	counts := make([]int, suitSize)
	for _, group := range handGroups {
		for _, tile := range group.GetTiles() {
			counts[tile.GetOrdinal()]++
		}
	}
	counts[context.OutTileSource.Tile.GetOrdinal()]--

	for i, count := range counts {
		expectedCount := 1
		if i == 0 || i == suitSize-1 {
			expectedCount = 3
		}
		if count != expectedCount {
			return nil
		}
	}
	return []*rules.Pattern{rules.NewPattern("九蓮寶燈", 480)}
}
