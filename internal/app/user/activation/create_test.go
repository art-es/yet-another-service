package activation

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/art-es/yet-another-service/internal/app/shared/models"
	"github.com/art-es/yet-another-service/internal/app/user/activation/mock"
	"github.com/art-es/yet-another-service/internal/core/mail"
	"github.com/art-es/yet-another-service/internal/core/transaction"
)

type createActivationMocks struct {
	activationRepository *mock.MockactivationRepository
	activationMailer     *mock.MockactivationMailer
}

func TestCreateActivation(t *testing.T) {
	for _, tt := range []struct {
		name   string
		setup  func(m createActivationMocks)
		assert func(t *testing.T, err error)
	}{
		{
			name: "save activation in repository error",
			setup: func(m createActivationMocks) {
				m.expectSaveActivation(errors.New("foo error"))
			},
			assert: func(t *testing.T, err error) {
				assert.EqualError(t, err, "save activation in repository: foo error")
			},
		},
		{
			name: "mail activation to user error",
			setup: func(m createActivationMocks) {
				m.expectSaveActivation(nil)
				m.expectMailActivation(errors.New("foo error"))
			},
			assert: func(t *testing.T, err error) {
				assert.EqualError(t, err, "mail activation to user: foo error")
			},
		},
		{
			name: "ok",
			setup: func(m createActivationMocks) {
				m.expectSaveActivation(nil)
				m.expectMailActivation(nil)
			},
			assert: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := newCreateActivationMocks(ctrl)
			tt.setup(m)

			baseAcivationURL, _ := url.Parse("http://localhost/activate?q=1")
			ctx := context.Background()
			tx := transaction.New(ctx)
			user := &models.User{ID: "user id", Email: "iivan@example.com"}

			service := NewService(*baseAcivationURL, m.activationRepository, nil, m.activationMailer)
			err := service.Create(ctx, tx, user)

			tt.assert(t, err)
		})
	}
}

func newCreateActivationMocks(ctrl *gomock.Controller) createActivationMocks {
	return createActivationMocks{
		activationRepository: mock.NewMockactivationRepository(ctrl),
		activationMailer:     mock.NewMockactivationMailer(ctrl),
	}
}

func (m createActivationMocks) expectSaveActivation(err error) {
	expectedActivation := &models.UserActivation{
		UserID: "user id",
	}

	m.activationRepository.EXPECT().
		Save(gomock.Any(), gomock.Any(), gomock.Eq(expectedActivation)).
		Do(func(_ context.Context, _ transaction.Transaction, activation *models.UserActivation) {
			activation.Token = "foo_token"
		}).
		Return(err)
}

func (m createActivationMocks) expectMailActivation(err error) {
	expectedMailData := mail.UserActivationData{
		ActivationURL: "http://localhost/activate?q=1&token=foo_token",
	}

	m.activationMailer.EXPECT().
		MailTo(gomock.Any(), gomock.Eq("iivan@example.com"), gomock.Eq(expectedMailData)).
		Return(err)
}
