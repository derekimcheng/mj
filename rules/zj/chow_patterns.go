package zj

import (
	"github.com/derekimcheng/mj/domain"
	"github.com/derekimcheng/mj/rules"
)

// 5.0 Identical Sets

// 5.1.1 Two Identical Sequences (一般高) : 10
// 5.1.2 Two Identical Sequences Twice (兩般高) : 60
// 5.1.3 Three Identical Sequences (一色三同順) : 120
// 5.1.4 Four Identical Sequences (一色四同順) : 480
func identicalSets(plan rules.OutPlan, context *rules.OutPlanScoringContext) []*rules.Pattern {
	allGroups := append(plan.GetHandGroups(), plan.GetMeldedGroups()...)
	chowGroupCounts := make(map[domain.TileBase]int)
	for _, group := range allGroups {
		if group.GetGroupType() != rules.TileGroupTypeChow {
			continue
		}
		// A chow group can be uniquely determined by the first tile of the sequence.
		headTileBase := group.GetTiles()[0].TileBase
		chowGroupCounts[headTileBase]++
	}

	numTwoIndenticalSeqs := 0
	for _, count := range chowGroupCounts {
		switch count {
		case 2:
			numTwoIndenticalSeqs++
			if numTwoIndenticalSeqs == 2 {
				return []*rules.Pattern{rules.NewPattern("兩般高", 60)}
			}
		case 3:
			return []*rules.Pattern{rules.NewPattern("一色三同順", 120)}
		case 4:
			return []*rules.Pattern{rules.NewPattern("一色四同順", 480)}
		}
	}

	if numTwoIndenticalSeqs == 1 {
		return []*rules.Pattern{rules.NewPattern("一般高", 10)}
	}

	return nil
}