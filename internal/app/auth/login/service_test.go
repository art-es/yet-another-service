package login

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/art-es/yet-another-service/internal/app/auth"
	"github.com/art-es/yet-another-service/internal/app/auth/login/mock"
	apperrors "github.com/art-es/yet-another-service/internal/app/shared/errors"
	"github.com/art-es/yet-another-service/internal/app/shared/models"
)

func TestService(t *testing.T) {
	type mocks struct {
		userRepository *mock.MockuserRepository
		hashChecker    *mock.MockhashChecker
		tokenGenerator *mock.MocktokenGenerator
	}

	now, _ := time.Parse(time.DateTime, "2000-01-01 10:00:00")
	nextHour, _ := time.Parse(time.DateTime, "2000-01-01 11:00:00")
	nextWeek, _ := time.Parse(time.DateTime, "2000-01-08 10:00:00")

	getCurrentTime = func() time.Time {
		return now
	}

	for _, tt := range []struct {
		name   string
		setup  func(t *testing.T, m mocks)
		assert func(t *testing.T, res *auth.LoginOut, err error)
	}{
		{
			name: "find user by email in repository error",
			setup: func(t *testing.T, m mocks) {
				m.userRepository.EXPECT().
					FindByEmail(gomock.Any(), gomock.Eq("iivan@example.com")).
					Return(nil, errors.New("dummy error"))
			},
			assert: func(t *testing.T, res *auth.LoginOut, err error) {
				assert.EqualError(t, err, "find user by email in repository: dummy error")
				assert.Nil(t, res)
			},
		},
		{
			name: "user not found",
			setup: func(t *testing.T, m mocks) {
				m.userRepository.EXPECT().
					FindByEmail(gomock.Any(), gomock.Eq("iivan@example.com")).
					Return(nil, nil)
			},
			assert: func(t *testing.T, res *auth.LoginOut, err error) {
				assert.ErrorIs(t, err, apperrors.ErrUserNotFound)
				assert.Nil(t, res)
			},
		},
		{
			name: "check password by hash error",
			setup: func(t *testing.T, m mocks) {
				user := &models.User{
					ID:           "dummy user id",
					Name:         "Ivanov Ivan",
					Email:        "iivan@example.com",
					PasswordHash: "dummy password hash",
				}

				m.userRepository.EXPECT().
					FindByEmail(gomock.Any(), gomock.Eq("iivan@example.com")).
					Return(user, nil)

				m.hashChecker.EXPECT().
					Check(gomock.Eq("secret123"), gomock.Eq("dummy password hash")).
					Return(errors.New("dummy error"))
			},
			assert: func(t *testing.T, res *auth.LoginOut, err error) {
				assert.EqualError(t, err, "check password by hash: dummy error")
				assert.Nil(t, res)
			},
		},
		{
			name: "wrong password",
			setup: func(t *testing.T, m mocks) {
				user := &models.User{
					ID:           "dummy user id",
					Name:         "Ivanov Ivan",
					Email:        "iivan@example.com",
					PasswordHash: "dummy password hash",
				}

				m.userRepository.EXPECT().
					FindByEmail(gomock.Any(), gomock.Eq("iivan@example.com")).
					Return(user, nil)

				m.hashChecker.EXPECT().
					Check(gomock.Eq("secret123"), gomock.Eq("dummy password hash")).
					Return(apperrors.ErrHashMismatched)
			},
			assert: func(t *testing.T, res *auth.LoginOut, err error) {
				assert.ErrorIs(t, err, apperrors.ErrWrongPassword)
				assert.Nil(t, res)
			},
		},
		{
			name: "generate access token error",
			setup: func(t *testing.T, m mocks) {
				user := &models.User{
					ID:           "dummy user id",
					Name:         "Ivanov Ivan",
					Email:        "iivan@example.com",
					PasswordHash: "dummy password hash",
				}

				m.userRepository.EXPECT().
					FindByEmail(gomock.Any(), gomock.Eq("iivan@example.com")).
					Return(user, nil)

				m.hashChecker.EXPECT().
					Check(gomock.Eq("secret123"), gomock.Eq("dummy password hash")).
					Return(nil)

				expAccessTokenClaims := &auth.TokenClaims{
					IssuedAt:  now,
					ExpiresAt: nextHour,
					UserID:    "dummy user id",
				}

				m.tokenGenerator.EXPECT().
					Generate(gomock.Eq(expAccessTokenClaims)).
					Return("", errors.New("dummy error"))
			},
			assert: func(t *testing.T, res *auth.LoginOut, err error) {
				assert.EqualError(t, err, "generate access token: dummy error")
				assert.Nil(t, res)
			},
		},
		{
			name: "generate refresh token error",
			setup: func(t *testing.T, m mocks) {
				user := &models.User{
					ID:           "dummy user id",
					Name:         "Ivanov Ivan",
					Email:        "iivan@example.com",
					PasswordHash: "dummy password hash",
				}

				m.userRepository.EXPECT().
					FindByEmail(gomock.Any(), gomock.Eq("iivan@example.com")).
					Return(user, nil)

				m.hashChecker.EXPECT().
					Check(gomock.Eq("secret123"), gomock.Eq("dummy password hash")).
					Return(nil)

				expAccessTokenClaims := &auth.TokenClaims{
					IssuedAt:  now,
					ExpiresAt: nextHour,
					UserID:    "dummy user id",
				}

				m.tokenGenerator.EXPECT().
					Generate(gomock.Eq(expAccessTokenClaims)).
					Return("dummy access token", nil)

				expRefreshTokenClaims := &auth.TokenClaims{
					IssuedAt:  now,
					ExpiresAt: nextWeek,
					UserID:    "dummy user id",
				}

				m.tokenGenerator.EXPECT().
					Generate(gomock.Eq(expRefreshTokenClaims)).
					Return("", errors.New("dummy error"))
			},
			assert: func(t *testing.T, res *auth.LoginOut, err error) {
				assert.EqualError(t, err, "generate refresh token: dummy error")
				assert.Nil(t, res)
			},
		},
		{
			name: "ok",
			setup: func(t *testing.T, m mocks) {
				user := &models.User{
					ID:           "dummy user id",
					Name:         "Ivanov Ivan",
					Email:        "iivan@example.com",
					PasswordHash: "dummy password hash",
				}

				m.userRepository.EXPECT().
					FindByEmail(gomock.Any(), gomock.Eq("iivan@example.com")).
					Return(user, nil)

				m.hashChecker.EXPECT().
					Check(gomock.Eq("secret123"), gomock.Eq("dummy password hash")).
					Return(nil)

				expAccessTokenClaims := &auth.TokenClaims{
					IssuedAt:  now,
					ExpiresAt: nextHour,
					UserID:    "dummy user id",
				}

				m.tokenGenerator.EXPECT().
					Generate(gomock.Eq(expAccessTokenClaims)).
					Return("dummy access token", nil)

				expRefreshTokenClaims := &auth.TokenClaims{
					IssuedAt:  now,
					ExpiresAt: nextWeek,
					UserID:    "dummy user id",
				}

				m.tokenGenerator.EXPECT().
					Generate(gomock.Eq(expRefreshTokenClaims)).
					Return("dummy refresh token", nil)
			},
			assert: func(t *testing.T, res *auth.LoginOut, err error) {
				assert.NoError(t, err)
				expResult := &auth.LoginOut{
					AccessToken:  "dummy access token",
					RefreshToken: "dummy refresh token",
				}
				assert.Equal(t, expResult, res)
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks{
				userRepository: mock.NewMockuserRepository(ctrl),
				hashChecker:    mock.NewMockhashChecker(ctrl),
				tokenGenerator: mock.NewMocktokenGenerator(ctrl),
			}

			if tt.setup != nil {
				tt.setup(t, m)
			}

			service := NewService(m.userRepository, m.hashChecker, m.tokenGenerator)
			res, err := service.Login(context.Background(), &auth.LoginIn{
				Email:    "iivan@example.com",
				Password: "secret123",
			})

			if tt.assert != nil {
				tt.assert(t, res, err)
			}
		})
	}
}
