package domain

import "errors"

var (
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidCommandArgs = errors.New("invalid command arguments")
	ErrNotCommand         = errors.New("not a command")
	ErrNotEnoughArguments = errors.New("not enough arguments provided")
	ErrTooManyArguments   = errors.New("too many arguments provided")

	ErrGroupDoesNotExist = errors.New("group does not exist")
)
