package logout

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/art-es/yet-another-service/internal/app/auth"
	"github.com/art-es/yet-another-service/internal/app/auth/logout/mock"
	"github.com/art-es/yet-another-service/internal/driver/zerolog"
	"github.com/art-es/yet-another-service/internal/util/pointer"
)

func TestLogout(t *testing.T) {
	for _, tt := range []struct {
		name   string
		input  auth.LogoutIn
		setup  func(tokenService *mock.MocktokenService)
		assert func(t *testing.T, err error, logs []string)
	}{
		{
			name: "invalidate refresh token error",
			input: auth.LogoutIn{
				AccessToken:  pointer.To("access token"),
				RefreshToken: "refresh token",
			},
			setup: func(tokenService *mock.MocktokenService) {
				tokenService.EXPECT().
					Invalidate(gomock.Any(), "refresh token").
					Return(errors.New("dummy error"))
			},
			assert: func(t *testing.T, err error, logs []string) {
				assert.EqualError(t, err, "invalidate refresh token: dummy error")
				assert.Empty(t, logs)
			},
		},
		{
			name: "ok, invalidate access token error",
			input: auth.LogoutIn{
				AccessToken:  pointer.To("access token"),
				RefreshToken: "refresh token",
			},
			setup: func(tokenService *mock.MocktokenService) {
				tokenService.EXPECT().
					Invalidate(gomock.Any(), "refresh token").
					Return(nil)

				tokenService.EXPECT().
					Invalidate(gomock.Any(), "access token").
					Return(errors.New("dummy error"))
			},
			assert: func(t *testing.T, err error, logs []string) {
				assert.NoError(t, err)
				assert.Len(t, logs, 1)
				assert.Equal(t, `{"level":"warn","error":"dummy error","message":"invalidate acccess token error"}`, logs[0])
			},
		},
		{
			name: "ok, no access token",
			input: auth.LogoutIn{
				RefreshToken: "refresh token",
			},
			setup: func(tokenService *mock.MocktokenService) {
				tokenService.EXPECT().
					Invalidate(gomock.Any(), "refresh token").
					Return(nil)
			},
			assert: func(t *testing.T, err error, logs []string) {
				assert.NoError(t, err)
				assert.Empty(t, logs)
			},
		},
		{
			name: "ok",
			input: auth.LogoutIn{
				AccessToken:  pointer.To("access token"),
				RefreshToken: "refresh token",
			},
			setup: func(tokenService *mock.MocktokenService) {
				tokenService.EXPECT().
					Invalidate(gomock.Any(), "refresh token").
					Return(nil)

				tokenService.EXPECT().
					Invalidate(gomock.Any(), "access token").
					Return(nil)
			},
			assert: func(t *testing.T, err error, logs []string) {
				assert.NoError(t, err)
				assert.Empty(t, logs)
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logbuf := &bytes.Buffer{}
			logger := zerolog.NewLoggerWithWriter(logbuf)
			tokenService := mock.NewMocktokenService(ctrl)

			tt.setup(tokenService)

			service := NewService(tokenService, logger)
			err := service.Logout(context.Background(), &tt.input)

			var logs []string
			for s := bufio.NewScanner(logbuf); s.Scan(); {
				logs = append(logs, s.Text())
			}

			tt.assert(t, err, logs)
		})
	}
}
