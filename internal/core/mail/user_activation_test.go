package mail

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/art-es/yet-another-service/internal/app/shared/dto"
	"github.com/art-es/yet-another-service/internal/core/mail/mock"
)

func TestUserActivationMailer(t *testing.T) {
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

	for _, tt := range []struct {
		name   string
		setup  func(mailRepository *mock.MockmailRepository)
		assert func(t *testing.T, err error)
	}{
		{
			name: "mail error",
			setup: func(mailRepository *mock.MockmailRepository) {
				expectedMails := []dto.Mail{
					{
						Address: address,
						Subject: userActivationSubject,
						Content: content,
					},
				}

				mailRepository.EXPECT().
					Save(gomock.Any(), gomock.Eq(expectedMails)).
					Return(errors.New("dummy error"))
			},
			assert: func(t *testing.T, err error) {
				assert.EqualError(t, err, "save mail: dummy error")
			},
		},
		{
			name: "ok",
			setup: func(mailRepository *mock.MockmailRepository) {
				expectedMails := []dto.Mail{
					{
						Address: address,
						Subject: userActivationSubject,
						Content: content,
					},
				}

				mailRepository.EXPECT().
					Save(gomock.Any(), gomock.Eq(expectedMails)).
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

			mockMailRepository := mock.NewMockmailRepository(ctrl)

			tt.setup(mockMailRepository)

			err := NewUserActivationMailer(mockMailRepository).
				MailTo(context.Background(), address, UserActivationData{ActivationURL: activationURL})

			tt.assert(t, err)
		})
	}
}
