package errors

import "errors"

// Common
var (
	ErrUserNotFound           = errors.New("user not found")
	ErrUserActivationNotFound = errors.New("user activation not found")
)

// Auth specific
var (
	ErrInvalidAuthToken  = errors.New("invalid auth token")
	ErrEmailAlreadyTaken = errors.New("email address is already taken")
	ErrWrongPassword     = errors.New("wrong password")
)

// Hash specific
var (
	ErrHashMismatched = errors.New("mismatched hash and string")
)
