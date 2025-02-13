package service

import "errors"

var (
	ErrNotEnoughCoins = errors.New("not enough coins")
	ErrSendToYourself = errors.New("can't send coins to yourself")
)
