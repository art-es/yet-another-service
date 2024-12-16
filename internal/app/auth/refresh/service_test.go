package refresh

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/art-es/yet-another-service/internal/app/auth"
	"github.com/art-es/yet-another-service/internal/app/auth/refresh/mock"
)

func TestService(t *testing.T) {
	type mocks struct {
		tokenService *mock.MocktokenService
	}

	now, _ := time.Parse(time.DateTime, "2000-01-01 10:00:00")
	nextHour, _ := time.Parse(time.DateTime, "2000-01-01 11:00:00")

	getCurrentTime = func() time.Time {
		return now
	}

	for _, tt := range []struct {
		name   string
		setup  func(t *testing.T, m mocks)
		assert func(t *testing.T, accessToken string, err error)
	}{
		{
			name: "parse refresh token error",
			setup: func(t *testing.T, m mocks) {
				m.tokenService.EXPECT().
					Parse(gomock.Eq("dummy refresh token")).
					Return(nil, errors.New("dummy error"))
			},
			assert: func(t *testing.T, accessToken string, err error) {
				assert.EqualError(t, err, "parse refresh token: dummy error")
			},
		},
		{
			name: "generate access token error",
			setup: func(t *testing.T, m mocks) {
				refreshTokenClaims := &auth.TokenClaims{
					UserID: "dummy user id",
				}

				m.tokenService.EXPECT().
					Parse(gomock.Eq("dummy refresh token")).
					Return(refreshTokenClaims, nil)

				expAccessTokenClaims := &auth.TokenClaims{
					IssuedAt:  now,
					ExpiresAt: nextHour,
					UserID:    "dummy user id",
				}

				m.tokenService.EXPECT().
					Generate(gomock.Eq(expAccessTokenClaims)).
					Return("", errors.New("dummy error"))
			},
			assert: func(t *testing.T, accessToken string, err error) {
				assert.EqualError(t, err, "generate access token: dummy error")
			},
		},
		{
			name: "ok",
			setup: func(t *testing.T, m mocks) {
				refreshTokenClaims := &auth.TokenClaims{
					UserID: "dummy user id",
				}

				m.tokenService.EXPECT().
					Parse(gomock.Eq("dummy refresh token")).
					Return(refreshTokenClaims, nil)

				expAccessTokenClaims := &auth.TokenClaims{
					IssuedAt:  now,
					ExpiresAt: nextHour,
					UserID:    "dummy user id",
				}

				m.tokenService.EXPECT().
					Generate(gomock.Eq(expAccessTokenClaims)).
					Return("dummy access token", nil)
			},
			assert: func(t *testing.T, accessToken string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "dummy access token", accessToken)
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks{
				tokenService: mock.NewMocktokenService(ctrl),
			}

			if tt.setup != nil {
				tt.setup(t, m)
			}

			service := NewService(m.tokenService)
			accessToken, err := service.Refresh("dummy refresh token")

			if tt.assert != nil {
				tt.assert(t, accessToken, err)
			}
		})
	}
}
