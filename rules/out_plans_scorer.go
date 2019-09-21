package rules

import (
	"fmt"
	"strings"
)

// ScoredOutPlan is a OutPlan with a score.
type ScoredOutPlan struct {
	Plan       OutPlan
	TotalScore int
	Patterns   Patterns
}

// ScoredOutPlans is a slice of ScoredOutPlan.
type ScoredOutPlans []*ScoredOutPlan

// Len ... (sort.Sort implementation)
func (ps ScoredOutPlans) Len() int {
	return len(ps)
}

// Swap ... (sort.Sort implementation)
func (ps ScoredOutPlans) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

// Less ... (sort.Sort implementation)
func (ps ScoredOutPlans) Less(i, j int) bool {
	if scoreDiff := ps[i].TotalScore - ps[j].TotalScore; scoreDiff != 0 {
		return scoreDiff > 0
	}
	// Assumes Patterns are already sorted. Whoever has the biggest individual pattern wins.
	return ComparePatterns(ps[i].Patterns[0], ps[j].Patterns[0]) > 0
}

// String ...
func (ps ScoredOutPlans) String() string {
	var str string
	str += fmt.Sprintf("Number of possible out plans: %d\n", len(ps))
	for i, plan := range ps {
		str += fmt.Sprintf("  Plan %d:\n", i+1)
		str += fmt.Sprintf("  %s\n", plan.Plan)
		str += fmt.Sprintf("  Total score: %d\n", plan.TotalScore)
		str += fmt.Sprintf("  Patterns:\n")
		for _, pattern := range plan.Patterns {
			str += fmt.Sprintf("    %s - %d\n", pattern.Name, pattern.Score)
		}
	}
	return str
}

// NewScoredOutPlan creates a new ScoredOutPlan.
func NewScoredOutPlan(plan OutPlan, totalScore int, patterns Patterns) *ScoredOutPlan {
	return &ScoredOutPlan{Plan: plan, TotalScore: totalScore, Patterns: patterns}
}

// Pattern is an opaque structure representing a scoring pattern.
type Pattern struct {
	Name  string
	Score int
}

// NewPattern creates a new Pattern.
func NewPattern(name string, score int) *Pattern {
	return &Pattern{Name: name, Score: score}
}

// ComparePatterns returns -1 if p1 precedes p2, 1 if p2 precedes p1, and 0 if p1 and p2 are
// equivalent in ordering.
func ComparePatterns(p1, p2 *Pattern) int {
	if scoreDiff := p1.Score - p2.Score; scoreDiff != 0 {
		return scoreDiff
	}
	return strings.Compare(p1.Name, p2.Name)
}

// Patterns is a list of Patterns.
type Patterns []*Pattern

// Len ... (sort.Sort implementation)
func (ps Patterns) Len() int {
	return len(ps)
}

// Swap ... (sort.Sort implementation)
func (ps Patterns) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

// Less ... (sort.Sort implementation)
func (ps Patterns) Less(i, j int) bool {
	return ComparePatterns(ps[i], ps[j]) > 0
}

// OutPlanScoringContext is a struct that contains context for scoring besides the out plan itself.
type OutPlanScoringContext struct {
	OutTileSource   *OutTileSource
	PlayerGameState *PlayerGameState
}

// NewOutPlanScoringContext ...
func NewOutPlanScoringContext(outTileSource *OutTileSource,
	playerGameState *PlayerGameState) *OutPlanScoringContext {
	return &OutPlanScoringContext{OutTileSource: outTileSource, PlayerGameState: playerGameState}
}

// OutPlansScorer scores a list of plans according to the implementation's rules.
type OutPlansScorer interface {
	// ScoreOutPlans scores the given plans and the context and returns them in descending
	// score order.
	ScoreOutPlans(plans OutPlans, context *OutPlanScoringContext) ScoredOutPlans
}
