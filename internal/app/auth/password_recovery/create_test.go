package password_recovery

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/art-es/yet-another-service/internal/app/auth/password_recovery/mock"
	apperrors "github.com/art-es/yet-another-service/internal/app/shared/errors"
	"github.com/art-es/yet-another-service/internal/app/shared/models"
	"github.com/art-es/yet-another-service/internal/core/mail"
	"github.com/art-es/yet-another-service/internal/core/transaction"
)

type createRecoveryState struct {
	txRollbacked bool
	txCommitted  bool
}

type createRecoveryMocks struct {
	userRepository     *mock.MockuserRepository
	recoveryRepository *mock.MockrecoveryRepository
	recoveryMailer     *mock.MockrecoveryMailer
	state              *createRecoveryState
}

func TestCreateRecovery(t *testing.T) {
	for _, tt := range []struct {
		name   string
		setup  func(m createRecoveryMocks)
		assert func(t *testing.T, err error, state createRecoveryState)
	}{
		{
			name: "find user by email in repository error",
			setup: func(m createRecoveryMocks) {
				m.expectFindUser(false, errors.New("foo error"))
			},
			assert: func(t *testing.T, err error, state createRecoveryState) {
				assert.EqualError(t, err, "find user by email in repository: foo error")
				assert.False(t, state.txRollbacked)
				assert.False(t, state.txCommitted)
			},
		},
		{
			name: "user not found",
			setup: func(m createRecoveryMocks) {
				m.expectFindUser(false, nil)
			},
			assert: func(t *testing.T, err error, state createRecoveryState) {
				assert.ErrorIs(t, err, apperrors.ErrUserNotFound)
				assert.False(t, state.txRollbacked)
				assert.False(t, state.txCommitted)
			},
		},
		{
			name: "save recovery in repository error",
			setup: func(m createRecoveryMocks) {
				m.expectFindUser(true, nil)
				m.expectSaveRecovery(errors.New("foo error"), nil)
			},
			assert: func(t *testing.T, err error, state createRecoveryState) {
				assert.EqualError(t, err, "save recovery in repository: foo error")
				assert.True(t, state.txRollbacked)
				assert.False(t, state.txCommitted)
			},
		},
		{
			name: "mail recovery to user error",
			setup: func(m createRecoveryMocks) {
				m.expectFindUser(true, nil)
				m.expectSaveRecovery(nil, nil)
				m.expectMailRecovery(errors.New("foo error"))
			},
			assert: func(t *testing.T, err error, state createRecoveryState) {
				assert.EqualError(t, err, "mail recovery to user: foo error")
				assert.True(t, state.txRollbacked)
				assert.False(t, state.txCommitted)
			},
		},
		{
			name: "commit transaction error",
			setup: func(m createRecoveryMocks) {
				m.expectFindUser(true, nil)
				m.expectSaveRecovery(nil, errors.New("foo error"))
				m.expectMailRecovery(nil)
			},
			assert: func(t *testing.T, err error, state createRecoveryState) {
				assert.EqualError(t, err, "commit transaction: foo error")
				assert.False(t, state.txRollbacked)
				assert.True(t, state.txCommitted)
			},
		},
		{
			name: "ok",
			setup: func(m createRecoveryMocks) {
				m.expectFindUser(true, nil)
				m.expectSaveRecovery(nil, nil)
				m.expectMailRecovery(nil)
			},
			assert: func(t *testing.T, err error, state createRecoveryState) {
				assert.NoError(t, err)
				assert.False(t, state.txRollbacked)
				assert.True(t, state.txCommitted)
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := newCreateRecoveryMocks(ctrl)
			tt.setup(m)

			baseRecoveryURL, _ := url.Parse("http://localhost/recover?some=foo")
			service := NewService(*baseRecoveryURL, m.userRepository, m.recoveryRepository, m.recoveryMailer, nil)
			err := service.CreateRecovery(context.Background(), "iivan@example.com")

			tt.assert(t, err, *m.state)
		})
	}
}

func newCreateRecoveryMocks(ctrl *gomock.Controller) createRecoveryMocks {
	return createRecoveryMocks{
		userRepository:     mock.NewMockuserRepository(ctrl),
		recoveryRepository: mock.NewMockrecoveryRepository(ctrl),
		recoveryMailer:     mock.NewMockrecoveryMailer(ctrl),
		state:              new(createRecoveryState),
	}
}

func (m *createRecoveryMocks) expectFindUser(found bool, err error) {
	var foundUser *models.User
	if found {
		foundUser = &models.User{
			ID:           "user id",
			Name:         "Ivanov Ivan",
			Email:        "iivan@example.com",
			PasswordHash: "password hash",
		}
	}

	m.userRepository.EXPECT().
		FindByEmail(gomock.Any(), gomock.Eq("iivan@example.com")).
		Return(foundUser, err)
}

func (m *createRecoveryMocks) expectSaveRecovery(recoverySaveErr, txCommitErr error) {
	expectedRecovery := &models.PasswordRecovery{
		UserID: "user id",
	}

	m.recoveryRepository.EXPECT().
		Save(gomock.Any(), gomock.Not(nil), gomock.Eq(expectedRecovery)).
		Do(func(_ context.Context, tx transaction.Transaction, r *models.PasswordRecovery) {
			r.Token = "foo_token"

			tx.AddRollback(func() {
				m.state.txRollbacked = true
			})

			tx.AddCommit(func() error {
				m.state.txCommitted = true
				return txCommitErr
			})
		}).
		Return(recoverySaveErr)
}

func (m *createRecoveryMocks) expectMailRecovery(err error) {
	expectedData := mail.PasswordRecoveryData{
		RecoveryURL: "http://localhost/recover?some=foo&token=foo_token",
	}

	m.recoveryMailer.EXPECT().
		MailTo(gomock.Any(), gomock.Eq("iivan@example.com"), gomock.Eq(expectedData)).
		Return(err)
}
