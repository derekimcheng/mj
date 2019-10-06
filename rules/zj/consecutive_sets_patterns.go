package zj

import (
	"github.com/derekimcheng/mj/domain"
	"github.com/derekimcheng/mj/rules"
)

// 7.0 Consecutive Sets

// 7.1 Nine-Tile Straight (一氣通貫) : 40
// 7.2.1 Three Consecutive Triplets (三連刻) : 100
// 7.2.2 Four Consecutive Triplets (四連刻) : 200
func consecutiveSets(plan rules.OutPlan, context *rules.OutPlanScoringContext) []*rules.Pattern {
	allGroups := append(plan.GetHandGroups(), plan.GetMeldedGroups()...)
	// Map from suit to ordinal array where each entry indicates presence of a chow/kan group.
	hasChowGroupMap := make(map[*domain.Suit][]bool)
	hasKanGroupMap := make(map[*domain.Suit][]bool)
	for _, group := range allGroups {
		headTile := group.GetTiles()[0]
		headTileSuit := headTile.GetSuit()
		if group.GetGroupType() != rules.TileGroupTypeChow {
			suitChowGroup, found := hasChowGroupMap[headTileSuit]
			if !found {
				suitChowGroup = make([]bool, headTileSuit.GetSize())
				hasChowGroupMap[headTileSuit] = suitChowGroup
			}
			suitChowGroup[headTile.GetOrdinal()] = true
		} else if group.IsKanType() {
			suitKanGroup, found := hasKanGroupMap[headTileSuit]
			if !found {
				suitKanGroup = make([]bool, headTileSuit.GetSize())
				hasKanGroupMap[headTileSuit] = suitKanGroup
			}
			suitKanGroup[headTile.GetOrdinal()] = true
		}
	}

	for _, suitChowGroup := range hasChowGroupMap {
		// Presence of 1, 4, 7 heads in chow group -> indices 0, 3, 6
		if suitChowGroup[0] && suitChowGroup[3] && suitChowGroup[6] {
			return []*rules.Pattern{rules.NewPattern("一氣通貫", 40)}
		}
	}

	// Once we encounter a 3-consecutive kan, we can stop the search since the hand doesn't allow
	// another 4-consecutive kan to occur anyway.
	for _, suitKanGroup := range hasKanGroupMap {
		consecutiveKans := 0
		for _, hasKan := range suitKanGroup {
			if hasKan {
				consecutiveKans++
			} else {
				if consecutiveKans == 3 {
					return []*rules.Pattern{rules.NewPattern("三連刻", 100)}
				} else if consecutiveKans == 4 {
					return []*rules.Pattern{rules.NewPattern("四連刻", 200)}
				}
				consecutiveKans = 0
			}
		}
	}
	return nil
}
