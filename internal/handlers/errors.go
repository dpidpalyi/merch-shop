package handlers

import (
	"errors"
)

var (
	ErrEmptyNamePassword    = errors.New("empty name or password specified")
	ErrEmptyToUser          = errors.New("empty toUser field")
	ErrZeroOrNegativeAmount = errors.New("amount to send should be positive")
	ErrEmptyItem            = errors.New("empty item parameter")
)
