package users

import "errors"

var (
	ErrNotFound         = errors.New("user not found")
	ErrEmailTaken       = errors.New("email already taken")
	ErrInvalidEmail     = errors.New("invalid email")
	ErrUserIDIdentifier = errors.New("invalid user identifier")
	// ErrInvalidRequestFormat the universal error for all situations in which incorrect information was sent to the server
	ErrInvalidRequestFormat = errors.New("invalid request format")
)
