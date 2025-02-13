package service

import "errors"

var (
	ErrSendToYourself = errors.New("can't send coins to yourself")
)
