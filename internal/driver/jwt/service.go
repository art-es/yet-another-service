package jwt

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"

	"github.com/art-es/yet-another-service/internal/app/auth"
	apperrors "github.com/art-es/yet-another-service/internal/app/shared/errors"
	"github.com/art-es/yet-another-service/internal/core/log"
)

type Service struct {
	secret []byte
	logger log.Logger
}

func NewService(secret string, logger log.Logger) *Service {
	return &Service{
		secret: []byte(secret),
		logger: logger.With().Str("package", "driver/jwt").Logger(),
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
			s.logger.Error().
				Str("signing_method", fmt.Sprintf("%v", token.Header["alg"])).
				Msg("unexpected signing method")

			return nil, errors.New("unexpected signing method")
		}

		return s.secret, nil
	})
	if err != nil {
		return nil, apperrors.ErrInvalidAuthToken
	}

	claims, ok := token.Claims.(*internalClaims)
	if !ok || !token.Valid {
		return nil, apperrors.ErrInvalidAuthToken
	}

	return &auth.TokenClaims{
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
		UserID:    claims.UserID,
	}, nil
}
