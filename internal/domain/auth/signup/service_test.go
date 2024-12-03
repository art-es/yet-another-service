package signup

import (
	"context"
	"errors"
	"net/url"
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
		mailSender        *mock.MockmailSender
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
					Return(nil, errors.New("dummy error"))
			},
			assert: func(t *testing.T, err error) {
				assert.EqualError(t, err, "create activation: dummy error")
			},
		},
		{
			name: "send activation mail error",
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

				activation := &auth.Activation{
					Token:  "dummy_token",
					UserID: "dummy user id",
				}

				m.activationCreator.EXPECT().
					Create(gomock.Any(), gomock.Not(nil), gomock.Eq("dummy user id")).
					Return(activation, nil)

				expActivationMailContent := "<!DOCTYPE html>\n<html>\n<body>\n    <p>To activate your account follow by link http://localhost:8080/foo?bar=1&token=dummy_token</p>\n</body>\n</html>"

				m.mailSender.EXPECT().
					SendMail(gomock.Eq("iivan@example.com"), gomock.Eq(activationMailSubject), gomock.Eq(expActivationMailContent)).
					Return(errors.New("dummy error"))
			},
			assert: func(t *testing.T, err error) {
				assert.EqualError(t, err, "send activation mail: dummy error")
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

				activation := &auth.Activation{
					Token:  "dummy_token",
					UserID: "dummy user id",
				}

				m.activationCreator.EXPECT().
					Create(gomock.Any(), gomock.Not(nil), gomock.Eq("dummy user id")).
					Return(activation, nil)

				expActivationMailContent := "<!DOCTYPE html>\n<html>\n<body>\n    <p>To activate your account follow by link http://localhost:8080/foo?bar=1&token=dummy_token</p>\n</body>\n</html>"

				m.mailSender.EXPECT().
					SendMail(gomock.Eq("iivan@example.com"), gomock.Eq(activationMailSubject), gomock.Eq(expActivationMailContent)).
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

				activation := &auth.Activation{
					Token:  "dummy_token",
					UserID: "dummy user id",
				}

				m.activationCreator.EXPECT().
					Create(gomock.Any(), gomock.Not(nil), gomock.Eq("dummy user id")).
					Return(activation, nil)

				expActivationMailContent := "<!DOCTYPE html>\n<html>\n<body>\n    <p>To activate your account follow by link http://localhost:8080/foo?bar=1&token=dummy_token</p>\n</body>\n</html>"

				m.mailSender.EXPECT().
					SendMail(gomock.Eq("iivan@example.com"), gomock.Eq(activationMailSubject), gomock.Eq(expActivationMailContent)).
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
				mailSender:        mock.NewMockmailSender(ctrl),
			}

			if tt.setup != nil {
				tt.setup(t, m)
			}

			activationURL, err := url.Parse("http://localhost:8080/foo?bar=1")
			assert.NoError(t, err)

			service := NewService(*activationURL, m.hashGenerator, m.userRepository, m.activationCreator, m.mailSender)
			err = service.Signup(context.Background(), &auth.SignupRequest{
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
