package zj

import (
	"github.com/derekimcheng/mj/rules"
)

// 4.0 Triplets and Kong

// 4.1 All Triplets (對對和) : 30
func allTriplets(plan rules.OutPlan, context *rules.OutPlanScoringContext) []*rules.Pattern {
	allGroups := append(plan.GetHandGroups(), plan.GetMeldedGroups()...)
	for _, group := range allGroups {
		if !group.IsKanType() && group.GetGroupType() != rules.TileGroupTypePair {
			return nil
		}
	}
	return []*rules.Pattern{rules.NewPattern("對對和", 30)}
}

// 4.2.1 Two Concealed Triplets (二暗刻) : 5
// 4.2.2 Three Concealed Triplets (三暗刻) : 30
// 4.2.3 Four Concealed Triplets (四暗刻) : 125
// 4.3.1 One Kong (一槓) : 5
// 4.3.2 Two Kong (二槓) : 20
// 4.3.3 Three Kong (三槓) : 120
// 4.3.4 Four Kong (四槓) : 480
func concealedTripletsAndKongs(plan rules.OutPlan, context *rules.OutPlanScoringContext) []*rules.Pattern {
	numConcealedTriplets := 0
	numKongs := 0
	for _, group := range plan.GetMeldedGroups() {
		switch group.GetGroupType() {
		case rules.TileGroupTypeConcealedKong:
			numConcealedTriplets++
			numKongs++
		case rules.TileGroupTypeKong:
			numKongs++
		}
	}
	for _, group := range plan.GetHandGroups() {
		if group.GetGroupType() == rules.TileGroupTypePong {
			numConcealedTriplets++
		}
	}

	var patterns []*rules.Pattern
	switch numConcealedTriplets {
	case 2:
		patterns = append(patterns, rules.NewPattern("二暗刻", 5))
	case 3:
		patterns = append(patterns, rules.NewPattern("三暗刻", 30))
	case 4:
		patterns = append(patterns, rules.NewPattern("四暗刻", 125))
	}
	switch numKongs {
	case 1:
		patterns = append(patterns, rules.NewPattern("一槓", 5))
	case 2:
		patterns = append(patterns, rules.NewPattern("二槓", 20))
	case 3:
		patterns = append(patterns, rules.NewPattern("三槓", 120))
	case 4:
		patterns = append(patterns, rules.NewPattern("四槓", 480))
	}
	return patterns
}
