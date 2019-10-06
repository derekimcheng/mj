package shorthand

import (
	"fmt"
	"github.com/derekimcheng/mj/domain"
	"github.com/derekimcheng/mj/rules"
	"github.com/pkg/errors"
)

const (
	// Letters used in shorthand notation.
	character = 'm'
	dot       = 'd'
	bamboo    = 'b'
	winds     = 'w'
	dragons   = 'y'
	flowers   = 'f'
	seasons   = 's'
)

var lettersToSuits = map[rune]*domain.Suit{
	character: rules.Characters,
	dot:       rules.Dots,
	bamboo:    rules.Bamboo,
	winds:     rules.Winds,
	dragons:   rules.Dragons,
	seasons:   rules.Seasons,
}

// Parser turns shorthand form strings into tiles / meld groups.
type Parser struct {
	nextID int
}

// NewParser creates a new Parser.
func NewParser() *Parser {
	return &Parser{nextID: 0}
}

// ParseTiles parses the given tiles shorthand form and returns the corresponding Tiles, or an error
// if the input is invalid.
func (r *Parser) ParseTiles(tilesStr string) (domain.Tiles, error) {
	var tiles domain.Tiles
	var ordinalsSoFar []int
	for _, c := range tilesStr {
		if c >= '1' && c <= '9' {
			ordinalsSoFar = append(ordinalsSoFar, int(c)-int('1'))
			continue
		}
		suit, found := lettersToSuits[c]
		if !found {
			return nil, fmt.Errorf("Unknown suit %c", c)
		}
		if len(ordinalsSoFar) == 0 {
			return nil, fmt.Errorf("Suit must be preceded with at least one number")
		}
		for _, ordinal := range ordinalsSoFar {
			if ordinal >= suit.GetSize() {
				return nil, fmt.Errorf("%d is out of range for suit %s",
					ordinal+1, suit.GetName())
			}
			tile, err := domain.NewTile(suit, ordinal, r.nextID)
			if err != nil {
				return nil, errors.Wrapf(err, "error constructing tile")
			}
			tiles = append(tiles, tile)
			r.nextID++
		}
		ordinalsSoFar = nil
	}
	return tiles, nil
}

// ParseMeldGroups parses the given tiles shorthand form and returns the corresponding meld groups,
// or an error if the input is invalid.
func (r *Parser) ParseMeldGroups(meldGroupsStr string) (rules.TileGroups, error) {
	// TODO: implement
	return nil, nil
}
