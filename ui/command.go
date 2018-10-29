package ui

// CommandType represents the set of possible commands.
type CommandType int

const (
	// SortHand sorts the tiles in the hand.
	SortHand CommandType = iota
	// ShowDiscardedTiles shows the tiles that have been discarded so far.
	ShowDiscardedTiles
	// DiscardTile discards a tile at the given index. Corresponds to DiscardTileCommand.
	DiscardTile
)

type discardTileCommand struct {
	// The index of the tile to be discarded.
	index int
}

// Command represents a command from an input that may modify the state of the game.
type Command struct {
	commandType CommandType
	// discardTile is only set if commandType is DiscardTile.
	discardTile *discardTileCommand
}

// SortHandCommand returns a new SortHand command.
func SortHandCommand() *Command {
	return &Command{commandType: SortHand}
}

// ShowDiscardedTilesCommand returns a new ShowDiscardedTiles command.
func ShowDiscardedTilesCommand() *Command {
	return &Command{commandType: ShowDiscardedTiles}
}

// DiscardTileCommand returns a new DiscardTile command with the given index.
func DiscardTileCommand(index int) *Command {
	return &Command{commandType: DiscardTile, discardTile: &discardTileCommand{index: index}}
}

// CommandReceiver is an interface for receiving commands from an input.
type CommandReceiver interface {
	// PromptForCommand prompts for a command from the input source. Returns the Command, or an
	// error if the receiver encountered an error.
	PromptForCommand() (*Command, error)
}
