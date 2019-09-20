package rules

import (
	"github.com/golang/glog"
	"fmt"
	"github.com/derekimcheng/mj/domain"
	"sort"
)

// TileGroup represents a group of tiles that make up an Out plan.
type TileGroup struct {
	// tiles is the sorted list of tiles that make up the group.
	tiles     domain.Tiles
	groupType TileGroupType
}

// UpgradeToKong upgrades the group from Pong to Kong as a result of additional Kong. It is an error
// to call this if the group cannot be upgraded with the given tile.
func (g *TileGroup) UpgradeToKong(tile *domain.Tile) {
	if g.groupType != TileGroupTypePong || domain.CompareTiles(g.tiles[0], tile) != 0 {
		panic(fmt.Errorf("UpgradeToKong failed for group %s, tile %s", g, tile))
	}
	g.tiles = append(g.tiles, tile)
	g.groupType = TileGroupTypeKong
}

// GetTiles ...
func (g *TileGroup) GetTiles() domain.Tiles {
	return g.tiles
}

// GetGroupType ...
func (g *TileGroup) GetGroupType() TileGroupType {
	return g.groupType
}

// String ...
func (g *TileGroup) String() string {
	return fmt.Sprintf("%s: %v", g.groupType, g.tiles)
}

// NewTileGroup creates a new TileGroup with the given parameters. The input tiles is not
// copied, and will be modified by sorting.
func NewTileGroup(tiles domain.Tiles, groupType TileGroupType) *TileGroup {
	sort.Sort(tiles)
	return &TileGroup{tiles: tiles, groupType: groupType}
}

// TileGroups is a slice of TileGroup.
type TileGroups []*TileGroup

// Len ... (implements sort.Interface)
func (groups TileGroups) Len() int {
	return len(groups)
}

// Swap ... (implements sort.Interface)
func (groups TileGroups) Swap(i, j int) {
	groups[i], groups[j] = groups[j], groups[i]
}

// Less ... (implements sort.Interface)
func (groups TileGroups) Less(i, j int) bool {
	res := domain.CompareTiles(groups[i].GetTiles()[0], groups[j].GetTiles()[0])
	if res != 0 {
		return res < 0
	}
	return groups[i].GetGroupType() < groups[j].GetGroupType()
}

// TileGroupType represents a type of Out tile group
type TileGroupType int

const (
	// TileGroupTypePair represents a pair of same suit+ordinal tiles. With the exception of
	// special hands, each out plan must contain exactly one Pair group.
	TileGroupTypePair TileGroupType = iota
	// TileGroupTypePong represents a triplet of the same suit+ordinal tiles.
	TileGroupTypePong
	// TileGroupTypeChow represents a triplet of consecutive tiles.
	TileGroupTypeChow
	// TileGroupTypeKong represents a quadruplet of the same suit+ordinal tiles. This is only
	// available in the meld area.
	TileGroupTypeKong
	// TileGroupTypeConcealedKong represents a concealed quadruplet of the same suit+ordinal
	// tiles. This is only available in the meld area.
	TileGroupTypeConcealedKong
	// TileGroupTypeSevenPairs is a special designation for "Seven Pairs". All of the tiles will
	// be represented as a single group.
	TileGroupTypeSevenPairs
	// TileGroupTypeThirteenOrphans is a special designation for "Thirteen Orphans". All of
	// the tiles will be represented as a single group.
	TileGroupTypeThirteenOrphans
)

func (t TileGroupType) String() string {
	switch t {
	case TileGroupTypePair:
		return "Pair"
	case TileGroupTypeChow:
		return "Chow"
	case TileGroupTypePong:
		return "Pong"
	case TileGroupTypeKong:
		return "Kong"
	case TileGroupTypeConcealedKong:
		return "ConcealedKong"
	case TileGroupTypeSevenPairs:
		return "SevenPairs"
	case TileGroupTypeThirteenOrphans:
		return "ThirteenOrphans"
	}
	glog.Errorf("Unhandled TileGroupType %d\n", t)
	return "?"
}