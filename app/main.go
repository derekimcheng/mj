package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("mj Hello world")
	rand.Seed(time.Now().UTC().UnixNano())
}
