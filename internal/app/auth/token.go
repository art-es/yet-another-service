package auth

import "time"

const (
	accessTokenExpiry  = time.Hour * 1
	refreshTokenExpiry = time.Hour * 24 * 7
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type TokenClaims struct {
	IssuedAt  time.Time
	ExpiresAt time.Time
	UserID    string
}

func (c *TokenClaims) ToAccessToken(from time.Time) *TokenClaims {
	return &TokenClaims{
		IssuedAt:  from,
		ExpiresAt: from.Add(accessTokenExpiry),
		UserID:    c.UserID,
	}
}

func NewAccessTokenClaims(from time.Time, userID string) *TokenClaims {
	return &TokenClaims{
		IssuedAt:  from,
		ExpiresAt: from.Add(accessTokenExpiry),
		UserID:    userID,
	}
}

func NewRefreshTokenClaims(from time.Time, userID string) *TokenClaims {
	return &TokenClaims{
		IssuedAt:  from,
		ExpiresAt: from.Add(refreshTokenExpiry),
		UserID:    userID,
	}
}
