package signup

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/art-es/yet-another-service/internal/core/mail"
	"github.com/art-es/yet-another-service/internal/core/transaction"
	"github.com/art-es/yet-another-service/internal/domain/auth"
	"github.com/art-es/yet-another-service/internal/domain/auth/signup/mock"
	errorsd "github.com/art-es/yet-another-service/internal/domain/shared/errors"
	"github.com/art-es/yet-another-service/internal/domain/shared/models"
)

const (
	userID            = "dummy user id"
	userName          = "Ivanov Ivan"
	userEmail         = "iivan@example.com"
	userPassword      = "secret123"
	userPasswordHash  = "dummy user password hash"
	activationToken   = "dummy activation token"
	dummyErrorMessage = "dummy error"

	baseActivationURL = "http://localhost:8080/foo?bar=1"
	userActivationURL = "http://localhost:8080/foo?bar=1&token=dummy+activation+token"
)

func TestService(t *testing.T) {
	type mocks struct {
		hashGenerator        *mock.MockhashGenerator
		userRepository       *mock.MockuserRepository
		activationRepository *mock.MockactivationRepository
		activationMailer     *mock.MockactivationMailer
	}

	for _, tt := range []struct {
		name   string
		setup  func(t *testing.T, m mocks)
		assert func(t *testing.T, err error)
	}{
		{
			name: "check user exists in repository error",
			setup: func(t *testing.T, m mocks) {
				m.userRepository.EXPECT().
					Exists(gomock.Any(), gomock.Eq(userEmail)).
					Return(false, errors.New(dummyErrorMessage))
			},
			assert: func(t *testing.T, err error) {
				assert.Errorf(t, err, "check user exists in repository: %s", dummyErrorMessage)
			},
		},
		{
			name: "email address is already taken",
			setup: func(t *testing.T, m mocks) {
				m.userRepository.EXPECT().
					Exists(gomock.Any(), gomock.Eq(userEmail)).
					Return(true, nil)
			},
			assert: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errorsd.ErrEmailAlreadyTaken)
			},
		},
		{
			name: "generate password hash error",
			setup: func(t *testing.T, m mocks) {
				m.userRepository.EXPECT().
					Exists(gomock.Any(), gomock.Eq(userEmail)).
					Return(false, nil)

				m.hashGenerator.EXPECT().
					Generate(gomock.Eq(userPassword)).
					Return("", errors.New(dummyErrorMessage))
			},
			assert: func(t *testing.T, err error) {
				assert.Errorf(t, err, "generate password hash: %s", dummyErrorMessage)
			},
		},
		{
			name: "add user to repository error",
			setup: func(t *testing.T, m mocks) {
				m.userRepository.EXPECT().
					Exists(gomock.Any(), gomock.Eq(userEmail)).
					Return(false, nil)

				m.hashGenerator.EXPECT().
					Generate(gomock.Eq(userPassword)).
					Return(userPasswordHash, nil)

				expectedUser := &models.User{
					Name:         userName,
					Email:        userEmail,
					PasswordHash: userPasswordHash,
				}

				m.userRepository.EXPECT().
					Save(gomock.Any(), gomock.Not(nil), gomock.Eq(expectedUser)).
					Return(errors.New(dummyErrorMessage))
			},
			assert: func(t *testing.T, err error) {
				assert.Errorf(t, err, "add user to repository: %s", dummyErrorMessage)
			},
		},
		{
			name: "create activation error",
			setup: func(t *testing.T, m mocks) {
				m.userRepository.EXPECT().
					Exists(gomock.Any(), gomock.Eq(userEmail)).
					Return(false, nil)

				m.hashGenerator.EXPECT().
					Generate(gomock.Eq(userPassword)).
					Return(userPasswordHash, nil)

				expectedUser := &models.User{
					Name:         userName,
					Email:        userEmail,
					PasswordHash: userPasswordHash,
				}

				m.userRepository.EXPECT().
					Save(gomock.Any(), gomock.Not(nil), gomock.Eq(expectedUser)).
					Do(func(_ context.Context, tx transaction.Transaction, user *models.User) {
						user.ID = userID
					}).
					Return(nil)

				expectedActivation := &models.UserActivation{
					UserID: userID,
				}

				m.activationRepository.EXPECT().
					Save(gomock.Any(), gomock.Not(nil), gomock.Eq(expectedActivation)).
					Return(errors.New(dummyErrorMessage))
			},
			assert: func(t *testing.T, err error) {
				assert.Errorf(t, err, "create activation: %s", dummyErrorMessage)
			},
		},
		{
			name: "send activation mail error",
			setup: func(t *testing.T, m mocks) {
				m.userRepository.EXPECT().
					Exists(gomock.Any(), gomock.Eq(userEmail)).
					Return(false, nil)

				m.hashGenerator.EXPECT().
					Generate(gomock.Eq(userPassword)).
					Return(userPasswordHash, nil)

				expectedUser := &models.User{
					Name:         userName,
					Email:        userEmail,
					PasswordHash: userPasswordHash,
				}

				m.userRepository.EXPECT().
					Save(gomock.Any(), gomock.Not(nil), gomock.Eq(expectedUser)).
					Do(func(_ context.Context, tx transaction.Transaction, user *models.User) {
						user.ID = userID
					}).
					Return(nil)

				expectedActivation := &models.UserActivation{
					UserID: userID,
				}

				m.activationRepository.EXPECT().
					Save(gomock.Any(), gomock.Not(nil), gomock.Eq(expectedActivation)).
					Do(func(_ context.Context, _ transaction.Transaction, activation *models.UserActivation) {
						activation.Token = activationToken
					}).
					Return(nil)

				expectedActivationData := mail.UserActivationData{
					ActivationURL: userActivationURL,
				}

				m.activationMailer.EXPECT().
					MailTo(gomock.Eq(userEmail), gomock.Eq(expectedActivationData)).
					Return(errors.New(dummyErrorMessage))
			},
			assert: func(t *testing.T, err error) {
				assert.Errorf(t, err, "send activation mail: %s", dummyErrorMessage)
			},
		},
		{
			name: "commit transaction error",
			setup: func(t *testing.T, m mocks) {
				m.userRepository.EXPECT().
					Exists(gomock.Any(), gomock.Eq(userEmail)).
					Return(false, nil)

				m.hashGenerator.EXPECT().
					Generate(gomock.Eq(userPassword)).
					Return(userPasswordHash, nil)

				expectedUser := &models.User{
					Name:         userName,
					Email:        userEmail,
					PasswordHash: userPasswordHash,
				}

				m.userRepository.EXPECT().
					Save(gomock.Any(), gomock.Not(nil), gomock.Eq(expectedUser)).
					Do(func(_ context.Context, tx transaction.Transaction, user *models.User) {
						user.ID = userID

						tx.AddCommit(func() error {
							return errors.New(dummyErrorMessage)
						})
					}).
					Return(nil)

				expectedActivation := &models.UserActivation{
					UserID: userID,
				}

				m.activationRepository.EXPECT().
					Save(gomock.Any(), gomock.Not(nil), gomock.Eq(expectedActivation)).
					Do(func(_ context.Context, _ transaction.Transaction, activation *models.UserActivation) {
						activation.Token = activationToken
					}).
					Return(nil)

				expectedActivationData := mail.UserActivationData{
					ActivationURL: userActivationURL,
				}

				m.activationMailer.EXPECT().
					MailTo(gomock.Eq(userEmail), gomock.Eq(expectedActivationData)).
					Return(nil)
			},
			assert: func(t *testing.T, err error) {
				assert.Errorf(t, err, "commit transaction: %s", dummyErrorMessage)
			},
		},
		{
			name: "ok",
			setup: func(t *testing.T, m mocks) {
				m.userRepository.EXPECT().
					Exists(gomock.Any(), gomock.Eq(userEmail)).
					Return(false, nil)

				m.hashGenerator.EXPECT().
					Generate(gomock.Eq(userPassword)).
					Return(userPasswordHash, nil)

				expectedUser := &models.User{
					Name:         userName,
					Email:        userEmail,
					PasswordHash: userPasswordHash,
				}

				m.userRepository.EXPECT().
					Save(gomock.Any(), gomock.Not(nil), gomock.Eq(expectedUser)).
					Do(func(_ context.Context, tx transaction.Transaction, user *models.User) {
						user.ID = userID
					}).
					Return(nil)

				expectedActivation := &models.UserActivation{
					UserID: userID,
				}

				m.activationRepository.EXPECT().
					Save(gomock.Any(), gomock.Not(nil), gomock.Eq(expectedActivation)).
					Do(func(_ context.Context, _ transaction.Transaction, activation *models.UserActivation) {
						activation.Token = activationToken
					}).
					Return(nil)

				expectedActivationData := mail.UserActivationData{
					ActivationURL: userActivationURL,
				}

				m.activationMailer.EXPECT().
					MailTo(gomock.Eq(userEmail), gomock.Eq(expectedActivationData)).
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
				hashGenerator:        mock.NewMockhashGenerator(ctrl),
				userRepository:       mock.NewMockuserRepository(ctrl),
				activationRepository: mock.NewMockactivationRepository(ctrl),
				activationMailer:     mock.NewMockactivationMailer(ctrl),
			}

			if tt.setup != nil {
				tt.setup(t, m)
			}

			parsedBaseActivationURL, err := url.Parse(baseActivationURL)
			assert.NoError(t, err)

			service := NewService(
				*parsedBaseActivationURL,
				m.hashGenerator,
				m.userRepository,
				m.activationRepository,
				m.activationMailer,
			)
			err = service.Signup(context.Background(), &auth.SignupIn{
				Name:     userName,
				Email:    userEmail,
				Password: userPassword,
			})

			if tt.assert != nil {
				tt.assert(t, err)
			}
		})
	}
}
