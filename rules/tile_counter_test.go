package rules

import (
	"fmt"
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
// 	return areOutTileGroupSlicesEquivalent(plan1.GetHandGroups(), plan2.GetHandGroups())
// }

// func areOutTileGroupSlicesEquivalent(groups1 []OutTileGroup, groups2 []OutTileGroup) bool {
// 	if len(groups1) != len(groups2) {
// 		return false
// 	}

// 	for _, group1 := range groups1 {
// 		for i, group2 := range groups2 {
// 			if areOutTileGroupsEquivalent(group1, group2) {
// 				groups2 = append(groups2[:i], groups2[i+1:]...)
// 				break
// 			}
// 			return false
// 		}
// 	}
// 	return true
// }

func Test_ComputeOutPlans_AllPongs(t *testing.T) {
	tiles := []*domain.Tile{
		domain.CreateTileForTest(t, dots, 0),
		domain.CreateTileForTest(t, dots, 0),
		domain.CreateTileForTest(t, dots, 0),
		domain.CreateTileForTest(t, dots, 2),
		domain.CreateTileForTest(t, dots, 2),
		domain.CreateTileForTest(t, dots, 2),
		domain.CreateTileForTest(t, dots, 4),
		domain.CreateTileForTest(t, dots, 4),
		domain.CreateTileForTest(t, dots, 4),
		domain.CreateTileForTest(t, dots, 6),
		domain.CreateTileForTest(t, dots, 6),
		domain.CreateTileForTest(t, dots, 6),
		domain.CreateTileForTest(t, dots, 8),
		domain.CreateTileForTest(t, dots, 8),
	}
	hand := domain.NewHand()
	hand.SetTiles(tiles)

	counter := NewHandTileCounter(GetSuitsForGame(), hand, nil)
	plans := counter.ComputeOutPlans()
	// Note: the tiles created below have the same data layout as the tiles above, but they are
	// distinct objects. The test assertions below rely on deep comparison for this test to pass.
	expected := OutPlans{
		NewOutPlan(OutTileGroups{
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 0),
			}, OutTileGroupTypePong),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 2),
				domain.CreateTileForTest(t, dots, 2),
				domain.CreateTileForTest(t, dots, 2),
			}, OutTileGroupTypePong),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 4),
				domain.CreateTileForTest(t, dots, 4),
				domain.CreateTileForTest(t, dots, 4),
			}, OutTileGroupTypePong),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 6),
				domain.CreateTileForTest(t, dots, 6),
				domain.CreateTileForTest(t, dots, 6),
			}, OutTileGroupTypePong),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 8),
				domain.CreateTileForTest(t, dots, 8),
			}, OutTileGroupTypePair),
		}),
	}

	assert.Equal(t, expected, plans)
}

func Test_ComputeOutPlans_PongsOrChows(t *testing.T) {
	tiles := []*domain.Tile{
		domain.CreateTileForTest(t, dots, 0),
		domain.CreateTileForTest(t, dots, 0),
		domain.CreateTileForTest(t, dots, 0),
		domain.CreateTileForTest(t, dots, 1),
		domain.CreateTileForTest(t, dots, 1),
		domain.CreateTileForTest(t, dots, 1),
		domain.CreateTileForTest(t, dots, 2),
		domain.CreateTileForTest(t, dots, 2),
		domain.CreateTileForTest(t, dots, 2),
		domain.CreateTileForTest(t, dots, 3),
		domain.CreateTileForTest(t, dots, 3),
		domain.CreateTileForTest(t, dots, 3),
		domain.CreateTileForTest(t, dots, 4),
		domain.CreateTileForTest(t, dots, 4),
	}
	hand := domain.NewHand()
	hand.SetTiles(tiles)

	counter := NewHandTileCounter(GetSuitsForGame(), hand, nil)
	plans := counter.ComputeOutPlans()

	expected := OutPlans{
		// Plan 1
		NewOutPlan(OutTileGroups{
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 0),
			}, OutTileGroupTypePong),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 1),
			}, OutTileGroupTypePair),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 2),
				domain.CreateTileForTest(t, dots, 3),
			}, OutTileGroupTypeChow),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 2),
				domain.CreateTileForTest(t, dots, 3),
				domain.CreateTileForTest(t, dots, 4),
			}, OutTileGroupTypeChow),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 2),
				domain.CreateTileForTest(t, dots, 3),
				domain.CreateTileForTest(t, dots, 4),
			}, OutTileGroupTypeChow),
		}),
		// Plan 2
		NewOutPlan(OutTileGroups{
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 0),
			}, OutTileGroupTypePong),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 1),
			}, OutTileGroupTypePong),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 2),
				domain.CreateTileForTest(t, dots, 2),
				domain.CreateTileForTest(t, dots, 2),
			}, OutTileGroupTypePong),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 3),
				domain.CreateTileForTest(t, dots, 3),
				domain.CreateTileForTest(t, dots, 3),
			}, OutTileGroupTypePong),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 4),
				domain.CreateTileForTest(t, dots, 4),
			}, OutTileGroupTypePair),
		}),
		// Plan 3
		NewOutPlan(OutTileGroups{
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 0),
			}, OutTileGroupTypePong),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 2),
				domain.CreateTileForTest(t, dots, 3),
			}, OutTileGroupTypeChow),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 2),
				domain.CreateTileForTest(t, dots, 3),
			}, OutTileGroupTypeChow),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 2),
				domain.CreateTileForTest(t, dots, 3),
			}, OutTileGroupTypeChow),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 4),
				domain.CreateTileForTest(t, dots, 4),
			}, OutTileGroupTypePair),
		}),
		// Plan 4
		NewOutPlan(OutTileGroups{
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 2),
			}, OutTileGroupTypeChow),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 2),
			}, OutTileGroupTypeChow),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 2),
			}, OutTileGroupTypeChow),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 3),
				domain.CreateTileForTest(t, dots, 3),
				domain.CreateTileForTest(t, dots, 3),
			}, OutTileGroupTypePong),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 4),
				domain.CreateTileForTest(t, dots, 4),
			}, OutTileGroupTypePair),
		}),
	}

	assert.Equal(t, expected, plans)
}

func Test_ComputeOutPlans_NineGates(t *testing.T) {
	tiles := []*domain.Tile{
		domain.CreateTileForTest(t, dots, 0),
		domain.CreateTileForTest(t, dots, 0),
		domain.CreateTileForTest(t, dots, 0),
		domain.CreateTileForTest(t, dots, 1),
		domain.CreateTileForTest(t, dots, 2),
		domain.CreateTileForTest(t, dots, 3),
		domain.CreateTileForTest(t, dots, 4),
		domain.CreateTileForTest(t, dots, 5),
		domain.CreateTileForTest(t, dots, 6),
		domain.CreateTileForTest(t, dots, 7),
		domain.CreateTileForTest(t, dots, 8),
		domain.CreateTileForTest(t, dots, 8),
		domain.CreateTileForTest(t, dots, 8),
	}

	// The Nine gates hand can go Out with any one of the tile of the same tile.
	for x := 0; x < dots.GetSize(); x++ {
		fmt.Printf("Extra tile is dots %d\n", x+1)
		hand := domain.NewHand()
		hand.SetTiles(tiles)
		outTile := domain.CreateTileForTest(t, dots, x)
		hand.AddTile(outTile)

		counter := NewHandTileCounter(GetSuitsForGame(), hand, nil)
		plans := counter.ComputeOutPlans()
		assert.NotEmpty(t, plans, "Nine gates failed with tile %s", outTile)
	}
}

func Test_ComputeOutPlans_ThreeQuadruples(t *testing.T) {
	tiles := []*domain.Tile{
		domain.CreateTileForTest(t, dots, 0),
		domain.CreateTileForTest(t, dots, 0),
		domain.CreateTileForTest(t, dots, 0),
		domain.CreateTileForTest(t, dots, 0),
		domain.CreateTileForTest(t, dots, 1),
		domain.CreateTileForTest(t, dots, 1),
		domain.CreateTileForTest(t, dots, 1),
		domain.CreateTileForTest(t, dots, 1),
		domain.CreateTileForTest(t, dots, 2),
		domain.CreateTileForTest(t, dots, 2),
		domain.CreateTileForTest(t, dots, 2),
		domain.CreateTileForTest(t, dots, 2),
		domain.CreateTileForTest(t, dots, 3),
		domain.CreateTileForTest(t, dots, 3),
	}
	hand := domain.NewHand()
	hand.SetTiles(tiles)

	counter := NewHandTileCounter(GetSuitsForGame(), hand, nil)
	plans := counter.ComputeOutPlans()

	expected := OutPlans{
		// Plan 1
		NewOutPlan(OutTileGroups{
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 2),
				domain.CreateTileForTest(t, dots, 2),
				domain.CreateTileForTest(t, dots, 2),
				domain.CreateTileForTest(t, dots, 2),
				domain.CreateTileForTest(t, dots, 3),
				domain.CreateTileForTest(t, dots, 3),
			}, OutTileGroupTypeSevenPairs),
		}),
		// Plan 2
		NewOutPlan(OutTileGroups{
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 0),
			}, OutTileGroupTypePair),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 2),
			}, OutTileGroupTypeChow),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 2),
			}, OutTileGroupTypeChow),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 2),
				domain.CreateTileForTest(t, dots, 3),
			}, OutTileGroupTypeChow),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 2),
				domain.CreateTileForTest(t, dots, 3),
			}, OutTileGroupTypeChow),
		}),
		// Plan 3
		NewOutPlan(OutTileGroups{
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 0),
			}, OutTileGroupTypePong),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 2),
			}, OutTileGroupTypeChow),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 1),
			}, OutTileGroupTypePong),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 2),
				domain.CreateTileForTest(t, dots, 2),
				domain.CreateTileForTest(t, dots, 2),
			}, OutTileGroupTypePong),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 3),
				domain.CreateTileForTest(t, dots, 3),
			}, OutTileGroupTypePair),
		}),
		// Plan 4
		NewOutPlan(OutTileGroups{
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 2),
			}, OutTileGroupTypeChow),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 2),
			}, OutTileGroupTypeChow),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 2),
			}, OutTileGroupTypeChow),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 2),
			}, OutTileGroupTypeChow),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 3),
				domain.CreateTileForTest(t, dots, 3),
			}, OutTileGroupTypePair),
		}),
	}

	assert.Equal(t, expected, plans)
}

func Test_ComputeOutPlans_PairAndAChow(t *testing.T) {
	tiles := []*domain.Tile{
		domain.CreateTileForTest(t, dots, 0),
		domain.CreateTileForTest(t, dots, 0),
		domain.CreateTileForTest(t, dots, 0),
		domain.CreateTileForTest(t, dots, 1),
		domain.CreateTileForTest(t, dots, 2),
	}
	hand := domain.NewHand()
	hand.SetTiles(tiles)

	counter := NewHandTileCounter(GetSuitsForGame(), hand, nil)
	plans := counter.ComputeOutPlans()

	expected := OutPlans{
		NewOutPlan(OutTileGroups{
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 0),
			}, OutTileGroupTypePair),
			NewOutTileGroup(domain.Tiles{
				domain.CreateTileForTest(t, dots, 0),
				domain.CreateTileForTest(t, dots, 1),
				domain.CreateTileForTest(t, dots, 2),
			}, OutTileGroupTypeChow),
		}),
	}

	assert.Equal(t, expected, plans)
}
