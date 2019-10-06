package zj

import (
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
	context *rules.OutPlanScoringContext) rules.ScoredOutPlans {
	var scoredPlans rules.ScoredOutPlans
	for _, plan := range plans {
		scoredPlans = append(scoredPlans, s.scoreOutPlan(plan, context))
	}

	sort.Sort(scoredPlans)
	return scoredPlans
}

func (s *OutPlansScorer) scoreOutPlan(plan rules.OutPlan, context *rules.OutPlanScoringContext) *rules.ScoredOutPlan {
	var patterns rules.Patterns
	for _, matchPattern := range matchPatternFuncList {
		patterns = append(patterns, matchPattern(plan, context)...)
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
type matchPatternFunc func(rules.OutPlan, *rules.OutPlanScoringContext) []*rules.Pattern

var matchPatternFuncList = []matchPatternFunc{
	// 1.0
	allSequences,
	concealedHand,
	noTerminals,
	// 2.0
	oneSuit,
	nineGates,
	// 3.0
	valueHonor,
	threeDragons,
	fourWinds,
	// 4.0
	allTriplets,
	concealedTripletsAndKongs,
	// 5.0
	identicalSets,
	// 6.0
	threeSimilarSequences,
	threeSimilarTriplets,
	// 7.0
	consecutiveSets,
	// 8.0
	terminals,
	// 9.0
	finalDrawOrDiscard,
	winOnKong,
	winOnInitialRound,
	// 10.0
	irregularHands,
}
