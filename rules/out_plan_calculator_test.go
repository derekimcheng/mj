package rules

import (
	"github.com/derekimcheng/mj/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

// func areOutPlansResultsEquivalent(plans1 []OutPlan, plans2 []OutPlan) bool {
// 	if len(plans1) != len(plans2) {
// 		return false
// 	}
// 	for _, plan1 := range plans1 {
// 		for i, plan2 := range plans2 {
// 			if (areOutPlansEquivalent(plan1, plan2)) {
// 				plans2 = append(plans2[:i], plans2[i+1:]...)
// 				break
// 			}
// 			return false
// 		}
// 	}
// 	return true
// }

// func areOutPlansEquivalent(plan1 OutPlan, plan2 OutPlan) bool {
// 	return areTileGroupSlicesEquivalent(plan1.GetHandGroups(), plan2.GetHandGroups())
// }

// func areTileGroupSlicesEquivalent(groups1 []TileGroup, groups2 []TileGroup) bool {
// 	if len(groups1) != len(groups2) {
// 		return false
// 	}

// 	for _, group1 := range groups1 {
// 		for i, group2 := range groups2 {
// 			if areTileGroupsEquivalent(group1, group2) {
// 				groups2 = append(groups2[:i], groups2[i+1:]...)
// 				break
// 			}
// 			return false
// 		}
// 	}
// 	return true
// }

func createOutTileSourceForTest(tile *domain.Tile) *OutTileSource {
	return NewOutTileSource(OutTileSourceTypeSelfDrawn, tile, nil)
}

func Test_ComputeOutPlans_AllPongs(t *testing.T) {
	tiles := []*domain.Tile{
		domain.CreateTileForTest(t, Dots, 0),
		domain.CreateTileForTest(t, Dots, 0),
		domain.CreateTileForTest(t, Dots, 0),
		domain.CreateTileForTest(t, Dots, 2),
		domain.CreateTileForTest(t, Dots, 2),
		domain.CreateTileForTest(t, Dots, 2),
		domain.CreateTileForTest(t, Dots, 4),
		domain.CreateTileForTest(t, Dots, 4),
		domain.CreateTileForTest(t, Dots, 4),
		domain.CreateTileForTest(t, Dots, 6),
		domain.CreateTileForTest(t, Dots, 6),
		domain.CreateTileForTest(t, Dots, 6),
		domain.CreateTileForTest(t, Dots, 8),
		domain.CreateTileForTest(t, Dots, 8),
	}
	hand := domain.NewHand()
	hand.SetTiles(tiles)

	player := NewPlayerGameState(hand, 0)
	calculator := NewOutPlanCalculator(GetSuitsForGame(), player,
		createOutTileSourceForTest(hand.GetTiles()[0]))
	plans := calculator.Calculate()
	// Note: the tiles created below have the same data layout as the tiles above, but they are
	// distinct objects. The test assertions below rely on deep comparison for this test to pass.
	handGroups := TileGroups{
		NewTileGroup(domain.Tiles{
			domain.CreateTileForTest(t, Dots, 0),
			domain.CreateTileForTest(t, Dots, 0),
			domain.CreateTileForTest(t, Dots, 0),
		}, TileGroupTypePong),
		NewTileGroup(domain.Tiles{
			domain.CreateTileForTest(t, Dots, 2),
			domain.CreateTileForTest(t, Dots, 2),
			domain.CreateTileForTest(t, Dots, 2),
		}, TileGroupTypePong),
		NewTileGroup(domain.Tiles{
			domain.CreateTileForTest(t, Dots, 4),
			domain.CreateTileForTest(t, Dots, 4),
			domain.CreateTileForTest(t, Dots, 4),
		}, TileGroupTypePong),
		NewTileGroup(domain.Tiles{
			domain.CreateTileForTest(t, Dots, 6),
			domain.CreateTileForTest(t, Dots, 6),
			domain.CreateTileForTest(t, Dots, 6),
		}, TileGroupTypePong),
		NewTileGroup(domain.Tiles{
			domain.CreateTileForTest(t, Dots, 8),
			domain.CreateTileForTest(t, Dots, 8),
		}, TileGroupTypePair)}
	expected := OutPlans{NewOutPlan(handGroups, nil)}

	assert.Equal(t, expected, plans)
}

func Test_ComputeOutPlans_PongsOrChows(t *testing.T) {
	tiles := []*domain.Tile{
		domain.CreateTileForTest(t, Dots, 0),
		domain.CreateTileForTest(t, Dots, 0),
		domain.CreateTileForTest(t, Dots, 0),
		domain.CreateTileForTest(t, Dots, 1),
		domain.CreateTileForTest(t, Dots, 1),
		domain.CreateTileForTest(t, Dots, 1),
		domain.CreateTileForTest(t, Dots, 2),
		domain.CreateTileForTest(t, Dots, 2),
		domain.CreateTileForTest(t, Dots, 2),
		domain.CreateTileForTest(t, Dots, 3),
		domain.CreateTileForTest(t, Dots, 3),
		domain.CreateTileForTest(t, Dots, 3),
		domain.CreateTileForTest(t, Dots, 4),
		domain.CreateTileForTest(t, Dots, 4),
	}
	hand := domain.NewHand()
	hand.SetTiles(tiles)

	player := NewPlayerGameState(hand, 0)
	calculator := NewOutPlanCalculator(GetSuitsForGame(), player,
		createOutTileSourceForTest(hand.GetTiles()[0]))
	plans := calculator.Calculate()

	expected := OutPlans{
		// Plan 1
		NewOutPlan(TileGroups{
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 0),
			}, TileGroupTypePong),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 1),
			}, TileGroupTypePair),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 2),
				domain.CreateTileForTest(t, Dots, 3),
			}, TileGroupTypeChow),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 2),
				domain.CreateTileForTest(t, Dots, 3),
				domain.CreateTileForTest(t, Dots, 4),
			}, TileGroupTypeChow),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 2),
				domain.CreateTileForTest(t, Dots, 3),
				domain.CreateTileForTest(t, Dots, 4),
			}, TileGroupTypeChow),
		}, /*meldedGroups*/ nil),
		// Plan 2
		NewOutPlan(TileGroups{
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 0),
			}, TileGroupTypePong),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 1),
			}, TileGroupTypePong),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 2),
				domain.CreateTileForTest(t, Dots, 2),
				domain.CreateTileForTest(t, Dots, 2),
			}, TileGroupTypePong),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 3),
				domain.CreateTileForTest(t, Dots, 3),
				domain.CreateTileForTest(t, Dots, 3),
			}, TileGroupTypePong),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 4),
				domain.CreateTileForTest(t, Dots, 4),
			}, TileGroupTypePair),
		}, /*meldedGroups*/ nil),
		// Plan 3
		NewOutPlan(TileGroups{
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 0),
			}, TileGroupTypePong),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 2),
				domain.CreateTileForTest(t, Dots, 3),
			}, TileGroupTypeChow),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 2),
				domain.CreateTileForTest(t, Dots, 3),
			}, TileGroupTypeChow),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 2),
				domain.CreateTileForTest(t, Dots, 3),
			}, TileGroupTypeChow),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 4),
				domain.CreateTileForTest(t, Dots, 4),
			}, TileGroupTypePair),
		}, /*meldedGroups*/ nil),
		// Plan 4
		NewOutPlan(TileGroups{
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 2),
			}, TileGroupTypeChow),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 2),
			}, TileGroupTypeChow),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 2),
			}, TileGroupTypeChow),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 3),
				domain.CreateTileForTest(t, Dots, 3),
				domain.CreateTileForTest(t, Dots, 3),
			}, TileGroupTypePong),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 4),
				domain.CreateTileForTest(t, Dots, 4),
			}, TileGroupTypePair),
		}, /*meldedGroups*/ nil),
	}

	assert.Equal(t, expected, plans)
}

func Test_ComputeOutPlans_NineGates(t *testing.T) {
	tiles := []*domain.Tile{
		domain.CreateTileForTest(t, Dots, 0),
		domain.CreateTileForTest(t, Dots, 0),
		domain.CreateTileForTest(t, Dots, 0),
		domain.CreateTileForTest(t, Dots, 1),
		domain.CreateTileForTest(t, Dots, 2),
		domain.CreateTileForTest(t, Dots, 3),
		domain.CreateTileForTest(t, Dots, 4),
		domain.CreateTileForTest(t, Dots, 5),
		domain.CreateTileForTest(t, Dots, 6),
		domain.CreateTileForTest(t, Dots, 7),
		domain.CreateTileForTest(t, Dots, 8),
		domain.CreateTileForTest(t, Dots, 8),
		domain.CreateTileForTest(t, Dots, 8),
	}

	// The Nine gates hand can go Out with any one of the tile of the same tile.
	for x := 0; x < Dots.GetSize(); x++ {
		hand := domain.NewHand()
		hand.SetTiles(tiles)
		outTile := domain.CreateTileForTest(t, Dots, x)
		hand.AddTile(outTile)

		player := NewPlayerGameState(hand, 0)
		calculator := NewOutPlanCalculator(GetSuitsForGame(), player,
			createOutTileSourceForTest(hand.GetTiles()[0]))
		plans := calculator.Calculate()
		assert.NotEmpty(t, plans, "Nine gates failed with tile %s", outTile)
	}
}

func Test_ComputeOutPlans_ThreeQuadruples(t *testing.T) {
	tiles := []*domain.Tile{
		domain.CreateTileForTest(t, Dots, 0),
		domain.CreateTileForTest(t, Dots, 0),
		domain.CreateTileForTest(t, Dots, 0),
		domain.CreateTileForTest(t, Dots, 0),
		domain.CreateTileForTest(t, Dots, 1),
		domain.CreateTileForTest(t, Dots, 1),
		domain.CreateTileForTest(t, Dots, 1),
		domain.CreateTileForTest(t, Dots, 1),
		domain.CreateTileForTest(t, Dots, 2),
		domain.CreateTileForTest(t, Dots, 2),
		domain.CreateTileForTest(t, Dots, 2),
		domain.CreateTileForTest(t, Dots, 2),
		domain.CreateTileForTest(t, Dots, 3),
		domain.CreateTileForTest(t, Dots, 3),
	}
	hand := domain.NewHand()
	hand.SetTiles(tiles)

	player := NewPlayerGameState(hand, 0)
	calculator := NewOutPlanCalculator(GetSuitsForGame(), player,
		createOutTileSourceForTest(hand.GetTiles()[0]))
	plans := calculator.Calculate()

	expected := OutPlans{
		// Plan 1
		NewOutPlan(TileGroups{
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 2),
				domain.CreateTileForTest(t, Dots, 2),
				domain.CreateTileForTest(t, Dots, 2),
				domain.CreateTileForTest(t, Dots, 2),
				domain.CreateTileForTest(t, Dots, 3),
				domain.CreateTileForTest(t, Dots, 3),
			}, TileGroupTypeSevenPairs),
		}, /*meldedGroups*/ nil),
		// Plan 2
		NewOutPlan(TileGroups{
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 0),
			}, TileGroupTypePair),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 2),
			}, TileGroupTypeChow),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 2),
			}, TileGroupTypeChow),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 2),
				domain.CreateTileForTest(t, Dots, 3),
			}, TileGroupTypeChow),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 2),
				domain.CreateTileForTest(t, Dots, 3),
			}, TileGroupTypeChow),
		}, /*meldedGroups*/ nil),
		// Plan 3
		NewOutPlan(TileGroups{
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 0),
			}, TileGroupTypePong),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 2),
			}, TileGroupTypeChow),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 1),
			}, TileGroupTypePong),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 2),
				domain.CreateTileForTest(t, Dots, 2),
				domain.CreateTileForTest(t, Dots, 2),
			}, TileGroupTypePong),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 3),
				domain.CreateTileForTest(t, Dots, 3),
			}, TileGroupTypePair),
		}, /*meldedGroups*/ nil),
		// Plan 4
		NewOutPlan(TileGroups{
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 2),
			}, TileGroupTypeChow),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 2),
			}, TileGroupTypeChow),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 2),
			}, TileGroupTypeChow),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 2),
			}, TileGroupTypeChow),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 3),
				domain.CreateTileForTest(t, Dots, 3),
			}, TileGroupTypePair),
		}, /*meldedGroups*/ nil),
	}

	assert.Equal(t, expected, plans)
}

func Test_ComputeOutPlans_PairAndAChow(t *testing.T) {
	tiles := []*domain.Tile{
		domain.CreateTileForTest(t, Dots, 0),
		domain.CreateTileForTest(t, Dots, 0),
		domain.CreateTileForTest(t, Dots, 0),
		domain.CreateTileForTest(t, Dots, 1),
		domain.CreateTileForTest(t, Dots, 2),
	}
	hand := domain.NewHand()
	hand.SetTiles(tiles)

	player := NewPlayerGameState(hand, 0)
	calculator := NewOutPlanCalculator(GetSuitsForGame(), player,
		createOutTileSourceForTest(hand.GetTiles()[0]))
	plans := calculator.Calculate()

	expected := OutPlans{
		NewOutPlan(TileGroups{
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 0),
			}, TileGroupTypePair),
			NewTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, Dots, 0),
				domain.CreateTileForTest(t, Dots, 1),
				domain.CreateTileForTest(t, Dots, 2),
			}, TileGroupTypeChow),
		}, /*meldedGroups*/ nil),
	}

	assert.Equal(t, expected, plans)
}
