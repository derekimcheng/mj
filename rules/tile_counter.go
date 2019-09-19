package rules

import (
	"github.com/derekimcheng/mj/domain"
)

var (
	thirteenOrphanTiles = []domain.TileBase{
		domain.NewTileBase(bamboo, 0),
		domain.NewTileBase(bamboo, 8),
		domain.NewTileBase(dots, 0),
		domain.NewTileBase(dots, 8),
		domain.NewTileBase(characters, 0),
		domain.NewTileBase(characters, 8),
		domain.NewTileBase(winds, 0),
		domain.NewTileBase(winds, 1),
		domain.NewTileBase(winds, 2),
		domain.NewTileBase(winds, 3),
		domain.NewTileBase(dragons, 0),
		domain.NewTileBase(dragons, 1),
		domain.NewTileBase(dragons, 2),
	}
)

type tileInventory = map[*domain.Suit][][]*domain.Tile
