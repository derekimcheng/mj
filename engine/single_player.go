package engine

import (
	"fmt"
	"github.com/derekimcheng/mj/domain"
	"github.com/derekimcheng/mj/flags"
	"github.com/derekimcheng/mj/rules"
	"github.com/derekimcheng/mj/rules/zj"
	"github.com/derekimcheng/mj/ui"
	"github.com/golang/glog"
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

var (
	commonCommands           = ui.CommandTypes{ui.SortHand, ui.ShowDiscardedTiles, ui.ShowMelded}
	commandsAfterDrawingTile = withCommands(ui.DiscardTile, ui.ConcealedKong, ui.AdditionalKong, ui.Out)
)

func withCommands(types ...ui.CommandType) ui.CommandTypes {
	return append(commonCommands, types...)
}

// SinglePlayerRunner is the runner for the single player mode game.
// The single player mode proceeds as follows:
// (1) The player is dealt tiles from the front of a shuffled deck.
// (2) For each bonus tile that was dealt, run step (D2), until there are no more bonus tiles in
//     the "hand".
// (3) The following consists of a single round: Drawing -> Player action -> Discard -> Burn. The
//     Burn phase may be executed multiple times (3 times would simulate a 4-player game), the last
//     of which allows the player to chow the tile being discarded. The exception is the first round
//     where Drawing is not performed.
// Drawing:
// (D1) The player draws a tile from the front of the deck.
// (D2) If the tile is a bonus tile, it is moved to the bonus area. Go to step (D2)'. Otherwise,
//      go to step (D3).
// (D2)' Draw a replacement tile from the back of the deck. Go to step (D2).
// (D3) If the drawn tile results in a winning hand, the player may declare Out. The game is over.
// Player action:
// (P1) The player may choose to do the following:
//        - Declare a concealed or additional Kong -> the kong tiles are moved to the meld area. Run
//        through step (D2)', then repeat (P1).
//        - Declare Out -> the game is over.
// Discard:
// (R1) The player discards a tile from their hand. That tile is moved to the discard area.
// Burn:
// (B1) Execute the Drawing algorithm to obtain a non-bonus tile.
// (B2) The player may choose to form a meld group by declaring a pong or kong with the tile. If
//      this Burn is the last of the current round, the player may also choose to declare a chow.
// (B2)' The player may also choose to declare an Out, in which case the game is over.
// (B3) If the player chooses to form a meld group, they must discard a tile from the hand.
//      If the player chooses to kong, execute step (D2)' before discarding.
//      Otherwise, the tile drawn in (B1) is discarded.
// The game ends if any of the conditions are met:
// - The player declares an Out.
// - The deck becomes empty AND a tile is required to be drawn.
type SinglePlayerRunner struct {
	receiver         ui.CommandReceiver
	numBurnsPerRound int

	started bool
	deck    domain.Deck
	player  *rules.PlayerGameState
	// pseudoOpponentGameState is used for populating DiscardInfo. It does not contain a real hand.
	// However, it is used to store "burn" tiles and bonus tiles drawn outside the player's turn.
	// Its wind ordinal is the one that proceeds from that of the player.
	pseudoOpponentGameState *rules.PlayerGameState

	// Tile drawn in the current Burn stage which can be used for a meld.
	currentBurnTile *domain.Tile
}

// NewSinglePlayerRunner returns a new instance of NewSinglePlayerRunner with the given input
// parameters.
func NewSinglePlayerRunner(receiver ui.CommandReceiver) *SinglePlayerRunner {
	if *flags.NumBurnsFlag < 0 || *flags.NumBurnsFlag > 3 {
		panic(fmt.Errorf("Invalid value for numBurnsFlag: %d", *flags.NumBurnsFlag))
	}
	return &SinglePlayerRunner{receiver: receiver, numBurnsPerRound: *flags.NumBurnsFlag}
}

// Start starts the game sequence. Returns an error if the game is already started (or ended), or
// if the game is unable to start. This function returns when the game ends.
func (r *SinglePlayerRunner) Start(deck domain.Deck) error {
	if r.started {
		return errors.New("Already started")
	}
	if deck.IsEmpty() {
		return errors.New("Deck is empty")
	}

	glog.V(2).Infof("Starting single player game")
	r.started = true

	r.deck = deck
	err := r.initializePlayer()
	if err != nil {
		return errors.Wrapf(err, "unable to start game")
	}

	// From here on, the game logic could throw exception to signal that the game is over.
	r.startGameSequence()
	return nil
}

func (r *SinglePlayerRunner) initializePlayer() error {
	glog.V(2).Infof("Initializing hand\n")
	hand := domain.NewHand()
	err := rules.PopulateHands(r.deck, []*domain.Hand{hand})
	if err != nil {
		return err
	}

	hand.Sort()
	// In single player mode, the player's prevailing wind is always East (0).
	windOrdinal := 0
	const numWindOrdinals = 4

	r.player = rules.NewPlayerGameState(hand, windOrdinal)
	r.pseudoOpponentGameState = rules.NewPlayerGameState(
		domain.NewHand(), (windOrdinal+1)%numWindOrdinals)
	glog.V(2).Infof("Populated hand: %s\n", hand)
	return nil
}

func (r *SinglePlayerRunner) startGameSequence() {
	defer func() {
		if r := recover(); r != nil {
			if gameOverErr, ok := r.(*gameOverError); ok {
				// TODO: implement proper panic recovery.
				glog.V(2).Infof("Game over: %s\n", gameOverErr)
				fmt.Printf("Game over: %s\n", gameOverErr)
			} else {
				panic(r)
			}
		}
	}()

	glog.V(2).Infof("Starting game sequence\n")
	// TODO: notify observer of game start

	numTilesToReplace := r.bulkMoveBonusTilesFromHand()
	if numTilesToReplace > 0 {
		replacementRound := 1
		for numTilesToReplace > 0 {
			glog.V(2).Infof("Replacing %d bonus tiles (round %d)\n",
				numTilesToReplace, replacementRound)
			for i := 0; i < numTilesToReplace; i++ {
				tile := r.drawFromDeckBack()
				r.addTileToHandNoCheck(tile)
			}
			numTilesToReplace = r.bulkMoveBonusTilesFromHand()
			replacementRound++
		}
	}

	r.sortHand()
	fmt.Printf("Hand after replacement: %s\n", r.player.GetHand())

	r.startPlayerRoundLoop()
}

func (r *SinglePlayerRunner) startPlayerRoundLoop() {
	round := 1
	playerMelded := false
	for {
		// TODO: Notify observer of round start
		fmt.Println(separator)
		fmt.Printf("Start of round %d\n", round)

		if round == 1 || !playerMelded {
			var source rules.OutTileSourceType
			var tile *domain.Tile
			if round == 1 {
				source = rules.OutTileSourceTypeInitialHand
			} else {
				// Draw phase
				tile = r.drawFromDeckFront()
				source = rules.OutTileSourceTypeSelfDrawn
				if round > 1 {
					fmt.Printf("Drawn tile %s\n", tile)
				}
				if rules.IsEligibleForHand(tile.GetSuit()) {
					r.addTileToHand(tile)
				} else {
					r.addTileToBonusArea(tile)
					tile = r.replaceTileLoop()
					source = rules.OutTileSourceTypeSelfDrawnReplacement
				}
			}

			fmt.Printf("Hand: %s\n", r.player.GetHand())
			// Player action phase
			r.promptAndExecutePlayerAction(commandsAfterDrawingTile,
				rules.NewOutTileSource(source, tile, nil))
		}

		playerMelded = false

		// Burn phase
		for x := 0; x < r.numBurnsPerRound; x++ {
			fmt.Printf("Burn %d of %d in round %d\n", x+1, r.numBurnsPerRound, round)
			// Allow chow in the last burn
			chowAllowed := x == r.numBurnsPerRound-1
			melded := r.burnSingleTile(chowAllowed)
			if melded {
				fmt.Printf("Exiting burn phase due to meld\n")
				playerMelded = true
				break
			}
		}
		round++
	}
}

func (r *SinglePlayerRunner) burnSingleTile(chowAllowed bool) bool {
	if r.currentBurnTile != nil {
		panic("There shouldn't be an active burn tile")
	}

	tile := r.drawFromDeckFront()
	round := 1
	for !rules.IsEligibleForHand(tile.GetSuit()) {
		r.addTileToOtherBonusArea(tile)

		fmt.Printf("Drawing a replacement tile from the back of deck (round %d)\n", round)
		tile = r.drawFromDeckBack()
		round++
	}

	// TODO: add observer for "other" player about to discard tile
	fmt.Printf("Burning tile %s\n", tile)
	r.currentBurnTile = tile

	discardInfo := rules.NewDiscardInfo(r.pseudoOpponentGameState)
	var cmdType ui.CommandType
	if chowAllowed {
		cmdType = r.promptAndExecutePlayerAction(
			withCommands(ui.Pong, ui.Kong, ui.Chow, ui.Pass, ui.Out),
			rules.NewOutTileSource(rules.OutTileSourceTypeDiscard, tile, discardInfo))
	} else {
		cmdType = r.promptAndExecutePlayerAction(
			withCommands(ui.Pong, ui.Kong, ui.Pass, ui.Out),
			rules.NewOutTileSource(rules.OutTileSourceTypeDiscard, tile, discardInfo))
	}

	melded := cmdType != ui.Pass
	if !melded {
		r.discardCurrentBurnTile()
	}

	return melded
}

func (r *SinglePlayerRunner) replaceTileLoop() *domain.Tile {
	for round := 1; ; /* no-op */ round++ {
		fmt.Printf("Drawing a replacement tile from the back of deck (round %d)\n", round)
		tile := r.drawFromDeckBack()
		if rules.IsEligibleForHand(tile.GetSuit()) {
			fmt.Printf("Adding replacement tile to hand: %s\n", tile)
			r.addTileToHand(tile)
			return tile
		}
		// Else tile is a bonus tile, add it and repeat.
		r.addTileToBonusArea(tile)
	}
}

func (r *SinglePlayerRunner) promptAndExecutePlayerAction(
	acceptedCommands ui.CommandTypes, outTileSource *rules.OutTileSource) ui.CommandType {
	for {
		cmd, err := r.receiver.PromptForCommand(acceptedCommands)
		if err != nil {
			// TODO: this should be its own error struct. Something like IOError.
			panic(newGameOverError(false))
		}

		proceed := r.executePlayerAction(cmd, outTileSource)
		if proceed {
			return cmd.GetCommandType()
		}
	}
}

func (r *SinglePlayerRunner) executePlayerAction(cmd *ui.Command,
	outTileSource *rules.OutTileSource) bool {
	switch cmd.GetCommandType() {
	case ui.SortHand:
		r.sortHand()
		return false
	case ui.ShowDiscardedTiles:
		r.showDiscardedTiles()
		return false
	case ui.ShowMelded:
		r.showMelded()
		return false
	case ui.DiscardTile:
		return r.discardTile(cmd.GetTileIndexCommand().GetIndex())
	case ui.Pong:
		return r.declarePong()
	case ui.Kong:
		return r.declareKong()
	case ui.ConcealedKong:
		return r.declareConcealedKong(cmd.GetTileIndexCommand().GetIndex())
	case ui.AdditionalKong:
		return r.declareAdditionalKong(cmd.GetTileIndexCommand().GetIndex())
	case ui.Chow:
		return r.declareChow(
			cmd.GetTileIndexCommand2().GetIndex1(), cmd.GetTileIndexCommand2().GetIndex2())
	case ui.Pass:
		return true
	case ui.Out:
		return r.checkForOut(outTileSource)
	}
	fmt.Printf("Unhandled command: %s\n", cmd.GetCommandType())
	return false
}

// checkForOut checks whether the current player state represents an Out. This function will
// panic if the hand is an out hand, or return false if it is not an out hand.
func (r *SinglePlayerRunner) checkForOut(outTileSource *rules.OutTileSource) bool {
	counter := rules.NewOutPlanCalculator(rules.GetSuitsForGame(), r.player, outTileSource)
	plans := counter.Calculate()

	if len(plans) > 0 {
		fmt.Printf("Out: %s.\n", outTileSource)
		if *flags.ReportScoringFlag {
			scorer := zj.NewOutPlansScorer()
			context := rules.NewOutPlanScoringContext(
				outTileSource, r.player, r.deck.NumRemainingTiles())
			scoredPlans := scorer.ScoreOutPlans(plans, context)
			fmt.Printf("Detailed scoring:\n")
			fmt.Printf("%s\n", scoredPlans)
		}
		panic(newGameOverError(true))
	}
	fmt.Println("Not an Out hand!")
	return false
}

// Methods that manipulate player / deck state that should notify the observer.

func (r *SinglePlayerRunner) sortHand() {
	r.player.SortHand()
	// TODO: Notify observer of hand sorted update.
	fmt.Printf("Hand: %s\n", r.player.GetHand())
}

func (r *SinglePlayerRunner) showDiscardedTiles() {
	fmt.Printf("Discarded tiles: %s\n", r.player.GetDiscardedTiles())
}

func (r *SinglePlayerRunner) showMelded() {
	fmt.Printf("Melded groups: %s\n", r.player.GetMeldGroups())
}

func (r *SinglePlayerRunner) addTileToHand(t *domain.Tile) {
	// TODO: add observer to update hand
	r.player.AddTileToHand(t)
}

func (r *SinglePlayerRunner) addTileToHandNoCheck(t *domain.Tile) {
	// TODO: add observer to update hand
	r.player.AddTileToHandNoCheck(t)
}

func (r *SinglePlayerRunner) discardTile(index int) bool {
	t, removed := r.player.DiscardTileAt(index)
	if removed {
		// TODO: notify observer
		fmt.Printf("Discarded tile at %d: %s\n", index, t)
	}
	return removed
}

func (r *SinglePlayerRunner) discardCurrentBurnTile() {
	if r.currentBurnTile == nil {
		panic(fmt.Errorf("There is no tile being burned"))
	}

	// TODO: notify observer
	fmt.Printf("Moving burn tile %s to other discards\n", r.currentBurnTile)
	r.pseudoOpponentGameState.AddTileToHand(r.currentBurnTile)
	r.pseudoOpponentGameState.DiscardTileAt(0)
	r.currentBurnTile = nil
}

func (r *SinglePlayerRunner) declarePong() bool {
	if r.currentBurnTile == nil {
		glog.V(2).Infof("There is no tile being burned")
		return false
	}

	removed := r.player.DeclarePong(r.currentBurnTile)
	if !removed {
		fmt.Printf("Failed to declare pong\n")
		return false
	}
	// TODO: notify observer of meld area / hand change
	fmt.Printf("Declared pong %s\n", r.currentBurnTile)
	r.currentBurnTile = nil
	r.promptAndExecutePlayerAction(withCommands(ui.DiscardTile), nil)
	return removed
}

func (r *SinglePlayerRunner) declareKong() bool {
	if r.currentBurnTile == nil {
		glog.V(2).Infof("There is no tile being burned")
		return false
	}

	removed := r.player.DeclareKong(r.currentBurnTile)
	if !removed {
		fmt.Printf("Failed to declare kong\n")
		return false
	}

	// TODO: notify observer of meld area / hand change
	fmt.Printf("Declared kong %s\n", r.currentBurnTile)
	r.currentBurnTile = nil

	// After drawing the replacement tile, the player may go out, or they must discard a tile.
	replacementTile := r.replaceTileLoop()
	r.promptAndExecutePlayerAction(commandsAfterDrawingTile,
		rules.NewOutTileSource(rules.OutTileSourceTypeSelfDrawnReplacement, replacementTile, nil))
	return removed
}

func (r *SinglePlayerRunner) declareConcealedKong(index int) bool {
	t, removed := r.player.DeclareConcealedKong(index)
	if !removed {
		glog.V(2).Infof("Failed to declare concealed kong with tile at %d\n", index)
		return false
	}

	// TODO: notify observer of meld area / hand change
	fmt.Printf("Declared concealed kong %s\n", t)

	// After drawing the replacement tile, the player may go out, or have another concealed kong.
	// Note this may result in a recursion.
	// TODO: don't do recursion?
	replacementTile := r.replaceTileLoop()
	r.promptAndExecutePlayerAction(commandsAfterDrawingTile,
		rules.NewOutTileSource(rules.OutTileSourceTypeSelfDrawnReplacement, replacementTile, nil))
	return removed
}

func (r *SinglePlayerRunner) declareAdditionalKong(index int) bool {
	t, removed := r.player.DeclareAdditionalKong(index)
	if !removed {
		glog.V(2).Infof("Failed to declare additional kong with tile at %d\n", index)
		return false
	}

	// TODO: notify observer of meld area / hand change
	fmt.Printf("Declared additional kong %s\n", t)

	// After drawing the replacement tile, the player may go out, or have another concealed kong.
	// Note this may result in a recursion.
	// TODO: don't do recursion?
	replacementTile := r.replaceTileLoop()
	r.promptAndExecutePlayerAction(commandsAfterDrawingTile,
		rules.NewOutTileSource(rules.OutTileSourceTypeSelfDrawnReplacement, replacementTile, nil))
	return removed
}

func (r *SinglePlayerRunner) declareChow(index1, index2 int) bool {
	if r.currentBurnTile == nil {
		glog.V(2).Infof("There is no tile being burned")
		return false
	}

	tiles, removed := r.player.DeclareChow(r.currentBurnTile, index1, index2)
	if !removed {
		fmt.Printf("Failed to declared chow\n")
		return false
	}

	fmt.Printf("Declared chow %s\n", tiles)
	// TODO: notify observer of meld area / hand change
	r.currentBurnTile = nil
	r.promptAndExecutePlayerAction(withCommands(ui.DiscardTile), nil)
	return true
}

func (r *SinglePlayerRunner) bulkMoveBonusTilesFromHand() int {
	// TODO: notify observer of moved tiles.
	return r.player.BulkMoveBonusTilesFromHand()
}

func (r *SinglePlayerRunner) addTileToBonusArea(t *domain.Tile) {
	fmt.Printf("Adding tile to bonus area: %s\n", t)
	r.player.AddTileToBonusArea(t)
	// TODO: notify observer
}

func (r *SinglePlayerRunner) addTileToOtherBonusArea(t *domain.Tile) {
	fmt.Printf("Adding tile to other bonus area: %s\n", t)
	r.pseudoOpponentGameState.AddTileToBonusArea(t)
	// TODO: notify observer
}

func (r *SinglePlayerRunner) drawFromDeckFront() *domain.Tile {
	// TODO: notify observer of deck pop
	tile, err := r.deck.PopFront()
	if err != nil {
		panic(newGameOverError(false))
	}
	return tile
}

func (r *SinglePlayerRunner) drawFromDeckBack() *domain.Tile {
	// TODO: notify observer of deck pop
	tile, err := r.deck.PopBack()
	if err != nil {
		panic(newGameOverError(false))
	}
	return tile
}
