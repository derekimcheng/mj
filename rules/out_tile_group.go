package rules

import (
	"sort"
	"fmt"
	"github.com/derekimcheng/mj/domain"
)

// OutTileGroup represents a group of tiles that make up an Out plan.
type OutTileGroup struct {
	// tiles is the sorted list of tiles that make up the group.
	tiles     domain.Tiles
	groupType OutTileGroupType
}

// UpgradeToKong upgrades the group from Pong to Kong as a result of additional Kong. It is an error
// to call this if the group cannot be upgraded with the given tile.
func (g *OutTileGroup) UpgradeToKong(tile *domain.Tile) {
	if g.groupType != OutTileGroupTypePong || domain.CompareTiles(g.tiles[0], tile) != 0 {
		panic(fmt.Errorf("UpgradeToKong failed for group %s, tile %s", g, tile))
	}
	g.tiles = append(g.tiles, tile)
	g.groupType = OutTileGroupTypeKong
}

// GetTiles ...
func (g *OutTileGroup) GetTiles() domain.Tiles {
	return g.tiles
}

// GetGroupType ...
func (g *OutTileGroup) GetGroupType() OutTileGroupType {
	return g.groupType
}

// String ...
func (g *OutTileGroup) String() string {
	return fmt.Sprintf("%s: %v", g.groupType, g.tiles)
}

// NewOutTileGroup creates a new OutTileGroup with the given parameters. The input tiles is not
// copied, and will be modified by sorting.
func NewOutTileGroup(tiles domain.Tiles, groupType OutTileGroupType) *OutTileGroup {
	sort.Sort(tiles)
	return &OutTileGroup{tiles: tiles, groupType: groupType}
}

// OutTileGroups is a slice of OutTileGroup.
type OutTileGroups []*OutTileGroup

// Len ... (implements sort.Interface)
func (groups OutTileGroups) Len() int {
	return len(groups)
}

// Swap ... (implements sort.Interface)
func (groups OutTileGroups) Swap(i, j int) {
	groups[i], groups[j] = groups[j], groups[i]
}

// Less ... (implements sort.Interface)
func (groups OutTileGroups) Less(i, j int) bool {
	res := domain.CompareTiles(groups[i].GetTiles()[0], groups[j].GetTiles()[0])
	if res != 0 {
		return res < 0
	}
	return groups[i].GetGroupType() < groups[j].GetGroupType()
}