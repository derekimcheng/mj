package rules

import (
	"fmt"
	"github.com/derekimcheng/mj/domain"
	"github.com/golang/glog"
)

// OutTileSource represents a tile that caused an Out to be declared.
type OutTileSource struct {
	SourceType OutTileSourceType
	// Tile is not set if SourceType == OutTileSourceTypeInitialHand.
	Tile *domain.Tile
	// DiscardInfo is only set when SourceType == OutTileSourceTypeDiscard.
	DiscardInfo *DiscardInfo
}

// NewOutTileSource creates a new OutTileSource.
func NewOutTileSource(sourceType OutTileSourceType, tile *domain.Tile,
	discardInfo *DiscardInfo) *OutTileSource {
	if (sourceType == OutTileSourceTypeInitialHand) != (tile == nil) {
		panic(fmt.Errorf("sourceType cannot be %s while tile == nil is %t",
			sourceType, tile == nil))
	}
	if (sourceType == OutTileSourceTypeDiscard) != (discardInfo != nil) {
		panic(fmt.Errorf("sourceType cannot be %s while discardInfo != nil is %t",
			sourceType, discardInfo != nil))
	}
	return &OutTileSource{
		SourceType:  sourceType,
		Tile:        tile,
		DiscardInfo: discardInfo,
	}
}

// String ...
func (s *OutTileSource) String() string {
	if s.SourceType == OutTileSourceTypeInitialHand {
		return s.SourceType.String()
	}
	return fmt.Sprintf("%s with %s", s.Tile, s.SourceType)
}

// OutTileSourceType enumerates the different kinds of outside tile sources.
type OutTileSourceType int

const (
	// OutTileSourceTypeDiscard - the tile came from another player's discard.
	OutTileSourceTypeDiscard OutTileSourceType = iota
	// OutTileSourceTypeSelfDrawn - the tile is self-drawn normally, from the front of deck.
	OutTileSourceTypeSelfDrawn
	// OutTileSourceTypeSelfDrawnReplacement - the tile is self-drawn as a result of replacement,
	// from the back of deck.
	OutTileSourceTypeSelfDrawnReplacement
	// OutTileSourceTypeAdditionalKong - the tile came from another player declaring
	// an Additional Kong. Only used in rules that allows "Robbing the Kong".
	OutTileSourceTypeAdditionalKong
	// OutTileSourceTypeInitialHand - out is declared from the initial hand.
	OutTileSourceTypeInitialHand
)

// IsExternalSource returns whether the out tile came from another player.
func (s *OutTileSource) IsExternalSource() bool {
	return s.SourceType == OutTileSourceTypeDiscard ||
		s.SourceType == OutTileSourceTypeAdditionalKong
}

// String ...
func (t OutTileSourceType) String() string {
	switch t {
	case OutTileSourceTypeDiscard:
		return "Discard"
	case OutTileSourceTypeSelfDrawn:
		return "Self-drawn"
	case OutTileSourceTypeSelfDrawnReplacement:
		return "Self-drawn replacement"
	case OutTileSourceTypeAdditionalKong:
		return "Robbed Kong"
	case OutTileSourceTypeInitialHand:
		return "Initial hand"
	}
	glog.Errorf("Unhandled OutTileSourceType %d\n", t)
	return "?"
}

// DiscardInfo stores information for a Out by discard.
type DiscardInfo struct {
	DiscardPlayer *PlayerGameState
}

// NewDiscardInfo creates a new DiscardInfo.
func NewDiscardInfo(discardPlayer *PlayerGameState) *DiscardInfo {
	return &DiscardInfo{
		DiscardPlayer: discardPlayer,
	}
}
