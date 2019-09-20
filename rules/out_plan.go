package rules

import (
	"fmt"
	"sort"
	"strings"
)

// OutPlans is a slice of OutPlan.
// TODO: OutPlans should be sorted according to score / type /etc.
type OutPlans []OutPlan

// String ...
func (plans OutPlans) String() string {
	if len(plans) == 0 {
		return ""
	}
	ret := ""
	for i, plan := range plans {
		if i > 0 {
			ret += ", "
		}
		ret += fmt.Sprintf("Plan %d: %s", i+1, plan)
	}
	return ret
}

// OutPlan consists of a strategy of declaring an Out using a combination of tile groups in the
// hand and meld area.
type OutPlan struct {
	handGroups   TileGroups
	meldedGroups TileGroups
}

// GetHandGroups ...
func (p OutPlan) GetHandGroups() TileGroups {
	return p.handGroups
}

// GetMeldedGroups ...
func (p OutPlan) GetMeldedGroups() TileGroups {
	return p.meldedGroups
}

// String ...
func (p OutPlan) String() string {
	handGroupStrs := []string{}
	for _, handGroup := range p.handGroups {
		handGroupStrs = append(handGroupStrs, handGroup.String())
	}
	hand := fmt.Sprintf("Hand: %s", strings.Join(handGroupStrs, ", "))

	meldedGroupStrs := []string{}
	for _, meldedGroup := range p.meldedGroups {
		meldedGroupStrs = append(meldedGroupStrs, meldedGroup.String())
	}
	melded := fmt.Sprintf("Melded: %s", strings.Join(meldedGroupStrs, ", "))

	return hand + " " + melded
}

// NewOutPlan creates a new OutPlan with the given parameters. The input groups is not
// copied, and will be modified by sorting.
func NewOutPlan(handGroups TileGroups, meldedGroups TileGroups) OutPlan {
	sort.Sort(handGroups)
	return OutPlan{handGroups: handGroups, meldedGroups: meldedGroups}
}