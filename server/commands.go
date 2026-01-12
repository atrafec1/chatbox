package server

import (
	. "chatbox/domain"
	"strings"
)

type CommandHandler func(args []string, session *Session) error

type Command struct {
	Name        string
	Handler     CommandHandler
	Usage       string
	Description string
	MinArgs     int
}

var Commands = map[string]Command{}

func CommandParser(command string) (*Command, error) {
	// /join <group_name>
	if !strings.HasPrefix(command, "/") {
		return nil, ErrNotCommand
	}
	cmd_without_prefix := command[1:]
	parts := strings.Split(cmd_without_prefix, " ")

	cmdName := parts[0]
	arguments := parts[1:]

	cmd, exists := Commands[cmdName]
	if !exists {
		return nil, ErrNonExistentCommand
	}

	if cmd.MinArgs > len(arguments) {
		return nil, ErrInvalidCommandArgs
	}
	return &cmd, nil
}

func (s* Session) ExecuteCommand(cmd *Command) error {
}
