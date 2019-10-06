package zj

import (
	"github.com/derekimcheng/mj/domain"
	"github.com/derekimcheng/mj/rules"
)

// 1.0 Trivial patterns

// 1.1 All Sequences (平和) : 5
func allSequences(plan rules.OutPlan, context *rules.OutPlanScoringContext) []*rules.Pattern {
	allGroups := append(plan.GetHandGroups(), plan.GetMeldedGroups()...)
	for _, group := range allGroups {
		if group.GetGroupType() != rules.TileGroupTypeChow &&
			group.GetGroupType() != rules.TileGroupTypePair {
			return nil
		}
	}
	return []*rules.Pattern{rules.NewPattern("平和", 5)}
}

// 1.2 Concealed Hand (門前清) : 5
func concealedHand(plan rules.OutPlan, context *rules.OutPlanScoringContext) []*rules.Pattern {
	for _, group := range plan.GetMeldedGroups() {
		if group.GetGroupType() != rules.TileGroupTypeConcealedKong {
			return nil
		}
	}
	for _, group := range plan.GetHandGroups() {
		if group.GetGroupType() == rules.TileGroupTypeSevenPairs ||
			group.GetGroupType() == rules.TileGroupTypeThirteenOrphans {
			return nil
		}
	}
	return []*rules.Pattern{rules.NewPattern("門前清", 5)}
}

// 1.3 No Terminals (斷么九) : 5
func noTerminals(plan rules.OutPlan, context *rules.OutPlanScoringContext) []*rules.Pattern {
	allGroups := append(plan.GetHandGroups(), plan.GetMeldedGroups()...)
	for _, group := range allGroups {
		for _, tile := range group.GetTiles() {
			if tile.GetSuit().GetSuitType() != domain.SuitTypeSimple || tile.IsTerminal() {
				return nil
			}

		}
	}
	return []*rules.Pattern{rules.NewPattern("斷么九", 5)}
}
