package repository

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrNotEnoughCoins = errors.New("not enough coins")
)
