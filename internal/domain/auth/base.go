package auth

import (
	"errors"
)

var (
	ErrEmailAlreadyTaken  = errors.New("email address is already taken")
	ErrUserNotFound       = errors.New("user not found")
	ErrActivationNotFound = errors.New("activation not found")
	ErrWrongPassword      = errors.New("wrong password")
	ErrInvalidToken       = errors.New("invalid token")
)

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
}

type Activation struct {
	Token  string
	UserID string
}
