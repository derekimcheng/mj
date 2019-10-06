package zj

import (
	"github.com/derekimcheng/mj/rules"
)

// 10.0 Irregular Hands
// 10.1 Thirteen Terminals (十三么九) : 160
// 10.2 Seven Pairs (七對子) : 30
func irregularHands(plan rules.OutPlan, context *rules.OutPlanScoringContext) []*rules.Pattern {
	handGroups := plan.GetHandGroups()
	if len(handGroups) != 1 {
		return nil
	}
	switch handGroups[0].GetGroupType() {
	case rules.TileGroupTypeThirteenOrphans:
		return []*rules.Pattern{rules.NewPattern("十三么九", 160)}
	case rules.TileGroupTypeSevenPairs:
		return []*rules.Pattern{rules.NewPattern("七對子", 30)}
	}
	return nil
}
