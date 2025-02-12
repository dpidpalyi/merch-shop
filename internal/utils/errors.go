package utils

import "errors"

var (
	ErrMismatchHashPassword       = errors.New("wrong password specified")
	ErrTooLongPassword            = errors.New("password is longer then 72 ASCII characters")
	ErrNoAuthorizationHeader      = errors.New("missing Authorization header")
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")
	ErrInvalidSigningMethod       = errors.New("invalid signing method")
	ErrInvalidToken               = errors.New("invalid token")
	ErrInvalidClaims              = errors.New("invalid claims")
	ErrExpiredToken               = errors.New("expired token")
	ErrInvalidUserID              = errors.New("invalid user ID")
)
