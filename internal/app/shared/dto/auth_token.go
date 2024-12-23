package dto

import "time"

const (
	accessTokenExpiry  = time.Hour * 1
	refreshTokenExpiry = time.Hour * 24 * 7
)

type AuthTokenPair struct {
	AccessToken  string
	RefreshToken string
}

type AuthTokenClaims struct {
	IssuedAt  time.Time
	ExpiresAt time.Time
	UserID    string
}

func (c *AuthTokenClaims) ToAccessTokenClaims(from time.Time) *AuthTokenClaims {
	return &AuthTokenClaims{
		IssuedAt:  from,
		ExpiresAt: from.Add(accessTokenExpiry),
		UserID:    c.UserID,
	}
}

func NewAccessTokenClaims(from time.Time, userID string) *AuthTokenClaims {
	return &AuthTokenClaims{
		IssuedAt:  from,
		ExpiresAt: from.Add(accessTokenExpiry),
		UserID:    userID,
	}
}

func NewRefreshTokenClaims(from time.Time, userID string) *AuthTokenClaims {
	return &AuthTokenClaims{
		IssuedAt:  from,
		ExpiresAt: from.Add(refreshTokenExpiry),
		UserID:    userID,
	}
}
