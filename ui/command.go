package ui

import (
	"fmt"
)

// CommandType represents the set of possible commands.
type CommandType = string

const (
	// SortHand sorts the tiles in the hand.
	SortHand CommandType = "sort"
	// ShowDiscardedTiles shows the tiles that have been discarded so far.
	ShowDiscardedTiles CommandType = "discarded"
	// DiscardTile discards a tile at the given index. Corresponds to DiscardTileCommand.
	DiscardTile CommandType = "discard"
	// Pong creates a meld from a pong tile group. Only available if the tile completing the
	// meld is discarded by other players.
	Pong CommandType = "pong"
	// Kong creates a meld from a kong tile group. Only available if the tile completing the
	// meld is discarded by other players.
	Kong CommandType = "kong"
	// ConcealedKong creates a concealed meld from a kong tile group. Only available if the tile
	// completing the meld is drawn.
	ConcealedKong CommandType = "ckong"
	// AdditionalKong converts a melded pong to a kong. Only available if the tile completing the
	// kong is drawn.
	AdditionalKong CommandType = "akong"
	// Chow creates a meld from a chow tile group. Only available if the tile completing the
	// meld is discarded by previous player.
	Chow CommandType = "chow"
	// Pass allows the player to pass on a tile about to be discarded by other players.
	Pass CommandType = "pass"
	// Out declares the player has reached an out hand.
	Out CommandType = "out"
)

// TileIndexCommand represents information of a command that specifies a tile index.
type TileIndexCommand struct {
	// index is the index of the tile to use in the command.
	index int
}

// GetIndex ...
func (c *TileIndexCommand) GetIndex() int {
	return c.index
}

// TileIndexCommand2 represents information of a command that specifies 2 indices.
type TileIndexCommand2 struct {
	// index1 is the index of the first of the 2 tiles to use in the command.
	index1 int
	// index2 is the index of the second of the 2 tiles to use in the command.
	index2 int
}

// GetIndex1 ...
func (c *TileIndexCommand2) GetIndex1() int {
	return c.index1
}

// GetIndex2 ...
func (c *TileIndexCommand2) GetIndex2() int {
	return c.index2
}

// Command represents a command from an input that may modify the state of the game.
type Command struct {
	commandType CommandType
	// tile is only set if commandType is DiscardTile / ConcealedKong / AdditionalKong.
	tile *TileIndexCommand
	// tile is only set if commandType is Chow.
	tile2 *TileIndexCommand2
}

// GetCommandType ...
func (c *Command) GetCommandType() CommandType {
	return c.commandType
}

// GetTileIndexCommand ...
func (c *Command) GetTileIndexCommand() *TileIndexCommand {
	if c.commandType != DiscardTile &&
		c.commandType != ConcealedKong &&
		c.commandType != AdditionalKong {
		panic(fmt.Errorf("invalid command type for TileIndexCommand: %s", c.commandType))
	}
	return c.tile
}

// GetTileIndexCommand2 ...
func (c *Command) GetTileIndexCommand2() *TileIndexCommand2 {
	if c.commandType != Chow {
		panic(fmt.Errorf("invalid command type for TileIndexCommand2: %s", c.commandType))
	}
	return c.tile2
}

// NewSortHandCommand returns a new SortHand command.
func NewSortHandCommand() *Command {
	return &Command{commandType: SortHand}
}

// NewShowDiscardedTilesCommand returns a new ShowDiscardedTiles command.
func NewShowDiscardedTilesCommand() *Command {
	return &Command{commandType: ShowDiscardedTiles}
}

// NewDiscardTileCommand returns a new DiscardTile command with the given index.
func NewDiscardTileCommand(index int) *Command {
	return &Command{commandType: DiscardTile, tile: &TileIndexCommand{index: index}}
}

// NewPongCommand returns a new Pong command.
func NewPongCommand() *Command {
	return &Command{commandType: Pong}
}

// NewKongCommand returns a new Kong command.
func NewKongCommand() *Command {
	return &Command{commandType: Kong}
}

// NewConcealedKongCommand returns a new ConcealedKong command with the given index.
func NewConcealedKongCommand(index int) *Command {
	return &Command{commandType: ConcealedKong, tile: &TileIndexCommand{index: index}}
}

// NewAdditionalKongCommand returns a new AdditionalKong command with the given index.
func NewAdditionalKongCommand(index int) *Command {
	return &Command{commandType: AdditionalKong, tile: &TileIndexCommand{index: index}}
}

// NewChowCommand returns a new Chow command with the given indices.
func NewChowCommand(index1, index2 int) *Command {
	return &Command{commandType: Chow,
		tile2: &TileIndexCommand2{index1: index1, index2: index2}}
}

// NewPassCommand retrurns a new Out command.
func NewPassCommand() *Command {
	return &Command{commandType: Pass}
}

// NewOutCommand retrurns a new Out command.
func NewOutCommand() *Command {
	return &Command{commandType: Out}
}

// CommandTypes is a slice of CommandTypes.
type CommandTypes []CommandType

// ContainsCommand returns true if the CommandTypes contains the given CommandType.
func (commands CommandTypes) ContainsCommand(command CommandType) bool {
	for _, c := range commands {
		if c == command {
			return true
		}
	}
	return false
}

// CommandReceiver is an interface for receiving commands from an input.
type CommandReceiver interface {
	// PromptForCommand prompts for a command from the input source. Returns the Command, or an
	// error if the receiver encountered an error. The returned command's type must be one of the
	// given CommandTypes.
	PromptForCommand(acceptedCommands CommandTypes) (*Command, error)
}
