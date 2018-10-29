package main

import (
	"flag"
	"fmt"
	"github.com/derekimcheng/mj/domain"
	"github.com/derekimcheng/mj/rules"
	"math/rand"
	"os"
	"time"
)

var modeFlag = flag.String("mode", "", "App mode")

func main() {
	fmt.Println("mj Hello world")
	initialize()

	switch *modeFlag {
	case "deck":
		createAndDumpDeck()
	case "single":
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
	fmt.Println("usage: ./app -mode={mode}")
	fmt.Println("mode can be one of: deck, single")
}

// createAndDumpDeck creates a game deck and empties it, logging each tile in the order they are
// drawn.
func createAndDumpDeck() {
	deck := rules.NewDeckForGame()
	deck.Shuffle()
	for !deck.IsEmpty() {
		tile, err := deck.PopFront()
		if err != nil {
			panic("Failed to draw from deck")
		}
		fmt.Printf("Got tile [%s]\n", tile)
	}
}

func simulateSingleHand() {
	deck := rules.NewDeckForGame()
	deck.Shuffle()

	hand := domain.NewHand()
	rules.PopulateHands(deck, []*domain.Hand{hand})
	fmt.Printf("Populated hand: %s\n", hand)
}
