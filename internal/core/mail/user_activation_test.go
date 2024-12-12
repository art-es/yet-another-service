package mail

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/art-es/yet-another-service/internal/core/mail/mock"
)

const (
	address       = "foo@example.com"
	activationURL = `http://example.com/activate?token=foo`
	content       = `<!DOCTYPE html>
<html>
<body>
    <p>To activate your account follow by link http://example.com/activate?token=foo</p>
</body>
</html>`
)

func TestUserActivationMailer(t *testing.T) {
	for _, tt := range []struct {
		name   string
		setup  func(mailer *mock.MockMailer)
		assert func(t *testing.T, err error)
	}{
		{
			name: "mail error",
			setup: func(mailer *mock.MockMailer) {
				mailer.EXPECT().
					MailTo(gomock.Eq(address), gomock.Eq(userActivationSubject), gomock.Eq(content)).
					Return(errors.New("dummy error"))
			},
			assert: func(t *testing.T, err error) {
				assert.EqualError(t, err, "mail: dummy error")
			},
		},
		{
			name: "ok",
			setup: func(mailer *mock.MockMailer) {
				mailer.EXPECT().
					MailTo(gomock.Eq(address), gomock.Eq(userActivationSubject), gomock.Eq(content)).
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

			mockBaseMailer := mock.NewMockMailer(ctrl)

			tt.setup(mockBaseMailer)

			err := NewUserActivationMailer(mockBaseMailer).
				MailTo(address, UserActivationData{ActivationURL: activationURL})

			tt.assert(t, err)
		})
	}
}
