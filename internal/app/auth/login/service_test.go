package login

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/art-es/yet-another-service/internal/app/auth/login/mock"
	"github.com/art-es/yet-another-service/internal/app/shared/dto"
	apperrors "github.com/art-es/yet-another-service/internal/app/shared/errors"
)

func TestService(t *testing.T) {
	type mocks struct {
		userRepository *mock.MockuserRepository
		hashChecker    *mock.MockhashChecker
		tokenGenerator *mock.MocktokenGenerator
	}

	for _, tt := range []struct {
		name   string
		setup  func(t *testing.T, m mocks)
		assert func(t *testing.T, res *dto.LoginOut, err error)
	}{
		{
			name: "find user by email in repository error",
			setup: func(t *testing.T, m mocks) {
				m.userRepository.EXPECT().
					FindByEmail(gomock.Any(), gomock.Eq("iivan@example.com")).
					Return(nil, errors.New("dummy error"))
			},
			assert: func(t *testing.T, res *dto.LoginOut, err error) {
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
			assert: func(t *testing.T, res *dto.LoginOut, err error) {
				assert.ErrorIs(t, err, apperrors.ErrUserNotFound)
				assert.Nil(t, res)
			},
		},
		{
			name: "check password by hash error",
			setup: func(t *testing.T, m mocks) {
				user := &dto.User{
					ID:           "dummy user id",
					DisplayName:  "Ivanov Ivan",
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
			assert: func(t *testing.T, res *dto.LoginOut, err error) {
				assert.EqualError(t, err, "check password by hash: dummy error")
				assert.Nil(t, res)
			},
		},
		{
			name: "wrong password",
			setup: func(t *testing.T, m mocks) {
				user := &dto.User{
					ID:           "dummy user id",
					DisplayName:  "Ivanov Ivan",
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
			assert: func(t *testing.T, res *dto.LoginOut, err error) {
				assert.ErrorIs(t, err, apperrors.ErrWrongPassword)
				assert.Nil(t, res)
			},
		},
		{
			name: "generate tokens error",
			setup: func(t *testing.T, m mocks) {
				user := &dto.User{
					ID:           "dummy user id",
					DisplayName:  "Ivanov Ivan",
					Email:        "iivan@example.com",
					PasswordHash: "dummy password hash",
				}

				m.userRepository.EXPECT().
					FindByEmail(gomock.Any(), gomock.Eq("iivan@example.com")).
					Return(user, nil)

				m.hashChecker.EXPECT().
					Check(gomock.Eq("secret123"), gomock.Eq("dummy password hash")).
					Return(nil)

				m.tokenGenerator.EXPECT().
					Generate(gomock.Eq("dummy user id")).
					Return(nil, errors.New("dummy error"))
			},
			assert: func(t *testing.T, res *dto.LoginOut, err error) {
				assert.EqualError(t, err, "generate tokens: dummy error")
				assert.Nil(t, res)
			},
		},
		{
			name: "ok",
			setup: func(t *testing.T, m mocks) {
				user := &dto.User{
					ID:           "dummy user id",
					DisplayName:  "Ivanov Ivan",
					Email:        "iivan@example.com",
					PasswordHash: "dummy password hash",
				}

				m.userRepository.EXPECT().
					FindByEmail(gomock.Any(), gomock.Eq("iivan@example.com")).
					Return(user, nil)

				m.hashChecker.EXPECT().
					Check(gomock.Eq("secret123"), gomock.Eq("dummy password hash")).
					Return(nil)

				tokensPair := &dto.AuthTokenPair{
					AccessToken:  "dummy access token",
					RefreshToken: "dummy refresh token",
				}

				m.tokenGenerator.EXPECT().
					Generate(gomock.Eq("dummy user id")).
					Return(tokensPair, nil)
			},
			assert: func(t *testing.T, res *dto.LoginOut, err error) {
				assert.NoError(t, err)
				expResult := &dto.LoginOut{
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
			res, err := service.Login(context.Background(), &dto.LoginIn{
				Email:    "iivan@example.com",
				Password: "secret123",
			})

			if tt.assert != nil {
				tt.assert(t, res, err)
			}
		})
	}
}
