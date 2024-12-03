package signup

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/art-es/yet-another-service/internal/core/transaction"
	"github.com/art-es/yet-another-service/internal/domain/auth"
	"github.com/art-es/yet-another-service/internal/domain/auth/signup/mock"
)

func TestService(t *testing.T) {
	type mocks struct {
		hashGenerator     *mock.MockhashGenerator
		userRepository    *mock.MockuserRepository
		activationCreator *mock.MockactivationCreator
	}

	for _, tt := range []struct {
		name   string
		setup  func(t *testing.T, m mocks)
		assert func(t *testing.T, err error)
	}{
		{
			name: "check user email exists in repository error",
			setup: func(t *testing.T, m mocks) {
				m.userRepository.EXPECT().
					EmailExists(gomock.Any(), gomock.Eq("iivan@example.com")).
					Return(false, errors.New("dummy error"))
			},
			assert: func(t *testing.T, err error) {
				assert.EqualError(t, err, "check user email exists in repository: dummy error")
			},
		},
		{
			name: "email address is already taken",
			setup: func(t *testing.T, m mocks) {
				m.userRepository.EXPECT().
					EmailExists(gomock.Any(), gomock.Eq("iivan@example.com")).
					Return(true, nil)
			},
			assert: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, auth.ErrEmailAlreadyTaken)
			},
		},
		{
			name: "generate password hash error",
			setup: func(t *testing.T, m mocks) {
				m.userRepository.EXPECT().
					EmailExists(gomock.Any(), gomock.Eq("iivan@example.com")).
					Return(false, nil)

				m.hashGenerator.EXPECT().
					Generate(gomock.Eq("secret123")).
					Return("", errors.New("dummy error"))
			},
			assert: func(t *testing.T, err error) {
				assert.EqualError(t, err, "generate password hash: dummy error")
			},
		},
		{
			name: "add user to repository error",
			setup: func(t *testing.T, m mocks) {
				m.userRepository.EXPECT().
					EmailExists(gomock.Any(), gomock.Eq("iivan@example.com")).
					Return(false, nil)

				m.hashGenerator.EXPECT().
					Generate(gomock.Eq("secret123")).
					Return("dummy password hash", nil)

				expUser := &auth.User{
					Name:         "Ivanov Ivan",
					Email:        "iivan@example.com",
					PasswordHash: "dummy password hash",
				}

				m.userRepository.EXPECT().
					Add(gomock.Any(), gomock.Not(nil), gomock.Eq(expUser)).
					Return(errors.New("dummy error"))
			},
			assert: func(t *testing.T, err error) {
				assert.EqualError(t, err, "add user to repository: dummy error")
			},
		},
		{
			name: "create activation error",
			setup: func(t *testing.T, m mocks) {
				m.userRepository.EXPECT().
					EmailExists(gomock.Any(), gomock.Eq("iivan@example.com")).
					Return(false, nil)

				m.hashGenerator.EXPECT().
					Generate(gomock.Eq("secret123")).
					Return("dummy password hash", nil)

				expUser := &auth.User{
					Name:         "Ivanov Ivan",
					Email:        "iivan@example.com",
					PasswordHash: "dummy password hash",
				}

				m.userRepository.EXPECT().
					Add(gomock.Any(), gomock.Not(nil), gomock.Eq(expUser)).
					Do(func(_ context.Context, tx transaction.Transaction, user *auth.User) {
						user.ID = "dummy user id"
					}).
					Return(nil)

				m.activationCreator.EXPECT().
					Create(gomock.Any(), gomock.Not(nil), gomock.Eq("dummy user id")).
					Return(errors.New("dummy error"))
			},
			assert: func(t *testing.T, err error) {
				assert.EqualError(t, err, "create activation: dummy error")
			},
		},
		{
			name: "commit transaction error",
			setup: func(t *testing.T, m mocks) {
				m.userRepository.EXPECT().
					EmailExists(gomock.Any(), gomock.Eq("iivan@example.com")).
					Return(false, nil)

				m.hashGenerator.EXPECT().
					Generate(gomock.Eq("secret123")).
					Return("dummy password hash", nil)

				expUser := &auth.User{
					Name:         "Ivanov Ivan",
					Email:        "iivan@example.com",
					PasswordHash: "dummy password hash",
				}

				m.userRepository.EXPECT().
					Add(gomock.Any(), gomock.Not(nil), gomock.Eq(expUser)).
					Do(func(_ context.Context, tx transaction.Transaction, user *auth.User) {
						user.ID = "dummy user id"

						tx.AddCommit(func() error {
							return errors.New("dummy error")
						})
					}).
					Return(nil)

				m.activationCreator.EXPECT().
					Create(gomock.Any(), gomock.Not(nil), gomock.Eq("dummy user id")).
					Return(nil)

			},
			assert: func(t *testing.T, err error) {
				assert.EqualError(t, err, "commit transaction: dummy error")
			},
		},
		{
			name: "ok",
			setup: func(t *testing.T, m mocks) {
				m.userRepository.EXPECT().
					EmailExists(gomock.Any(), gomock.Eq("iivan@example.com")).
					Return(false, nil)

				m.hashGenerator.EXPECT().
					Generate(gomock.Eq("secret123")).
					Return("dummy password hash", nil)

				expUser := &auth.User{
					Name:         "Ivanov Ivan",
					Email:        "iivan@example.com",
					PasswordHash: "dummy password hash",
				}

				m.userRepository.EXPECT().
					Add(gomock.Any(), gomock.Not(nil), gomock.Eq(expUser)).
					Do(func(_ context.Context, tx transaction.Transaction, user *auth.User) {
						user.ID = "dummy user id"
					}).
					Return(nil)

				m.activationCreator.EXPECT().
					Create(gomock.Any(), gomock.Not(nil), gomock.Eq("dummy user id")).
					Return(nil)
			},
			assert: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks{
				hashGenerator:     mock.NewMockhashGenerator(ctrl),
				userRepository:    mock.NewMockuserRepository(ctrl),
				activationCreator: mock.NewMockactivationCreator(ctrl),
			}

			if tt.setup != nil {
				tt.setup(t, m)
			}

			service := NewService(m.hashGenerator, m.userRepository, m.activationCreator)
			err := service.Signup(context.Background(), &auth.SignupRequest{
				Name:     "Ivanov Ivan",
				Email:    "iivan@example.com",
				Password: "secret123",
			})

			if tt.assert != nil {
				tt.assert(t, err)
			}
		})
	}
}