package flags

import (
	"flag"
)

//// App level flags

// ModeFlag specifies the app mode.
var ModeFlag = flag.String("mj.mode", "", "App mode")

// AppMode specifies the app mode.
type AppMode = string

const (
	// AppModeDeck dumps a shuffled deck and exits.
	AppModeDeck AppMode = "deck"
	// AppModeSingle runs the single player mode.
	AppModeSingle AppMode = "single"
)

// RuleNameFlag specifies the MJ rule name.
var RuleNameFlag = flag.String("mj.ruleName", "zj", "Name of MJ rule to use")

// RuleName specifies the MJ rule name.
type RuleName = string

const (
	// RuleNameHK is Hong Kong MJ.
	RuleNameHK RuleName = "hk"
	// RuleNameZJ is Zung Jung MJ.
	RuleNameZJ RuleName = "zj"
)

//// Single player mode flags

// NumBurnsFlag specifies number of tiles to burn in each round in single player mode.
var NumBurnsFlag = flag.Int("mj.numBurns", 0, "Number of tiles to burn in each round")
