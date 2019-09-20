package main

import (
	"flag"
	"fmt"
	"github.com/derekimcheng/mj/domain"
	"github.com/derekimcheng/mj/engine"
	"github.com/derekimcheng/mj/flags"
	"github.com/derekimcheng/mj/rules"
	"github.com/derekimcheng/mj/ui"
	"github.com/pkg/errors"
	"math/rand"
	"os"
	"time"
)

func main() {
	fmt.Println("mj Hello world")
	initialize()

	switch *flags.ModeFlag {
	case flags.AppModeDeck:
		createAndDumpDeck()
	case flags.AppModeSingle:
		simulateSingleHand()
	default:
		printUsage()
		os.Exit(1)
	}
	os.Exit(0)
}

func initialize() {
	flag.Parse()
	rand.Seed(time.Now().UTC().UnixNano())
}

func printUsage() {
	fmt.Println("usage: see ./app -help")
}

// createAndDumpDeck creates a game deck and empties it, logging each tile in the order they are
// drawn.
func createAndDumpDeck() {
	deck := createDeck()
	for !deck.IsEmpty() {
		tile, err := deck.PopFront()
		if err != nil {
			panic("Failed to draw from deck")
		}
		fmt.Printf("Got tile [%s]\n", tile)
	}
}

func simulateSingleHand() {
	runner := engine.NewSinglePlayerRunner(ui.NewConsoleCommandReceiver(os.Stdin))
	err := runner.Start(createDeck())
	if err != nil {
		fmt.Printf("Encountered error while running single player game: %s\n", err)
	}
}

func createDeck() domain.Deck {
	deck, err := rules.NewDeckForGame(*flags.RuleNameFlag)
	if err != nil {
		panic(errors.Wrapf(err, "failed to initialize deck"))
	}
	if deck.IsEmpty() {
		panic(errors.New("Deck is empty"))
	}
	deck.Shuffle()
	return deck
}