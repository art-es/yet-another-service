package password_recovery

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/art-es/yet-another-service/internal/app/auth"
	"github.com/art-es/yet-another-service/internal/app/auth/password_recovery/mock"
	apperrors "github.com/art-es/yet-another-service/internal/app/shared/errors"
	"github.com/art-es/yet-another-service/internal/app/shared/models"
	"github.com/art-es/yet-another-service/internal/core/transaction"
)

type recoverState struct {
	txRollbacked bool
	txCommitted  bool
}

type recoverMocks struct {
	userRepository     *mock.MockuserRepository
	recoveryRepository *mock.MockrecoveryRepository
	hashService        *mock.MockhashService
	state              *recoverState
}

func TestRecover(t *testing.T) {
	for _, tt := range []struct {
		name   string
		setup  func(m recoverMocks)
		assert func(t *testing.T, err error, state recoverState)
	}{
		{
			name: "find recovery in repository error",
			setup: func(m recoverMocks) {
				m.expectFindRecovery(false, errors.New("foo error"))
			},
			assert: func(t *testing.T, err error, state recoverState) {
				assert.EqualError(t, err, "find recovery in repository: foo error")
				assert.False(t, state.txRollbacked)
				assert.False(t, state.txCommitted)
			},
		},
		{
			name: "recovery not found",
			setup: func(m recoverMocks) {
				m.expectFindRecovery(false, nil)
			},
			assert: func(t *testing.T, err error, state recoverState) {
				assert.ErrorIs(t, err, apperrors.ErrUserPasswordRecoveryNotFound)
				assert.False(t, state.txRollbacked)
				assert.False(t, state.txCommitted)
			},
		},
		{
			name: "find user in repository error",
			setup: func(m recoverMocks) {
				m.expectFindRecovery(true, nil)
				m.expectFindUser(false, errors.New("foo error"))
			},
			assert: func(t *testing.T, err error, state recoverState) {
				assert.EqualError(t, err, "find user in repository: foo error")
				assert.False(t, state.txRollbacked)
				assert.False(t, state.txCommitted)
			},
		},
		{
			name: "user not found",
			setup: func(m recoverMocks) {
				m.expectFindRecovery(true, nil)
				m.expectFindUser(false, nil)
			},
			assert: func(t *testing.T, err error, state recoverState) {
				assert.ErrorIs(t, err, apperrors.ErrUserNotFound)
				assert.False(t, state.txRollbacked)
				assert.False(t, state.txCommitted)
			},
		},
		{
			name: "check old password with hash error",
			setup: func(m recoverMocks) {
				m.expectFindRecovery(true, nil)
				m.expectFindUser(true, nil)
				m.expectCheckOldPasswordHash(errors.New("foo error"))
			},
			assert: func(t *testing.T, err error, state recoverState) {
				assert.EqualError(t, err, "check old password with hash: foo error")
				assert.False(t, state.txRollbacked)
				assert.False(t, state.txCommitted)
			},
		},
		{
			name: "old password and hash mismatch",
			setup: func(m recoverMocks) {
				m.expectFindRecovery(true, nil)
				m.expectFindUser(true, nil)
				m.expectCheckOldPasswordHash(apperrors.ErrHashMismatched)
			},
			assert: func(t *testing.T, err error, state recoverState) {
				assert.ErrorIs(t, err, apperrors.ErrHashMismatched)
				assert.False(t, state.txRollbacked)
				assert.False(t, state.txCommitted)
			},
		},
		{
			name: "save user in repository error",
			setup: func(m recoverMocks) {
				m.expectFindRecovery(true, nil)
				m.expectFindUser(true, nil)
				m.expectCheckOldPasswordHash(nil)
				m.expectGenerateNewPasswordHash(errors.New("foo error"))
			},
			assert: func(t *testing.T, err error, state recoverState) {
				assert.EqualError(t, err, "generate new password hash: foo error")
				assert.False(t, state.txRollbacked)
				assert.False(t, state.txCommitted)
			},
		},
		{
			name: "save user in repository error",
			setup: func(m recoverMocks) {
				m.expectFindRecovery(true, nil)
				m.expectFindUser(true, nil)
				m.expectCheckOldPasswordHash(nil)
				m.expectGenerateNewPasswordHash(nil)
				m.expectSaveUser(errors.New("foo error"), nil)
			},
			assert: func(t *testing.T, err error, state recoverState) {
				assert.EqualError(t, err, "save user in repository: foo error")
				assert.True(t, state.txRollbacked)
				assert.False(t, state.txCommitted)
			},
		},
		{
			name: "delete recovery in repository error",
			setup: func(m recoverMocks) {
				m.expectFindRecovery(true, nil)
				m.expectFindUser(true, nil)
				m.expectCheckOldPasswordHash(nil)
				m.expectGenerateNewPasswordHash(nil)
				m.expectSaveUser(nil, nil)
				m.expectDeleteRecovery(errors.New("foo error"))
			},
			assert: func(t *testing.T, err error, state recoverState) {
				assert.EqualError(t, err, "delete recovery in repository: foo error")
				assert.True(t, state.txRollbacked)
				assert.False(t, state.txCommitted)
			},
		},
		{
			name: "commit transaction error",
			setup: func(m recoverMocks) {
				m.expectFindRecovery(true, nil)
				m.expectFindUser(true, nil)
				m.expectCheckOldPasswordHash(nil)
				m.expectGenerateNewPasswordHash(nil)
				m.expectSaveUser(nil, errors.New("foo error"))
				m.expectDeleteRecovery(nil)
			},
			assert: func(t *testing.T, err error, state recoverState) {
				assert.EqualError(t, err, "commit transaction: foo error")
				assert.False(t, state.txRollbacked)
				assert.True(t, state.txCommitted)
			},
		},
		{
			name: "ok",
			setup: func(m recoverMocks) {
				m.expectFindRecovery(true, nil)
				m.expectFindUser(true, nil)
				m.expectCheckOldPasswordHash(nil)
				m.expectGenerateNewPasswordHash(nil)
				m.expectSaveUser(nil, nil)
				m.expectDeleteRecovery(nil)
			},
			assert: func(t *testing.T, err error, state recoverState) {
				assert.NoError(t, err)
				assert.False(t, state.txRollbacked)
				assert.True(t, state.txCommitted)
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := newRecoverMocks(ctrl)
			tt.setup(m)

			service := NewService(url.URL{}, m.userRepository, m.recoveryRepository, nil, m.hashService)
			err := service.Recover(context.Background(), &auth.PasswordRecoverIn{
				Token:       "foo_token",
				OldPassword: "old password",
				NewPassword: "new password",
			})

			tt.assert(t, err, *m.state)
		})
	}
}

func newRecoverMocks(ctrl *gomock.Controller) recoverMocks {
	return recoverMocks{
		userRepository:     mock.NewMockuserRepository(ctrl),
		recoveryRepository: mock.NewMockrecoveryRepository(ctrl),
		hashService:        mock.NewMockhashService(ctrl),
		state:              new(recoverState),
	}
}

func (m *recoverMocks) expectFindRecovery(found bool, err error) {
	var foundRecovery *models.PasswordRecovery
	if found {
		foundRecovery = &models.PasswordRecovery{
			Token:  "foo_token",
			UserID: "user id",
		}
	}

	m.recoveryRepository.EXPECT().
		Find(gomock.Any(), gomock.Eq("foo_token")).
		Return(foundRecovery, err)
}

func (m *recoverMocks) expectFindUser(found bool, err error) {
	var foundUser *models.User
	if found {
		foundUser = &models.User{
			ID:           "user id",
			Name:         "Ivanov Ivan",
			Email:        "iivan@example.com",
			PasswordHash: "old password hash",
		}
	}

	m.userRepository.EXPECT().
		Find(gomock.Any(), gomock.Eq("user id")).
		Return(foundUser, err)
}

func (m *recoverMocks) expectCheckOldPasswordHash(err error) {
	m.hashService.EXPECT().
		Check(gomock.Eq("old password"), gomock.Eq("old password hash")).
		Return(err)
}

func (m *recoverMocks) expectGenerateNewPasswordHash(err error) {
	var generatedHash string
	if err == nil {
		generatedHash = "new password hash"
	}

	m.hashService.EXPECT().
		Generate(gomock.Eq("new password")).
		Return(generatedHash, err)
}

func (m *recoverMocks) expectSaveUser(userSaveErr, txCommitErr error) {
	expectedUser := &models.User{
		ID:           "user id",
		Name:         "Ivanov Ivan",
		Email:        "iivan@example.com",
		PasswordHash: "new password hash",
	}

	m.userRepository.EXPECT().
		Save(gomock.Any(), gomock.Not(nil), gomock.Eq(expectedUser)).
		Do(func(_ context.Context, tx transaction.Transaction, u *models.User) {
			tx.AddRollback(func() {
				m.state.txRollbacked = true
			})

			tx.AddCommit(func() error {
				m.state.txCommitted = true
				return txCommitErr
			})
		}).
		Return(userSaveErr)
}

func (m *recoverMocks) expectDeleteRecovery(err error) {
	m.recoveryRepository.EXPECT().
		Delete(gomock.Any(), gomock.Not(nil), gomock.Eq("foo_token")).
		Return(err)
}
