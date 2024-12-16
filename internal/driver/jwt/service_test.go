package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/art-es/yet-another-service/internal/app/auth"
	apperrors "github.com/art-es/yet-another-service/internal/app/shared/errors"
)

const month = 30 * 24 * time.Hour

func TestService_Expired(t *testing.T) {
	service := NewService("secret")

	t.Run("expired", func(t *testing.T) {
		prevMonth := time.Now().Add(-month)

		token, err := service.Generate(auth.NewAccessTokenClaims(prevMonth, "dummy user id"))
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := service.Parse(token)
		assert.ErrorIs(t, err, apperrors.ErrInvalidAuthToken)
		assert.Nil(t, claims)
	})

	t.Run("ok", func(t *testing.T) {
		token, err := service.Generate(auth.NewAccessTokenClaims(time.Now(), "dummy user id"))
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := service.Parse(token)
		assert.NoError(t, err)
		assert.Equal(t, "dummy user id", claims.UserID)
	})
}
