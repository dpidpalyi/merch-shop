package utils

import "errors"

var (
	ErrMismatchHashPassword = errors.New("user password mismatch")
	ErrTooLongPassword      = errors.New("password is longer then 72 ASCII charcters")
)
