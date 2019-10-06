package rules

import (
	"github.com/derekimcheng/mj/domain"
)

var (
	thirteenOrphanTiles = []domain.TileBase{
		domain.NewTileBase(Bamboo, 0),
		domain.NewTileBase(Bamboo, 8),
		domain.NewTileBase(Dots, 0),
		domain.NewTileBase(Dots, 8),
		domain.NewTileBase(Characters, 0),
		domain.NewTileBase(Characters, 8),
		domain.NewTileBase(Winds, 0),
		domain.NewTileBase(Winds, 1),
		domain.NewTileBase(Winds, 2),
		domain.NewTileBase(Winds, 3),
		domain.NewTileBase(Dragons, 0),
		domain.NewTileBase(Dragons, 1),
		domain.NewTileBase(Dragons, 2),
	}
)

type tileInventory = map[*domain.Suit][][]*domain.Tile
