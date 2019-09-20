package zj

import (
	"github.com/derekimcheng/mj/domain"
	"github.com/derekimcheng/mj/rules"
	"sort"
)

// OutPlansScorer is an implementation of rules.OutPLansScorer based on ZJ rules.
// The implementation assumes each plan contains a valid combination of tiles. Any invalid
// combination may result in incorrect scoring.
type OutPlansScorer struct{}

// NewOutPlansScorer creates a new OutPlansScorer.
func NewOutPlansScorer() *OutPlansScorer {
	return &OutPlansScorer{}
}

// ScoreOutPlans ... (rules.OutPlansScorer implementation)
func (s *OutPlansScorer) ScoreOutPlans(plans rules.OutPlans,
	outTileSource *rules.OutTileSource) rules.ScoredOutPlans {
	var scoredPlans rules.ScoredOutPlans
	for _, plan := range plans {
		scoredPlans = append(scoredPlans, s.scoreOutPlan(plan, outTileSource))
	}

	sort.Sort(scoredPlans)
	return scoredPlans
}

func (s *OutPlansScorer) scoreOutPlan(plan rules.OutPlan, outTileSource *rules.OutTileSource) *rules.ScoredOutPlan {
	var patterns rules.Patterns
	for _, matchPattern := range matchPatternFuncList {
		patterns = append(patterns, matchPattern(plan, outTileSource)...)
	}

	if len(patterns) == 0 {
		patterns = append(patterns, rules.NewPattern("雞和", 1))
	}

	sort.Sort(patterns)
	totalScore := 0
	for _, pattern := range patterns {
		totalScore += pattern.Score
	}

	return rules.NewScoredOutPlan(plan, totalScore, patterns)
}

// TODO: maybe this can be common?
type matchPatternFunc func(rules.OutPlan, *rules.OutTileSource) []*rules.Pattern

var matchPatternFuncList = []matchPatternFunc{
	allSequences,
	concealedHand,
	noTerminals,
}

// 1.0 Trivial patterns

// 1.1 All Sequences (平和) : 5
func allSequences(plan rules.OutPlan, outTileSource *rules.OutTileSource) []*rules.Pattern {
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
func concealedHand(plan rules.OutPlan, outTileSource *rules.OutTileSource) []*rules.Pattern {
	if len(plan.GetMeldedGroups()) == 0 {
		return []*rules.Pattern{rules.NewPattern("門前清", 5)}
	}
	return nil
}

// 1.3 No Terminals (斷么九) : 5
func noTerminals(plan rules.OutPlan, outTileSource *rules.OutTileSource) []*rules.Pattern {
	allGroups := append(plan.GetHandGroups(), plan.GetMeldedGroups()...)
	for _, group := range allGroups {
		for _, tile := range group.GetTiles() {
			if tile.GetSuit().GetSuitType() != domain.SuitTypeSimple || tile.GetOrdinal() == 0 || tile.GetOrdinal() == 8 {
				return nil
			}

		}
	}
	return []*rules.Pattern{rules.NewPattern("斷么九", 5)}
}
