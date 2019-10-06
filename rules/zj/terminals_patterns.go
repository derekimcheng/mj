package zj

import (
	"github.com/derekimcheng/mj/domain"
	"github.com/derekimcheng/mj/rules"
)

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