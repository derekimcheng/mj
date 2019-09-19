package rules

import (
	"errors"
	"fmt"
	"github.com/derekimcheng/mj/domain"
	"github.com/golang/glog"
	"sort"
)

// PlayerGameState represents a player state in a single game.
type PlayerGameState struct {
	hand           *domain.Hand
	bonusTiles     domain.Tiles
	discardedTiles domain.Tiles
	meldGroups     OutTileGroups
}

// NewPlayerGameState creates a blank PlayerGameState object.
func NewPlayerGameState(hand *domain.Hand) *PlayerGameState {
	if hand == nil {
		panic(errors.New("Given hand cannot be nil"))
	}
	return &PlayerGameState{hand: hand, bonusTiles: nil, discardedTiles: nil}
}

// SortHand sorts the tiles in the player's hand.
func (s *PlayerGameState) SortHand() {
	s.hand.Sort()
}

// BulkMoveBonusTilesFromHand removes all bonus tiles from the initial hand and moves them into
// the bonus area. Returns the number of tiles removed. This may be called multiple times before
// the start of the first round.
func (s *PlayerGameState) BulkMoveBonusTilesFromHand() int {
	tiles := RemoveIneligibleTilesFromHand(s.hand)
	if len(tiles) > 0 {
		glog.V(2).Infof("Adding tiles to bonus area: %s\n", tiles)
	}
	s.bonusTiles = append(s.bonusTiles, tiles...)
	return len(tiles)
}

// AddTileToHand adds the given tile to the player's hand. It is an error to call this method with
// a tile that's not eligible for the hand.
func (s *PlayerGameState) AddTileToHand(t *domain.Tile) {
	if !IsEligibleForHand(t.GetSuit()) {
		panic(fmt.Errorf("Attempting to add ineligible tile %s to hand", t))
	}
	s.hand.AddTile(t)
}

// AddTileToHandNoCheck adds the given tile to the player's hand, but does not check whether the
// tile is eligible for hand. This function is appropriate to call when bulk replacing bonus tiles
// at the beginning of a game.
func (s *PlayerGameState) AddTileToHandNoCheck(t *domain.Tile) {
	s.hand.AddTile(t)
}

// AddTileToBonusArea adds the given tile to the bonus area. It is an error to call this method
// with a non-bonus tile.
func (s *PlayerGameState) AddTileToBonusArea(t *domain.Tile) {
	if t.GetSuit().GetSuitType() != domain.SuitTypeBonus {
		panic(fmt.Errorf("Attempting to add non-bonus tile %s to bonus area", t))
	}
	s.bonusTiles = append(s.bonusTiles, t)
	sort.Sort(s.bonusTiles)
}

// DiscardTileAt discards the tile at the given index of the player's hand, and moves it into
// the discard area. Returns whether if the operation was successful, and if so, also returns the
// removed tile.
func (s *PlayerGameState) DiscardTileAt(index int) (*domain.Tile, bool) {
	t, err := s.hand.RemoveTile(index)
	if err != nil {
		glog.V(2).Infof("Failed to remove tile at %d: %s\n", index, err)
		return nil, false
	}
	s.discardedTiles = append(s.discardedTiles, t)
	return t, true
}

// DeclarePong declares a pong using the given tile, which is being discarded. The tiles
// are moved to the meld area. Returns whether the operation was successful.
func (s *PlayerGameState) DeclarePong(t *domain.Tile) bool {
	if !CanPong(t.GetSuit()) {
		glog.V(2).Infof("Tile %s cannot be used in a pong\n", t)
		return false
	}
	count := s.countSimilarTilesInHand(t)
	if count < 2 {
		glog.V(2).Infof("Hand only has %d of %s which is insufficient for pong\n", count, t)
		return false
	}

	// Remove the tiles from hand. Form a meld group with the given tile.
	pongTiles := s.removeSimilarTilesFromHand(t, 2)
	pongTiles = append(pongTiles, t)

	// Add the tiles to a meld group.
	s.meldGroups = append(s.meldGroups, NewOutTileGroup(pongTiles, OutTileGroupTypePong))
	return true
}

// DeclareKong declares a kong using the given tile, which is being discarded. The tiles
// are moved to the meld area. Returns whether the operation was successful.
func (s *PlayerGameState) DeclareKong(t *domain.Tile) bool {
	if !CanPong(t.GetSuit()) {
		glog.V(2).Infof("Tile %s cannot be used in a kong\n", t)
		return false
	}
	count := s.countSimilarTilesInHand(t)
	if count < 3 {
		glog.V(2).Infof("Hand only has %d of %s which is insufficient for kong\n", count, t)
		return false
	}

	// Remove the tiles from hand. Form a meld group with the given tile.
	kongTiles := s.removeSimilarTilesFromHand(t, 3)
	kongTiles = append(kongTiles, t)

	// Add the tiles to a meld group.
	s.meldGroups = append(s.meldGroups, NewOutTileGroup(kongTiles, OutTileGroupTypeKong))
	return true
}

// DeclareConcealedKong declares concealed kong using the tile at the given index. The tiles are
// moved to the meld area. Returns whether if the operation was successful, and if so, also returns
// one of the moved tiles.
func (s *PlayerGameState) DeclareConcealedKong(index int) (*domain.Tile, bool) {
	theTile, err := s.hand.GetTileAt(index)
	if err != nil {
		glog.V(2).Infof("Failed to get tile at %d: %s\n", index, err)
		return nil, false
	}
	if !CanPong(theTile.GetSuit()) {
		glog.V(2).Infof("Tile %s cannot be used in a kong\n", theTile)
		return nil, false
	}

	count := s.countSimilarTilesInHand(theTile)
	if count < 4 {
		glog.V(2).Infof("Hand only has %d of %s which is insufficient for kong\n", count, theTile)
		return nil, false
	}

	// Remove the tiles from hand.
	kongTiles := s.removeSimilarTilesFromHand(theTile, 4)

	// Add the tiles to a meld group.
	s.meldGroups = append(s.meldGroups, NewOutTileGroup(kongTiles, OutTileGroupTypeConcealedKong))
	return kongTiles[0], true
}

// DeclareAdditionalKong declares an additional kong using the tile at the given index. The tile is
// moved to the meld area. Returns whether if the operation was successful, and if so, also returns
// the moved tile.
func (s *PlayerGameState) DeclareAdditionalKong(index int) (*domain.Tile, bool) {
	theTile, err := s.hand.GetTileAt(index)
	if err != nil {
		glog.V(2).Infof("Failed to get tile at %d: %s\n", index, err)
		return nil, false
	}
	if !CanPong(theTile.GetSuit()) {
		glog.V(2).Infof("Tile %s cannot be used in a kong\n", theTile)
		return nil, false
	}

	for _, group := range s.meldGroups {
		if group.GetGroupType() == OutTileGroupTypePong &&
			domain.CompareTiles(group.GetTiles()[0], theTile) == 0 {
			group.UpgradeToKong(theTile)
			theTile, err := s.hand.RemoveTile(index)
			if err != nil {
				panic(err)
			}
			return theTile, true
		}
	}

	glog.V(2).Infof("Failed to find existing Pong group for Tile %s\n", theTile)
	return nil, false
}

// DeclareChow declares chow using the given tile, which is being discarded, and two additional
// tiles at the given indices. The tiles are moved to the meld area. Returns whether if the
// operation was successful, and if so, also returns the tiles used in the chow.
func (s *PlayerGameState) DeclareChow(t *domain.Tile, index1, index2 int) (domain.Tiles, bool) {
	if !CanChow(t.GetSuit()) {
		glog.V(2).Infof("Tile %s cannot be used in a chow\n", t)
		return nil, false
	}
	if index1 == index2 {
		glog.V(2).Infof("The same index is provided twice: %d\n", index1)
		return nil, false
	}

	// Swap indices for convenience.
	if index1 > index2 {
		index1, index2 = index2, index1
	}

	tile1, err := s.hand.GetTileAt(index1)
	if err != nil {
		glog.V(2).Infof("Failed to get tile at %d: %s\n", index1, err)
		return nil, false
	}
	tile2, err := s.hand.GetTileAt(index2)
	if err != nil {
		glog.V(2).Infof("Failed to get tile at %d: %s\n", index2, err)
		return nil, false
	}

	chowTiles := domain.Tiles{tile1, tile2, t}
	sort.Sort(chowTiles)
	if chowTiles[0].GetSuit() != chowTiles[1].GetSuit() ||
		chowTiles[0].GetSuit() != chowTiles[2].GetSuit() {
		glog.V(2).Infof("Not a chow - tiles are not of same suit: %s\n", chowTiles)
		return nil, false
	}
	if chowTiles[1].GetOrdinal() != chowTiles[0].GetOrdinal()+1 ||
		chowTiles[2].GetOrdinal() != chowTiles[1].GetOrdinal()+1 {
		glog.V(2).Infof("Not a chow - tiles are not consecutive: %s\n", chowTiles)
		return nil, false
	}

	// Remove the tiles from hand.
	tiles := s.hand.GetTiles()
	updatedTiles := append(tiles[:index1], tiles[index1+1:index2]...)
	updatedTiles = append(updatedTiles, tiles[index2+1:]...)
	s.hand.SetTiles(updatedTiles)

	// Add the tiles to a meld group.
	s.meldGroups = append(s.meldGroups, NewOutTileGroup(chowTiles, OutTileGroupTypeChow))
	return chowTiles, true
}

// GetHand ...
func (s *PlayerGameState) GetHand() *domain.Hand {
	return s.hand
}

// GetBonusTiles ...
func (s *PlayerGameState) GetBonusTiles() domain.Tiles {
	return s.bonusTiles
}

// GetDiscardedTiles ...
func (s *PlayerGameState) GetDiscardedTiles() domain.Tiles {
	return s.discardedTiles
}

// GetMeldGroups ...
func (s *PlayerGameState) GetMeldGroups() OutTileGroups {
	return s.meldGroups
}

func (s *PlayerGameState) countSimilarTilesInHand(t *domain.Tile) int {
	count := 0
	for _, tile := range s.hand.GetTiles() {
		if tile.GetSuit() == t.GetSuit() && tile.GetOrdinal() == t.GetOrdinal() {
			count++
		}
	}
	return count
}

// removeSimilarTilesFromHand removes |count| tiles of the same suit+ordinal as |t| from the hand
// and returns them. It is invalid to call this function if the hand does not have at least |count|
// tiles similar to |t|.
func (s *PlayerGameState) removeSimilarTilesFromHand(t *domain.Tile, count int) domain.Tiles {
	// Remove the tiles from hand.
	index1 := 0
	tiles := s.hand.GetTiles()
	removedTiles := domain.Tiles{}
	remaining := count
	for index2, tile := range tiles {
		if remaining == 0 || domain.CompareTiles(tile, t) != 0 {
			if index1 != index2 {
				tiles[index1] = tiles[index2]
			}
			index1++
			continue
		}

		removedTiles = append(removedTiles, tile)
		remaining--
	}
	tiles = tiles[:len(tiles)-count]

	if len(removedTiles) != count {
		panic(fmt.Errorf("Expected to have remove %d tiles after loop, got %d",
			count, len(removedTiles)))
	}

	s.hand.SetTiles(tiles)
	return removedTiles
}
