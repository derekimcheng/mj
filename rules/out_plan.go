package rules

import (
	"fmt"
	"github.com/golang/glog"
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
	handGroups   OutTileGroups
	meldedGroups OutTileGroups
}

// GetHandGroups ...
func (p OutPlan) GetHandGroups() OutTileGroups {
	return p.handGroups
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
func NewOutPlan(handGroups OutTileGroups, meldedGroups OutTileGroups) OutPlan {
	sort.Sort(handGroups)
	return OutPlan{handGroups: handGroups, meldedGroups: meldedGroups}
}

// OutTileGroupType represents a type of Out tile group
type OutTileGroupType int

const (
	// OutTileGroupTypePair represents a pair of same suit+ordinal tiles. With the exception of
	// special hands, each out plan must contain exactly one Pair group.
	OutTileGroupTypePair OutTileGroupType = iota
	// OutTileGroupTypePong represents a triplet of the same suit+ordinal tiles.
	OutTileGroupTypePong
	// OutTileGroupTypeChow represents a triplet of consecutive tiles.
	OutTileGroupTypeChow
	// OutTileGroupTypeKong represents a quadruplet of the same suit+ordinal tiles. This is only
	// available in the meld area.
	OutTileGroupTypeKong
	// OutTileGroupTypeConcealedKong represents a concealed quadruplet of the same suit+ordinal
	// tiles. This is only available in the meld area.
	OutTileGroupTypeConcealedKong
	// OutTileGroupTypeSevenPairs is a special designation for "Seven Pairs". All of the tiles will
	// be represented as a single group.
	OutTileGroupTypeSevenPairs
	// OutTileGroupTypeThirteenOrphans is a special designation for "Thirteen Orphans". All of
	// the tiles will be represented as a single group.
	OutTileGroupTypeThirteenOrphans
)

func (t OutTileGroupType) String() string {
	switch t {
	case OutTileGroupTypePair:
		return "Pair"
	case OutTileGroupTypeChow:
		return "Chow"
	case OutTileGroupTypePong:
		return "Pong"
	case OutTileGroupTypeKong:
		return "Kong"
	case OutTileGroupTypeConcealedKong:
		return "ConcealedKong"
	case OutTileGroupTypeSevenPairs:
		return "SevenPairs"
	case OutTileGroupTypeThirteenOrphans:
		return "ThirteenOrphans"
	}
	glog.Errorf("Unhandled OutTileGroupType %d\n", t)
	return "?"
}