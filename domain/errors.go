package domain

import "errors"

var (
	ErrInvalidPassword    = errors.New("invalid password")
	ErrNonExistentCommand = errors.New("command does not exist")
	ErrInvalidCommandArgs = errors.New("invalid command arguments")
	ErrNotCommand         = errors.New("not a command")
)
