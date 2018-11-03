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
	// Out declares the player has reached an out hand.
	Out CommandType = "out"
)

// DiscardTileCommand represents information of a DiscardTile command.
type DiscardTileCommand struct {
	// The index of the tile to be discarded.
	index int
}

// GetIndex ...
func (d *DiscardTileCommand) GetIndex() int {
	return d.index
}

// Command represents a command from an input that may modify the state of the game.
type Command struct {
	commandType CommandType
	// discardTile is only set if commandType is DiscardTile.
	discardTile *DiscardTileCommand
}

// GetCommandType ...
func (c *Command) GetCommandType() CommandType {
	return c.commandType
}

// GetDiscardTileCommand ...
func (c *Command) GetDiscardTileCommand() *DiscardTileCommand {
	if c.commandType != DiscardTile {
		fmt.Printf("ERROR: Command is not of type DiscardTile: %s\n", c.commandType)
		return nil
	}
	return c.discardTile
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
	return &Command{commandType: DiscardTile, discardTile: &DiscardTileCommand{index: index}}
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
