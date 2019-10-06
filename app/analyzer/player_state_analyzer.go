package analyzer

import (
	"bufio"
	"fmt"
	"github.com/derekimcheng/mj/domain"
	"github.com/derekimcheng/mj/rules"
	"github.com/derekimcheng/mj/rules/zj"
	"github.com/derekimcheng/mj/shorthand"
	"io"
	"strconv"
)

var outSourceTypeOptionsStr = "d, sd, sdr, ak, ih"
var outSourceTypeMap = map[string]rules.OutTileSourceType{
	"d":   rules.OutTileSourceTypeDiscard,
	"sd":  rules.OutTileSourceTypeSelfDrawn,
	"sdr": rules.OutTileSourceTypeSelfDrawnReplacement,
	"ak":  rules.OutTileSourceTypeAdditionalKong,
	"ih":  rules.OutTileSourceTypeInitialHand,
}

// PlayerStateAnalyzer analyzes the given player state and scores it.
type PlayerStateAnalyzer struct {
	scanner         *bufio.Scanner
	shortHandParser *shorthand.Parser
}

// NewPlayerStateAnalyzer returns a new PlayerStateAnalyzer.
// TODO: add rules to input?
func NewPlayerStateAnalyzer(reader io.Reader) *PlayerStateAnalyzer {
	return &PlayerStateAnalyzer{
		scanner:         bufio.NewScanner(reader),
		shortHandParser: shorthand.NewParser(),
	}
}

// Start ...
func (p *PlayerStateAnalyzer) Start() {
	if err := p.doStart(); err != nil {
		fmt.Printf("Encountered error: %s\n", err)
	}
}

func (p *PlayerStateAnalyzer) doStart() error {
	// Input hand
	hand, err := p.inputHand()
	if err != nil {
		return err
	}

	// Input melded area
	meldGroups, err := p.inputMeldGroups()
	if err != nil {
		return err
	}

	// Input player wind seat
	windOrdinal, err := p.inputWind("player")
	if err != nil {
		return err
	}

	// Input out criteria (including the out tile)
	outTileSource, isLastTile, err := p.inputOutTileSource()
	if err != nil {
		return err
	}

	// Validate (melds are valid, number of tiles in hand+meld is valid, max 4 tiles each)
	// TODO: implement a validator.
	// Add self-drawn tile to hand so it is picked up by the calculator.
	if rules.IsSelfDrawnType(outTileSource.SourceType) {
		hand.AddTile(outTileSource.Tile)
	}

	// Score and list out plans
	playerGameState := rules.NewExistingPlayerGameState(hand, windOrdinal, nil, meldGroups)
	calc := rules.NewOutPlanCalculator(rules.GetSuitsForGame(), playerGameState, outTileSource)
	plans := calc.Calculate()

	fmt.Printf("Found %d out plans\n", len(plans))
	if len(plans) > 0 {
		// Hack to simulate last-tile.
		numRemainingTiles := 42
		if isLastTile {
			numRemainingTiles = 0
		}

		scorer := zj.NewOutPlansScorer()
		context := rules.NewOutPlanScoringContext(outTileSource, playerGameState, numRemainingTiles)
		scoredPlans := scorer.ScoreOutPlans(plans, context)
		fmt.Printf("Detailed scoring:\n")
		fmt.Printf("%s\n", scoredPlans)
	}

	return nil
}

func (p *PlayerStateAnalyzer) inputHand() (*domain.Hand, error) {
	for {
		str, err := p.promptForInput("Enter hand tiles, not including the out tile", "")
		if err != nil {
			return nil, err
		}
		tiles, err := p.shortHandParser.ParseTiles(str)
		if err != nil {
			fmt.Printf("Error parsing hand tiles: %s\n", err)
			continue
		}
		numTiles := len(tiles)
		// This value is typically 1, but could also be 2 if the player wins with initial hand.
		if numTiles%3 == 0 {
			fmt.Printf("Invalid number of tiles in hand: %d\n", numTiles)
			continue
		}
		hand := domain.NewHand()
		hand.SetTiles(tiles)
		return hand, nil
	}
}

func (p *PlayerStateAnalyzer) inputMeldGroups() (rules.TileGroups, error) {
	for {
		str, err := p.promptForInput("Enter meld groups", "")
		if err != nil {
			return nil, err
		}
		groups, err := p.shortHandParser.ParseMeldGroups(str)
		if err != nil {
			fmt.Printf("Error parsing meld groups: %s\n", err)
			continue
		}
		if len(groups) > 4 {
			fmt.Printf("Invalid number of meld groups: %d\n", len(groups))
			continue
		}
		return groups, nil
	}
}

func (p *PlayerStateAnalyzer) inputWind(who string) (int, error) {
	for {
		str, err := p.promptForInput(
			fmt.Sprintf("Input %s wind seat (1=E, 2=S, 3=W, 4=N)", who), "1")
		if err != nil {
			return 0, err
		}
		windSeat, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Invalid wind seat %s\n", str)
			continue
		}

		windOrdinal := windSeat - 1
		if windOrdinal < 0 || windOrdinal >= rules.Winds.GetSize() {
			fmt.Printf("Wind seat out of range: %d\n", windSeat)
			continue
		}
		return windOrdinal, nil
	}
}

func (p *PlayerStateAnalyzer) inputOutTileSource() (*rules.OutTileSource, bool, error) {
	for {
		str, err := p.promptForInput(fmt.Sprintf("Input out source (%s)", outSourceTypeOptionsStr), "")
		if err != nil {
			return nil, false, err
		}
		outSourceType, found := outSourceTypeMap[str]
		if !found {
			fmt.Printf("Unknown out source type %s\n", str)
			continue
		}

		var outTile *domain.Tile
		isLastTile := false
		if outSourceType != rules.OutTileSourceTypeInitialHand {
			outTile, err = p.inputOutTile()
			if err != nil {
				fmt.Printf("Error parsing out tile: %s\n", err)
				return nil, false, err
			}
			// Parse isLastTile
			isLastTile, err = p.inputBool("Is last tile", false)
			if err != nil {
				fmt.Printf("Error parsing isLastTile bool: %s\n", err)
				return nil, false, err
			}
		}

		// Parse discard player
		var discardInfo *rules.DiscardInfo
		if rules.IsExternalOutSourceType(outSourceType) {
			discardPlayerWindOrdinal, err := p.inputWind("discarder")
			if err != nil {
				fmt.Printf("Error parsing discarder wind: %s\n", err)
				return nil, false, err
			}
			isFirstDiscard, err := p.inputBool("Is first discard", false)
			if err != nil {
				fmt.Printf("Error parsing isFirstDiscard bool: %s\n", err)
				return nil, false, err
			}
			var discardedTiles domain.Tiles
			// Hack to emulate that this is not a first discard.
			if !isFirstDiscard {
				discardedTiles = append(discardedTiles, nil)
			}
			discardPlayer := rules.NewExistingPlayerGameState(
				domain.NewHand(), discardPlayerWindOrdinal, discardedTiles, nil)
			discardInfo = rules.NewDiscardInfo(discardPlayer)
		}

		return rules.NewOutTileSource(outSourceType, outTile, discardInfo), isLastTile, nil
	}
}

func (p *PlayerStateAnalyzer) inputOutTile() (*domain.Tile, error) {
	for {
		str, err := p.promptForInput("Enter out tile", "")
		if err != nil {
			return nil, err
		}
		tiles, err := p.shortHandParser.ParseTiles(str)
		if err != nil {
			fmt.Printf("Error parsing out tile: %s\n", err)
			continue
		}
		numTiles := len(tiles)
		if numTiles != 1 {
			fmt.Printf("Invalid number of tiles in hand, expected 1: %d\n", numTiles)
			continue
		}
		return tiles[0], nil
	}
}

func (p *PlayerStateAnalyzer) inputBool(prompt string, defaultValue bool) (bool, error) {
	for {
		str, err := p.promptForInput(prompt, strconv.FormatBool(defaultValue))
		if err != nil {
			return false, err
		}
		value, err := strconv.ParseBool(str)
		if err != nil {
			fmt.Printf("Uknown bool value %s\n", str)
			continue
		}
		return value, nil
	}
}

func (p *PlayerStateAnalyzer) promptForInput(prompt string, defaultValue string) (string, error) {
	for {
		fmt.Printf("%s [default='%s']: ", prompt, defaultValue)
		success := p.scanner.Scan()
		if !success {
			err := p.scanner.Err()
			if err == nil {
				err = io.EOF
			}
			return "", err
		}

		text := p.scanner.Text()
		if len(text) == 0 {
			return defaultValue, nil
		}
		return text, nil
	}
}
