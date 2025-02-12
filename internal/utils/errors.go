package utils

import "errors"

var (
	ErrMismatchHashPassword = errors.New("wrong password specified")
	ErrTooLongPassword      = errors.New("password is longer then 72 ASCII characters")
)
