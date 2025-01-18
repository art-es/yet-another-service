package signup

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/art-es/yet-another-service/internal/app/auth/signup/mock"
	"github.com/art-es/yet-another-service/internal/app/shared/dto"
	apperrors "github.com/art-es/yet-another-service/internal/app/shared/errors"
	"github.com/art-es/yet-another-service/internal/core/transaction"
)

const (
	userID            = "dummy user id"
	userName          = "Ivanov Ivan"
	userEmail         = "iivan@example.com"
	userPassword      = "secret123"
	userPasswordHash  = "dummy user password hash"
	dummyErrorMessage = "dummy error"
)

func TestService(t *testing.T) {
	type mocks struct {
		hashGenerator     *mock.MockhashGenerator
		userRepository    *mock.MockuserRepository
		activationService *mock.MockactivationService
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
				assert.ErrorIs(t, err, apperrors.ErrEmailAlreadyTaken)
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

				expectedUser := &dto.User{
					DisplayName:  userName,
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

				expectedUser := &dto.User{
					DisplayName:  userName,
					Email:        userEmail,
					PasswordHash: userPasswordHash,
				}

				m.userRepository.EXPECT().
					Save(gomock.Any(), gomock.Not(nil), gomock.Eq(expectedUser)).
					Do(func(_ context.Context, tx transaction.Transaction, user *dto.User) {
						user.ID = userID
					}).
					Return(nil)

				expectedUser = &dto.User{
					ID:           userID,
					DisplayName:  userName,
					Email:        userEmail,
					PasswordHash: userPasswordHash,
				}

				m.activationService.EXPECT().
					Create(gomock.Any(), gomock.Not(nil), gomock.Eq(expectedUser)).
					Return(errors.New(dummyErrorMessage))
			},
			assert: func(t *testing.T, err error) {
				assert.Errorf(t, err, "create activation: %s", dummyErrorMessage)
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

				expectedUser := &dto.User{
					DisplayName:  userName,
					Email:        userEmail,
					PasswordHash: userPasswordHash,
				}

				m.userRepository.EXPECT().
					Save(gomock.Any(), gomock.Not(nil), gomock.Eq(expectedUser)).
					Do(func(_ context.Context, tx transaction.Transaction, user *dto.User) {
						user.ID = userID

						tx.AddCommit(func() error {
							return errors.New(dummyErrorMessage)
						})
					}).
					Return(nil)

				expectedUser = &dto.User{
					ID:           userID,
					DisplayName:  userName,
					Email:        userEmail,
					PasswordHash: userPasswordHash,
				}

				m.activationService.EXPECT().
					Create(gomock.Any(), gomock.Not(nil), gomock.Eq(expectedUser)).
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

				expectedUser := &dto.User{
					DisplayName:  userName,
					Email:        userEmail,
					PasswordHash: userPasswordHash,
				}

				m.userRepository.EXPECT().
					Save(gomock.Any(), gomock.Not(nil), gomock.Eq(expectedUser)).
					Do(func(_ context.Context, tx transaction.Transaction, user *dto.User) {
						user.ID = userID
					}).
					Return(nil)

				expectedUser = &dto.User{
					ID:           userID,
					DisplayName:  userName,
					Email:        userEmail,
					PasswordHash: userPasswordHash,
				}

				m.activationService.EXPECT().
					Create(gomock.Any(), gomock.Not(nil), gomock.Eq(expectedUser)).
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
				activationService: mock.NewMockactivationService(ctrl),
			}

			if tt.setup != nil {
				tt.setup(t, m)
			}

			service := NewService(
				m.hashGenerator,
				m.userRepository,
				m.activationService,
			)
			err := service.Signup(context.Background(), &dto.SignupIn{
				DisplayName: userName,
				Email:       userEmail,
				Password:    userPassword,
			})

			if tt.assert != nil {
				tt.assert(t, err)
			}
		})
	}
}
