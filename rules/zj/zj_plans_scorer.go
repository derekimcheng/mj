package zj

import (
	"errors"
	"fmt"
	"github.com/derekimcheng/mj/domain"
	"github.com/derekimcheng/mj/rules"
	"github.com/derekimcheng/mj/util"
	"github.com/golang/glog"
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
	allSequences,
	concealedHand,
	noTerminals,
	oneSuit,
	nineGates,
	valueHonor,
	threeDragons,
	fourWinds,
	allTriplets,
	concealedTripletsAndKongs,
	identicalSets,
	threeSimilarSequences,
	threeSimilarTriplets,
	consecutiveSets,
	terminals,
}

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

// 2.0 One-Suit patterns

// 2.1.1 Mixed One-Suit (混一色) : 40
// 2.1.2 Pure One-Suit (清一色) : 80
// Also for optimization:
// 3.4 All Honors (字一色) : 320
func oneSuit(plan rules.OutPlan, context *rules.OutPlanScoringContext) []*rules.Pattern {
	allGroups := append(plan.GetHandGroups(), plan.GetMeldedGroups()...)
	var tilesToCheck domain.Tiles
	for _, group := range allGroups {
		if group.GetGroupType() == rules.TileGroupTypeThirteenOrphans {
			return nil
		}
		if group.GetGroupType() == rules.TileGroupTypeSevenPairs {
			tilesToCheck = append(tilesToCheck, getTilesToCheckForSevenPairs(group)...)
		} else {
			tilesToCheck = append(tilesToCheck, group.GetTiles()[0])
		}
	}
	suitCount, hasHonorTiles := oneSuitHonorHelper(tilesToCheck)
	switch suitCount {
	case noSimpleSuits:
		if hasHonorTiles {
			return []*rules.Pattern{rules.NewPattern("字一色", 320)}
		}
		glog.V(2).Infof("Detected invalid plan having neither simple nor honor tiles: %s\n", plan)
		return nil
	case oneSimpleSuit:
		if hasHonorTiles {
			return []*rules.Pattern{rules.NewPattern("混一色", 40)}
		}
		return []*rules.Pattern{rules.NewPattern("清一色", 80)}
	case moreThanOneSimpleSuits:
		return nil
	}
	panic(errors.New("Unreachable code"))
}

// 2.2 Nine Gates (九蓮寶燈) : 480
func nineGates(plan rules.OutPlan, context *rules.OutPlanScoringContext) []*rules.Pattern {
	if len(plan.GetMeldedGroups()) > 0 {
		return nil
	}
	handGroups := plan.GetHandGroups()
	var tilesToCheck domain.Tiles
	for _, group := range handGroups {
		if group.GetGroupType() == rules.TileGroupTypeThirteenOrphans ||
			group.GetGroupType() == rules.TileGroupTypeSevenPairs {
			return nil
		}
		tilesToCheck = append(tilesToCheck, group.GetTiles()[0])
	}
	suitCount, hasHonorTiles := oneSuitHonorHelper(tilesToCheck)
	if hasHonorTiles || suitCount != oneSimpleSuit {
		return nil
	}
	// Nine Gates must wait from a 1112345678999 hand. To compute what was in the hand, subtract
	// the outTileSource from the tiles in the plan.
	suitSize := handGroups[0].GetTiles()[0].GetSuit().GetSize()
	counts := make([]int, suitSize)
	for _, group := range handGroups {
		for _, tile := range group.GetTiles() {
			counts[tile.GetOrdinal()]++
		}
	}
	counts[context.OutTileSource.Tile.GetOrdinal()]--

	for i, count := range counts {
		expectedCount := 1
		if i == 0 || i == suitSize-1 {
			expectedCount = 3
		}
		if count != expectedCount {
			return nil
		}
	}
	return []*rules.Pattern{rules.NewPattern("九蓮寶燈", 480)}
}

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

// 6.0 Similar Sets

// 6.1 Three Similar Sequences (三色同順) : 35
func threeSimilarSequences(plan rules.OutPlan, context *rules.OutPlanScoringContext) []*rules.Pattern {
	allGroups := append(plan.GetHandGroups(), plan.GetMeldedGroups()...)
	// Map from ordinal to set of suits seen for chow groups.
	chowGroupCounts := make(map[int]map[*domain.Suit]struct{})
	for _, group := range allGroups {
		if group.GetGroupType() != rules.TileGroupTypeChow {
			continue
		}
		// A chow group can be uniquely determined by the first tile of the sequence.
		headTile := group.GetTiles()[0]
		suitSet, found := chowGroupCounts[headTile.GetOrdinal()]
		if !found {
			suitSet = make(map[*domain.Suit]struct{})
			chowGroupCounts[headTile.GetOrdinal()] = suitSet
		}
		suitSet[headTile.GetSuit()] = struct{}{}
	}
	for _, suitSet := range chowGroupCounts {
		if len(suitSet) == 3 {
			return []*rules.Pattern{rules.NewPattern("三色同順", 35)}
		}
	}
	return nil
}

// 6.2.1 Small Three Similar Triplets (三色小同刻) : 30
// 6.2.2 Three Similar Triplets (三色同刻) : 120
func threeSimilarTriplets(plan rules.OutPlan, context *rules.OutPlanScoringContext) []*rules.Pattern {
	allGroups := append(plan.GetHandGroups(), plan.GetMeldedGroups()...)
	// Map from ordinal to map of suits to either a value of 1 or 2. 1 indicates there is only pair
	// of that suit, 2 indicates there is at least one "kan" of that suit.
	similarTripletsPoints := make(map[int]map[*domain.Suit]int)
	for _, group := range allGroups {
		tile := group.GetTiles()[0]
		if tile.GetSuit().GetSuitType() != domain.SuitTypeSimple {
			continue
		}
		points := 0
		if group.GetGroupType() == rules.TileGroupTypePair {
			points = 1
		} else if group.IsKanType() {
			points = 2
		} else {
			continue
		}
		pointsMap, found := similarTripletsPoints[tile.GetOrdinal()]
		if !found {
			pointsMap = make(map[*domain.Suit]int)
			similarTripletsPoints[tile.GetOrdinal()] = pointsMap
		}
		pointsMap[tile.GetSuit()] = util.MaxInt(pointsMap[tile.GetSuit()], points)
	}
	for _, pointsMap := range similarTripletsPoints {
		totalPoints := 0
		for _, points := range pointsMap {
			totalPoints += points
		}
		if totalPoints >= 6 {
			return []*rules.Pattern{rules.NewPattern("三色同刻", 120)}
		} else if totalPoints >= 5 {
			return []*rules.Pattern{rules.NewPattern("三色小同刻", 30)}
		}
	}
	return nil
}

// 7.0 Consecutive Sets

// 7.1 Nine-Tile Straight (一氣通貫) : 40
// 7.2.1 Three Consecutive Triplets (三連刻) : 100
// 7.2.2 Four Consecutive Triplets (四連刻) : 200
func consecutiveSets(plan rules.OutPlan, context *rules.OutPlanScoringContext) []*rules.Pattern {
	allGroups := append(plan.GetHandGroups(), plan.GetMeldedGroups()...)
	// Map from suit to ordinal array where each entry indicates presence of a chow/kan group.
	hasChowGroupMap := make(map[*domain.Suit][]bool)
	hasKanGroupMap := make(map[*domain.Suit][]bool)
	for _, group := range allGroups {
		headTile := group.GetTiles()[0]
		headTileSuit := headTile.GetSuit()
		if group.GetGroupType() != rules.TileGroupTypeChow {
			suitChowGroup, found := hasChowGroupMap[headTileSuit]
			if !found {
				suitChowGroup = make([]bool, headTileSuit.GetSize())
				hasChowGroupMap[headTileSuit] = suitChowGroup
			}
			suitChowGroup[headTile.GetOrdinal()] = true
		} else if group.IsKanType() {
			suitKanGroup, found := hasKanGroupMap[headTileSuit]
			if !found {
				suitKanGroup = make([]bool, headTileSuit.GetSize())
				hasKanGroupMap[headTileSuit] = suitKanGroup
			}
			suitKanGroup[headTile.GetOrdinal()] = true
		}
	}

	for _, suitChowGroup := range hasChowGroupMap {
		// Presence of 1, 4, 7 heads in chow group -> indices 0, 3, 6
		if suitChowGroup[0] && suitChowGroup[3] && suitChowGroup[6] {
			return []*rules.Pattern{rules.NewPattern("一氣通貫", 40)}
		}
	}

	// Once we encounter a 3-consecutive kan, we can stop the search since the hand doesn't allow
	// another 4-consecutive kan to occur anyway.
	for _, suitKanGroup := range hasKanGroupMap {
		consecutiveKans := 0
		for _, hasKan := range suitKanGroup {
			if hasKan {
				consecutiveKans++
			} else {
				if consecutiveKans == 3 {
					return []*rules.Pattern{rules.NewPattern("三連刻", 100)}
				} else if consecutiveKans == 4 {
					return []*rules.Pattern{rules.NewPattern("四連刻", 200)}
				}
				consecutiveKans = 0
			}
		}
	}
	return nil
}

// 8.0 Terminals
// 8.1.1 Mixed Lesser Terminals (混全帶么) : 40
// 8.1.2 Pure Lesser Terminals (純全帶么) : 50
// 8.1.3 Mixed Greater Terminals (混么九) : 100
// 8.1.4 Pure Greater Terminals (清么九) : 400
func terminals(plan rules.OutPlan, context *rules.OutPlanScoringContext) []*rules.Pattern {
	allGroups := append(plan.GetHandGroups(), plan.GetMeldedGroups()...)
	hasChowGroup := false
	var tilesToCheck domain.Tiles
	for _, group := range allGroups {
		tiles := group.GetTiles()
		headTile := tiles[0]
		if group.IsKanType() || group.GetGroupType() == rules.TileGroupTypePair {
			tilesToCheck = append(tilesToCheck, headTile)
		} else if group.GetGroupType() == rules.TileGroupTypeChow {
			hasChowGroup = true
			if headTile.GetOrdinal() == 0 {
				tilesToCheck = append(tilesToCheck, headTile)
			} else {
				tilesToCheck = append(tilesToCheck, tiles[len(tiles)-1])
			}
		} else if group.GetGroupType() == rules.TileGroupTypeSevenPairs {
			tilesToCheck = append(tilesToCheck, getTilesToCheckForSevenPairs(group)...)
		} else {
			return nil
		}
	}

	hasTerminal, hasNonTerminal, hasHonor := terminalTileHelper(tilesToCheck)
	if hasNonTerminal || !hasTerminal {
		return nil
	}
	if hasChowGroup {
		if hasHonor {
			return []*rules.Pattern{rules.NewPattern("混全帶么", 40)}
		}
		return []*rules.Pattern{rules.NewPattern("純全帶么", 50)}
	}
	if hasHonor {
		return []*rules.Pattern{rules.NewPattern("混么九", 100)}
	}
	return []*rules.Pattern{rules.NewPattern("清么九", 400)}
}

// getTilesToCheckForSevenPairs returns for the given Seven Pairs tile group, list of tiles to
// check using the helper functions provided below.
func getTilesToCheckForSevenPairs(group *rules.TileGroup) domain.Tiles {
	// We only need to check every other tile.
	var tiles domain.Tiles
	groupTiles := group.GetTiles()
	for i := 0; i < len(groupTiles); i += 2 {
		tiles = append(tiles, groupTiles[i])
	}
	return tiles
}

// terminalTileHelper returns three bools given a set of tiles:
// The first bool indicates where there is at least one terminal tile.
// The second bool indicates whether there is at least one non-terminal tile.
// The third bool indicates whether there is at least one honor tile.
func terminalTileHelper(tiles domain.Tiles) (bool, bool, bool) {
	hasTerminal := false
	hasNonTerminal := false
	hasHonor := false
	for _, tile := range tiles {
		suitType := tile.GetSuit().GetSuitType()
		if suitType == domain.SuitTypeHonor {
			hasHonor = true
		} else if suitType == domain.SuitTypeSimple {
			if tile.IsTerminal() {
				hasTerminal = true
			} else {
				hasNonTerminal = true
			}
		}
	}
	return hasTerminal, hasNonTerminal, hasHonor
}

type simpleSuitCount int

const (
	noSimpleSuits simpleSuitCount = iota
	oneSimpleSuit
	moreThanOneSimpleSuits
)

// oneSuitHonorHelper returns:
// a simple suit name, if there is exactly one simple suit amongst all tiles. Empty otherwise.
// a bool indicating whether there is at least one honor tile AND there's at most one simple suit.
func oneSuitHonorHelper(tiles domain.Tiles) (simpleSuitCount, bool) {
	hasHonorTiles := false
	var suitName string
	for _, tile := range tiles {
		switch tile.GetSuit().GetSuitType() {
		case domain.SuitTypeSimple:
			if len(suitName) == 0 {
				suitName = tile.GetSuit().GetName()
			} else if suitName != tile.GetSuit().GetName() {
				// Encountered more than one simple suit.
				return moreThanOneSimpleSuits, false
			}
		case domain.SuitTypeHonor:
			hasHonorTiles = true
		}
	}
	if len(suitName) > 0 {
		return oneSimpleSuit, hasHonorTiles
	}
	return noSimpleSuits, hasHonorTiles
}
