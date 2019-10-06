package zj

import (
	"fmt"
	"github.com/derekimcheng/mj/domain"
	"github.com/derekimcheng/mj/rules"
)

// 3.0 Honor Tiles

// 3.1 Value Honor (番牌) : 10 per set
func valueHonor(plan rules.OutPlan, context *rules.OutPlanScoringContext) []*rules.Pattern {
	allGroups := append(plan.GetHandGroups(), plan.GetMeldedGroups()...)
	numHonors := 0
	for _, group := range allGroups {
		if !group.IsKanType() {
			continue
		}
		firstTile := group.GetTiles()[0]
		suit := firstTile.GetSuit()
		if suit.GetSuitType() != domain.SuitTypeHonor {
			continue
		}
		// Prevailing wind
		if rules.IsWindSuit(suit) {
			if firstTile.GetOrdinal() == context.PlayerGameState.GetWindOrdinal() {
				numHonors++
			}
		} else {
			numHonors++
		}
	}
	if numHonors == 0 {
		return nil
	}
	patternName := fmt.Sprintf("番牌 (%d)", numHonors)
	return []*rules.Pattern{rules.NewPattern(patternName, numHonors*10)}
}

// 3.2.1 Small Three Dragons (小三元) : 40
// 3.2.2 Big Three Dragons (大三元) : 130
func threeDragons(plan rules.OutPlan, context *rules.OutPlanScoringContext) []*rules.Pattern {
	allGroups := append(plan.GetHandGroups(), plan.GetMeldedGroups()...)
	numKans := 0
	numPairs := 0
	for _, group := range allGroups {
		firstTile := group.GetTiles()[0]
		suit := firstTile.GetSuit()
		if suit.GetSuitType() != domain.SuitTypeHonor || rules.IsWindSuit(suit) {
			continue
		}
		if group.IsKanType() {
			numKans++
		} else {
			// Else assume it is a pair since honor tiles can only be kans or pairs.
			numPairs++
		}
	}

	if numKans == 2 && numPairs == 1 {
		return []*rules.Pattern{rules.NewPattern("小三元", 40)}
	} else if numKans == 3 {
		return []*rules.Pattern{rules.NewPattern("大三元", 130)}
	}
	return nil
}

// 3.3.1 Small Three Winds (小三風) : 30
// 3.3.2 Big Three Winds (大三風) : 120
// 3.3.3 Small Four Winds (小四喜) : 320
// 3.3.4 Big Four Winds (大四喜) : 400
func fourWinds(plan rules.OutPlan, context *rules.OutPlanScoringContext) []*rules.Pattern {
	allGroups := append(plan.GetHandGroups(), plan.GetMeldedGroups()...)
	numKans := 0
	numPairs := 0
	for _, group := range allGroups {
		firstTile := group.GetTiles()[0]
		suit := firstTile.GetSuit()
		if !rules.IsWindSuit(suit) {
			continue
		}
		if group.IsKanType() {
			numKans++
		} else {
			// Else assume it is a pair since honor tiles can only be kans or pairs.
			numPairs++
		}
	}

	if numKans == 2 && numPairs == 1 {
		return []*rules.Pattern{rules.NewPattern("小三風", 30)}
	} else if numKans == 3 {
		if numPairs == 0 {
			return []*rules.Pattern{rules.NewPattern("大三風", 120)}
		}
		return []*rules.Pattern{rules.NewPattern("小四喜", 320)}
	} else if numKans == 4 {
		return []*rules.Pattern{rules.NewPattern("大四喜", 400)}
	}
	return nil
}

// 3.4 All Honors (字一色) : 320 - covered by oneSuit().