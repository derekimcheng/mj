package ui

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var commandMapping = map[string]CommandType{
	"sort":      SortHand,
	"discarded": ShowDiscardedTiles,
	"discard":   DiscardTile,
}

// ConsoleCommandReceiver receives command from the an input stream, such as the console.
type ConsoleCommandReceiver struct {
	scanner *bufio.Scanner
}

// NewConsoleCommandReceiver creates a new ConsoleCommandReceiver with the given input source.
func NewConsoleCommandReceiver(r io.Reader) *ConsoleCommandReceiver {
	return &ConsoleCommandReceiver{
		scanner: bufio.NewScanner(r),
	}
}

// PromptForCommand ... (CommandReceiver implementation)
func (recver *ConsoleCommandReceiver) PromptForCommand() (*Command, error) {
	// Repeat until an error is encountered or a valid Command is obtained.
	for {
		fmt.Printf("Enter a command: ")
		success := recver.scanner.Scan()
		if !success {
			err := recver.scanner.Err()
			if err == nil {
				err = io.EOF
			}
			return nil, err
		}

		text := recver.scanner.Text()
		cmd, err := parseCommand(text)
		if err == nil {
			return cmd, nil
		}
		fmt.Printf("Received error from parsing command: %s\n", err)
	}
}

func parseCommand(input string) (*Command, error) {
	fields := strings.Fields(input)
	cmdStr := fields[0]
	args := fields[1:]
	cmd, found := commandMapping[cmdStr]
	if !found {
		return nil, fmt.Errorf("Received unknown command %s", cmdStr)
	}

	switch cmd {
	case SortHand:
		return SortHandCommand(), nil
	case ShowDiscardedTiles:
		return ShowDiscardedTilesCommand(), nil
	case DiscardTile:
		if len(args) < 1 {
			return nil, fmt.Errorf("Not enough args for DiscardTile")
		}
		index, err := strconv.Atoi(args[0])
		if err != nil || index < 0 {
			return nil, fmt.Errorf("Invalid arg for DiscardTile: %s", args[0])
		}
		return DiscardTileCommand(index), nil
	}
	return nil, fmt.Errorf("Unhandled command %s", cmdStr)
}
