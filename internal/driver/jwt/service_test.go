package jwt

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"

	"github.com/art-es/yet-another-service/internal/app/auth"
	apperrors "github.com/art-es/yet-another-service/internal/app/shared/errors"
	"github.com/art-es/yet-another-service/internal/driver/zerolog"
)

const month = 30 * 24 * time.Hour

func TestService_Expired(t *testing.T) {
	service := &Service{secret: []byte("secret")}

	t.Run("unexpected signing method", func(t *testing.T) {
		logbuf := &bytes.Buffer{}
		service.logger = zerolog.NewLoggerWithWriter(logbuf)

		token := createUnexpectedMethodToken(t)

		claims, err := service.Parse(token)
		assert.ErrorIs(t, err, apperrors.ErrInvalidAuthToken)
		assert.Nil(t, claims)

		logs := getLogs(logbuf)
		assert.Len(t, logs, 1)
		assert.Equal(t, `{"level":"error","signing_method":"PS256","message":"unexpected signing method"}`, logs[0])
	})

	t.Run("expired", func(t *testing.T) {
		logbuf := &bytes.Buffer{}
		service.logger = zerolog.NewLoggerWithWriter(logbuf)
		prevMonth := time.Now().Add(-month)

		token, err := service.Generate(auth.NewAccessTokenClaims(prevMonth, "dummy user id"))
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := service.Parse(token)
		assert.ErrorIs(t, err, apperrors.ErrInvalidAuthToken)
		assert.Nil(t, claims)

		logs := getLogs(logbuf)
		assert.Empty(t, logs)
	})

	t.Run("ok", func(t *testing.T) {
		logbuf := &bytes.Buffer{}
		service.logger = zerolog.NewLoggerWithWriter(logbuf)

		token, err := service.Generate(auth.NewAccessTokenClaims(time.Now(), "dummy user id"))
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := service.Parse(token)
		assert.NoError(t, err)
		assert.Equal(t, "dummy user id", claims.UserID)

		logs := getLogs(logbuf)
		assert.Empty(t, logs)
	})
}

func getLogs(buf *bytes.Buffer) []string {
	var logs []string
	for s := bufio.NewScanner(buf); s.Scan(); {
		logs = append(logs, s.Text())
	}
	return logs
}

func createUnexpectedMethodToken(t *testing.T) string {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	assert.NoError(t, err)

	token := jwt.NewWithClaims(jwt.SigningMethodPS256, jwt.RegisteredClaims{})
	signedToken, err := token.SignedString(privateKey)
	assert.NoError(t, err)

	return signedToken
}
