package zj

import (
	"github.com/derekimcheng/mj/domain"
	"github.com/derekimcheng/mj/rules"
	"github.com/derekimcheng/mj/util"
)

// 6.0 Similar Sets

// 6.1 Three Similar Sequences (三色同順) : 35
func threeSimilarSequences(plan rules.OutPlan, context *rules.OutPlanScoringContext) []*rules.Pattern {
	allGroups := append(plan.GetHandGroups(), plan.GetMeldedGroups()...)
	// Map from ordinal to set of suits seen for chow groups.
	chowGroupCounts := make(map[int]map[*domain.Suit]struct{})
	for _, group := range allGroups {
		if group.GetGroupType() != rules.TileGroupTypeChow {
			continue
		}
		// A chow group can be uniquely determined by the first tile of the sequence.
		headTile := group.GetTiles()[0]
		suitSet, found := chowGroupCounts[headTile.GetOrdinal()]
		if !found {
			suitSet = make(map[*domain.Suit]struct{})
			chowGroupCounts[headTile.GetOrdinal()] = suitSet
		}
		suitSet[headTile.GetSuit()] = struct{}{}
	}
	for _, suitSet := range chowGroupCounts {
		if len(suitSet) == 3 {
			return []*rules.Pattern{rules.NewPattern("三色同順", 35)}
		}
	}
	return nil
}

// 6.2.1 Small Three Similar Triplets (三色小同刻) : 30
// 6.2.2 Three Similar Triplets (三色同刻) : 120
func threeSimilarTriplets(plan rules.OutPlan, context *rules.OutPlanScoringContext) []*rules.Pattern {
	allGroups := append(plan.GetHandGroups(), plan.GetMeldedGroups()...)
	// Map from ordinal to map of suits to either a value of 1 or 2. 1 indicates there is only pair
	// of that suit, 2 indicates there is at least one "kan" of that suit.
	similarTripletsPoints := make(map[int]map[*domain.Suit]int)
	for _, group := range allGroups {
		tile := group.GetTiles()[0]
		if tile.GetSuit().GetSuitType() != domain.SuitTypeSimple {
			continue
		}
		points := 0
		if group.GetGroupType() == rules.TileGroupTypePair {
			points = 1
		} else if group.IsKanType() {
			points = 2
		} else {
			continue
		}
		pointsMap, found := similarTripletsPoints[tile.GetOrdinal()]
		if !found {
			pointsMap = make(map[*domain.Suit]int)
			similarTripletsPoints[tile.GetOrdinal()] = pointsMap
		}
		pointsMap[tile.GetSuit()] = util.MaxInt(pointsMap[tile.GetSuit()], points)
	}
	for _, pointsMap := range similarTripletsPoints {
		totalPoints := 0
		for _, points := range pointsMap {
			totalPoints += points
		}
		if totalPoints >= 6 {
			return []*rules.Pattern{rules.NewPattern("三色同刻", 120)}
		} else if totalPoints >= 5 {
			return []*rules.Pattern{rules.NewPattern("三色小同刻", 30)}
		}
	}
	return nil
}
