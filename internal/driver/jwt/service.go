package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"

	"github.com/art-es/yet-another-service/internal/domain/auth"
)

type Service struct {
	secret []byte
}

func NewService(secret string) *Service {
	return &Service{
		secret: []byte(secret),
	}
}

type internalClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"uid,omitempty"`
}

func (s *Service) Generate(claims *auth.TokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &internalClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(claims.IssuedAt),
			ExpiresAt: jwt.NewNumericDate(claims.ExpiresAt),
		},
		UserID: claims.UserID,
	})

	signedToken, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	return signedToken, nil
}

func (s *Service) Parse(signedToken string) (*auth.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(signedToken, &internalClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.secret, nil
	})
	if err != nil {
		return nil, auth.ErrInvalidToken
	}

	claims, ok := token.Claims.(*internalClaims)
	if !ok || !token.Valid {
		return nil, auth.ErrInvalidToken
	}

	return &auth.TokenClaims{
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
		UserID:    claims.UserID,
	}, nil
}
