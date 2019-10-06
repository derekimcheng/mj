package zj

import (
	"github.com/derekimcheng/mj/rules"
)

// 9.0 Incidental bonuses
// 9.1.1 Final Draw (海底撈月) : 10
// 9.1.2 Final Discard (河底撈魚) : 10
func finalDrawOrDiscard(plan rules.OutPlan, context *rules.OutPlanScoringContext) []*rules.Pattern {
	if context.NumRemainingTilesInDeck > 0 {
		return nil
	}
	outTileSourceType := context.OutTileSource.SourceType
	if outTileSourceType == rules.OutTileSourceTypeSelfDrawn ||
		outTileSourceType == rules.OutTileSourceTypeSelfDrawnReplacement {
		return []*rules.Pattern{rules.NewPattern("海底撈月", 10)}
	}
	if outTileSourceType == rules.OutTileSourceTypeDiscard {
		return []*rules.Pattern{rules.NewPattern("河底撈魚", 10)}
	}
	return nil
}

// 9.2 Win on Kong (嶺上開花) : 10
// 9.3 Robbing a Kong (搶槓) : 10
func winOnKong(plan rules.OutPlan, context *rules.OutPlanScoringContext) []*rules.Pattern {
	if context.OutTileSource.SourceType == rules.OutTileSourceTypeSelfDrawnReplacement {
		return []*rules.Pattern{rules.NewPattern("嶺上開花", 10)}
	}
	if context.OutTileSource.SourceType == rules.OutTileSourceTypeAdditionalKong {
		return []*rules.Pattern{rules.NewPattern("搶槓", 10)}
	}
	return nil
}

// 9.4.1 Blessing of Heaven (天和) : 155
// 9.4.2 Blessing of Earth (地和) : 155
func winOnInitialRound(plan rules.OutPlan, context *rules.OutPlanScoringContext) []*rules.Pattern {
	if context.OutTileSource.SourceType == rules.OutTileSourceTypeInitialHand {
		return []*rules.Pattern{rules.NewPattern("天和", 155)}
	}
	if context.OutTileSource.SourceType == rules.OutTileSourceTypeDiscard {
		discardPlayer := context.OutTileSource.DiscardInfo.DiscardPlayer
		if discardPlayer.GetWindOrdinal() == 0 && len(discardPlayer.GetDiscardedTiles()) == 0 {
			return []*rules.Pattern{rules.NewPattern("地和", 155)}
		}
	}
	return nil
}
