package main

import (
	"fmt"
	"github.com/derekimcheng/mj/rules"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("mj Hello world")
	rand.Seed(time.Now().UTC().UnixNano())

	CreateAndDumpDeck()
}

// CreateAndDumpDeck creates a game deck and empties it, logging each tile in the order they are
// drawn.
func CreateAndDumpDeck() {
	deck := rules.NewDeckForGame()
	deck.Shuffle()
	for !deck.IsEmpty() {
		tile, err := deck.PopFront()
		if err != nil {
			panic("Failed to draw from deck")
		}
		fmt.Printf("Got tile %s\n", tile)
	}
}
