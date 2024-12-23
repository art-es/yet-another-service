package jwt

import (
	"errors"
	"fmt"

	"github.com/art-es/yet-another-service/internal/app/shared/dto"

	"github.com/golang-jwt/jwt/v4"

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

func (s *Service) Generate(claims *dto.AuthTokenClaims) (string, error) {
	tokenObject := jwt.NewWithClaims(jwt.SigningMethodHS256, &internalClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(claims.IssuedAt),
			ExpiresAt: jwt.NewNumericDate(claims.ExpiresAt),
		},
		UserID: claims.UserID,
	})

	signedToken, err := tokenObject.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	return signedToken, nil
}

func (s *Service) Parse(signedToken string) (*dto.AuthTokenClaims, error) {
	tokenObject, err := jwt.ParseWithClaims(signedToken, &internalClaims{}, func(token *jwt.Token) (any, error) {
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

	claims, ok := tokenObject.Claims.(*internalClaims)
	if !ok || !tokenObject.Valid {
		return nil, apperrors.ErrInvalidAuthToken
	}

	return &dto.AuthTokenClaims{
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
		UserID:    claims.UserID,
	}, nil
}
