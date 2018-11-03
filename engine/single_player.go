package engine

import (
	"fmt"
	"github.com/derekimcheng/mj/domain"
	"github.com/derekimcheng/mj/rules"
	"github.com/derekimcheng/mj/ui"
	"github.com/pkg/errors"
)

// TODO: put this in rules/?
type gameOverError struct {
	outDeclared bool
}

func newGameOverError(outDeclared bool) *gameOverError {
	return &gameOverError{outDeclared: outDeclared}
}

// Error ... (error implementation)
func (e *gameOverError) Error() string {
	return fmt.Sprintf("Game over. Declared out? %t", e.outDeclared)
}

const (
	separator = "==============================================================="
)

// SinglePlayerRunner is the runner for the single player mode game.
// The single player mode proceeds as follows:
// (1) The player is dealt tiles from the front of a shuffled deck.
// (2) For each bonus tile that was dealt, run step (D2), until there are no more bonus tiles in
//     the "hand".
// (3) The following consists of a single round: Drawing -> Player action -> Discard.
// Drawing:
// (D1) The player draws a tile from the front of the deck.
// (D2) If the tile is a bonus tile, it is moved to the bonus area. Go to step (D2)'. Otherwise,
//      go to step (D3).
// (D2)' Draw a replacement tile from the back of the deck. Go to step (D2).
// (D3) If the drawn tile results in a winning hand, the player may declare Out. The game is over.
// Player action:
// (P1) The player may choose to do the following:
//        Declare a (concealed) Kong -> the kong tiles are moved to the meld area. Run through step
//        (D2)', then repeat (P1).
//        Declare a Pong or a Chow -> the pong or chow tiles are moved to the meld area.
//        Declare Out -> the game is over.
// Discard:
// (R1) The player discards a tile from their hand. That tile is moved to the discard area.
//
// The game ends if any of the conditions are met:
// - The player declares an Out.
// - The deck becomes empty AND a tile is required to be drawn.
type SinglePlayerRunner struct {
	receiver ui.CommandReceiver

	started bool
	// TODO: maybe this would be better encapsulated in a single struct?
	deck domain.Deck
	hand *domain.Hand

	bonusTiles     []*domain.Tile
	discardedTiles []*domain.Tile
}

// NewSinglePlayerRunner returns a new instance of NewSinglePlayerRunner with the given input
// parameters.
func NewSinglePlayerRunner(receiver ui.CommandReceiver) *SinglePlayerRunner {
	return &SinglePlayerRunner{receiver: receiver}
}

// Start starts the game sequence. Returns an error if the game is already started (or ended), or
// if the game is unable to start. This function returns when the game ends.
func (r *SinglePlayerRunner) Start() error {
	if r.started {
		return fmt.Errorf("Already started")
	}

	fmt.Printf("Starting single player game\n")
	r.started = true

	r.initializeDeck()
	err := r.initializeHandNoReplacement()
	if err != nil {
		return errors.Wrapf(err, "unable to start game")
	}

	// From here on, the game logic could throw exception to signal that the game is over.
	r.startGameSequence()
	return nil
}

func (r *SinglePlayerRunner) initializeDeck() {
	fmt.Printf("Initializing deck\n")
	r.deck = rules.NewDeckForGame()
	r.deck.Shuffle()
}

func (r *SinglePlayerRunner) initializeHandNoReplacement() error {
	fmt.Printf("Initializing hand\n")
	r.hand = domain.NewHand()
	err := rules.PopulateHands(r.deck, []*domain.Hand{r.hand})
	if err != nil {
		return err
	}

	r.sortHand()
	fmt.Printf("Populated hand: %s\n", r.hand)
	return nil
}

func (r *SinglePlayerRunner) startGameSequence() {
	defer func() {
		if r := recover(); r != nil {
			if gameOverErr, ok := r.(*gameOverError); ok {
				// TODO: implement proper panic recovery.
				fmt.Printf("Game over: %s\n", gameOverErr)
			} else {
				panic(r)
			}
		}
	}()

	fmt.Printf("Starting game sequence\n")

	numTilesToReplace := r.bulkMoveBonusTilesFromHand()
	if numTilesToReplace > 0 {
		replacementRound := 1
		for numTilesToReplace > 0 {
			fmt.Printf("Replacing %d bonus tiles (round %d)\n", numTilesToReplace, replacementRound)
			for i := 0; i < numTilesToReplace; i++ {
				tile := r.drawFromDeckFront()
				r.addTileToHand(tile)
			}
			numTilesToReplace = r.bulkMoveBonusTilesFromHand()
			replacementRound++
		}
	}

	r.sortHand()
	fmt.Printf("Hand after replacement: %s\n", r.hand)

	r.startPlayerRoundLoop()
}

func (r *SinglePlayerRunner) startPlayerRoundLoop() {
	round := 1
	for {
		fmt.Println(separator)
		fmt.Printf("Start of round %d\n", round)

		// Draw phase
		tile := r.drawFromDeckFront()
		fmt.Printf("Drawn tile %s\n", tile)

		if rules.IsEligibleForHand(tile.GetSuit()) {
			r.addTileToHand(tile)
		} else {
			r.replaceTileLoop()
		}

		fmt.Printf("Hand: %s\n", r.hand)
		// Player action phase
		r.promptAndExecutePlayerAction(
			ui.CommandTypes{ui.SortHand, ui.ShowDiscardedTiles, ui.DiscardTile, ui.Out})

		round++
	}
}

func (r *SinglePlayerRunner) promptAndExecutePlayerAction(acceptedCommands ui.CommandTypes) {
	for {
		cmd, err := r.receiver.PromptForCommand(acceptedCommands)
		if err != nil {
			// TODO: this should be its own error struct. Something like IOError.
			panic(newGameOverError(false))
		}

		proceedToNextRound := r.executePlayerAction(cmd)
		if proceedToNextRound {
			break
		}
	}
}

func (r *SinglePlayerRunner) executePlayerAction(cmd *ui.Command) bool {
	switch cmd.GetCommandType() {
	case ui.SortHand:
		r.sortAndShowHand()
		return false
	case ui.ShowDiscardedTiles:
		r.showDiscardedTiles()
		return false
	case ui.DiscardTile:
		return r.discardTile(cmd.GetDiscardTileCommand().GetIndex())
	case ui.Out:
		return r.checkForOut()
	}
	fmt.Printf("Unhandled command: %s\n", cmd.GetCommandType())
	return false
}

// TODO: probably a PlayerState method
func (r *SinglePlayerRunner) bulkMoveBonusTilesFromHand() int {
	tiles := rules.RemoveIneligibleTilesFromHand(r.hand)
	if len(tiles) > 0 {
		fmt.Printf("Adding tiles to bonus area: %s\n", tiles)
	}
	r.bonusTiles = append(r.bonusTiles, tiles...)
	return len(tiles)
}

func (r *SinglePlayerRunner) addTileToBonusArea(t *domain.Tile) {
	// TODO: how to sort r.bonusTiles?
	fmt.Printf("Adding tile to bonus area: %s\n", t)
	r.bonusTiles = append(r.bonusTiles, t)
}

func (r *SinglePlayerRunner) drawFromDeckFront() *domain.Tile {
	tile, err := r.deck.PopFront()
	if err != nil {
		panic(newGameOverError(false))
	}
	return tile
}

func (r *SinglePlayerRunner) replaceTileLoop() {
	for round := 1; ; /* no-op */ round++ {
		fmt.Printf("Drawing a replacement tile from the back of deck (round %d)\n", round)
		tile, err := r.deck.PopBack()
		if err != nil {
			panic(newGameOverError(false))
		}
		if rules.IsEligibleForHand(tile.GetSuit()) {
			fmt.Printf("Adding replacement tile to hand: %s\n", tile)
			r.addTileToHand(tile)
			return
		}
		// Else tile is a bonus tile, add it and repeat.
		r.addTileToBonusArea(tile)
	}
}

func (r *SinglePlayerRunner) addTileToHand(t *domain.Tile) {
	r.hand.AddTile(t)
}

func (r *SinglePlayerRunner) discardTile(index int) bool {
	t, err := r.hand.RemoveTile(index)
	if err != nil {
		fmt.Printf("Failed to remove tile at %d: %s\n", index, err)
		return false
	}
	fmt.Printf("Discarded tile at %d: %s\n", index, t)
	r.discardedTiles = append(r.discardedTiles, t)
	return true
}

func (r *SinglePlayerRunner) sortHand() {
	r.hand.Sort()
}

func (r *SinglePlayerRunner) sortAndShowHand() {
	r.sortHand()
	r.showHand()
}

func (r *SinglePlayerRunner) showHand() {
	fmt.Printf("Hand: %s\n", r.hand)
}

func (r *SinglePlayerRunner) showDiscardedTiles() {
	// TODO: fix printing
	fmt.Printf("Discarded tiles: %s\n", r.discardedTiles)
}

// checkForOut checks whether the current player state represents an Out. This function will
// panic if the hand is an out hand, or return false if it is not an out hand.
// TODO: score the out.
func (r *SinglePlayerRunner) checkForOut() bool {
	counter := rules.NewHandTileCounter(rules.GetSuitsForGame(), r.hand)
	plans := counter.ComputeOutPlans()

	if len(plans) > 0 {
		panic(newGameOverError(true))
	}
	fmt.Println("Not an Out hand!")
	return false
}
