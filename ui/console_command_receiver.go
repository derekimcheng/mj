package ui

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

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
func (recver *ConsoleCommandReceiver) PromptForCommand(acceptedCommands CommandTypes) (*Command, error) {
	// Repeat until an error is encountered or a valid Command is obtained.
	for {
		fmt.Printf("Enter a command [%s]: ", strings.Join(acceptedCommands, "|"))
		success := recver.scanner.Scan()
		if !success {
			err := recver.scanner.Err()
			if err == nil {
				err = io.EOF
			}
			return nil, err
		}

		text := recver.scanner.Text()
		if len(text) == 0 {
			continue
		}

		cmd, err := parseCommand(text)
		if err != nil {
			fmt.Printf("Received error from parsing command: %s\n", err)
			continue
		}
		if !acceptedCommands.ContainsCommand(cmd.GetCommandType()) {
			fmt.Printf("Unacceptable command %s\n", cmd.GetCommandType())
			continue
		}
		return cmd, nil
	}
}

func parseCommand(input string) (*Command, error) {
	fields := strings.Fields(input)
	if len(fields) < 1 {
		return nil, fmt.Errorf("Fewer than 1 field in input")
	}
	cmdStr := fields[0]
	args := fields[1:]

	switch cmdStr {
	case SortHand:
		return NewSortHandCommand(), nil
	case ShowDiscardedTiles:
		return NewShowDiscardedTilesCommand(), nil
	case DiscardTile:
		if len(args) < 1 {
			return nil, fmt.Errorf("Not enough args for DiscardTile")
		}
		index, err := strconv.Atoi(args[0])
		if err != nil || index < 0 {
			return nil, fmt.Errorf("Invalid arg for DiscardTile: %s", args[0])
		}
		return NewDiscardTileCommand(index), nil
	case Pong:
		return NewPongCommand(), nil
	case Kong:
		return NewKongCommand(), nil
	case ConcealedKong:
		if len(args) < 1 {
			return nil, fmt.Errorf("Not enough args for ConcealedKong")
		}
		index, err := strconv.Atoi(args[0])
		if err != nil || index < 0 {
			return nil, fmt.Errorf("Invalid arg for ConcealedKong: %s", args[0])
		}
		return NewConcealedKongCommand(index), nil
	case Chow:
		if len(args) < 1 {
			return nil, fmt.Errorf("Not enough args for Chow")
		}
		index1, err := strconv.Atoi(args[0])
		if err != nil || index1 < 0 {
			return nil, fmt.Errorf("Invalid arg for Chow: %s", args[0])
		}
		index2, err := strconv.Atoi(args[1])
		if err != nil || index2 < 0 {
			return nil, fmt.Errorf("Invalid arg for Chow: %s", args[1])
		}
		if index1 == index2 {
			return nil, fmt.Errorf("Two different indices must be specified: %d", index1)
		}
		if index1 > index2 {
			index1, index2 = index2, index1
		}
		return NewChowCommand(index1, index2), nil
	case Pass:
		return NewPassCommand(), nil
	case Out:
		return NewOutCommand(), nil
	}
	return nil, fmt.Errorf("Unhandled command %s", cmdStr)
}
