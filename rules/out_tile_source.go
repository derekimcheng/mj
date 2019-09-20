package rules

import (
	"github.com/golang/glog"
	"fmt"
	"github.com/derekimcheng/mj/domain"
)

// OutTileSource represents a tile that caused an Out to be declared.
type OutTileSource struct {
	Tile       *domain.Tile
	SourceType OutTileSourceType
}

// NewOutTileSource creates a new OutTileSource.
func NewOutTileSource(tile *domain.Tile, sourceType OutTileSourceType) *OutTileSource {
	return &OutTileSource{Tile: tile, SourceType: sourceType}
}

// String ...
func (s *OutTileSource) String() string {
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
	}
	glog.Errorf("Unhandled OutTileSourceType %d\n", t)
	return "?"
}