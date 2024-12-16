package mailing

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mockretrier "github.com/art-es/yet-another-service/internal/core/retrier/mock"
	"github.com/art-es/yet-another-service/internal/domain/mailing/mock"
	"github.com/art-es/yet-another-service/internal/domain/shared/models"
	"github.com/art-es/yet-another-service/internal/driver/zerolog"
)

func TestRun(t *testing.T) {
	for _, tt := range []struct {
		name   string
		setup  func(mailRepository *mock.MockmailRepository, mailer *mock.Mockmailer)
		assert func(t *testing.T, err error, logs []string)
	}{
		{
			name: "ok",
			setup: func(mailRepository *mock.MockmailRepository, mailer *mock.Mockmailer) {
				gotMails := []models.Mail{
					{ID: "mail id 1", Address: "foo@example.com", Subject: "mail 1", Content: "lorem ipsum 1"},
					{ID: "mail id 2", Address: "bar@example.com", Subject: "mail 2", Content: "lorem ipsum 2"},
				}
				mailRepository.EXPECT().
					Get(gomock.Any()).
					Return(gotMails, nil)

				mailer.EXPECT().
					MailTo(gomock.Eq("foo@example.com"), gomock.Eq("mail 1"), gomock.Eq("lorem ipsum 1")).
					Return(nil)

				mailer.EXPECT().
					MailTo(gomock.Eq("bar@example.com"), gomock.Eq("mail 2"), gomock.Eq("lorem ipsum 2")).
					Return(nil)

				expectedMails := []models.Mail{
					{Mailed: true, ID: "mail id 1", Address: "foo@example.com", Subject: "mail 1", Content: "lorem ipsum 1"},
					{Mailed: true, ID: "mail id 2", Address: "bar@example.com", Subject: "mail 2", Content: "lorem ipsum 2"},
				}
				mailRepository.EXPECT().
					Save(gomock.Any(), gomock.Eq(expectedMails)).
					Return(nil)

				// 2nd cycle
				mailRepository.EXPECT().
					Get(gomock.Any()).
					Return(nil, nil)
			},
			assert: func(t *testing.T, err error, logs []string) {
				assert.NoError(t, err)
				assert.Len(t, logs, 0)
			},
		},
		{
			name: "partially mails error",
			setup: func(mailRepository *mock.MockmailRepository, mailer *mock.Mockmailer) {
				gotMails := []models.Mail{
					{ID: "mail id 1", Address: "foo@example.com", Subject: "mail 1", Content: "lorem ipsum 1"},
					{ID: "mail id 2", Address: "bar@example.com", Subject: "mail 2", Content: "lorem ipsum 2"},
				}
				mailRepository.EXPECT().
					Get(gomock.Any()).
					Return(gotMails, nil)

				mailer.EXPECT().
					MailTo(gomock.Eq("foo@example.com"), gomock.Eq("mail 1"), gomock.Eq("lorem ipsum 1")).
					Return(errors.New("foo error"))

				mailer.EXPECT().
					MailTo(gomock.Eq("bar@example.com"), gomock.Eq("mail 2"), gomock.Eq("lorem ipsum 2")).
					Return(nil)

				expectedMails := []models.Mail{
					{Mailed: false, ID: "mail id 1", Address: "foo@example.com", Subject: "mail 1", Content: "lorem ipsum 1"},
					{Mailed: true, ID: "mail id 2", Address: "bar@example.com", Subject: "mail 2", Content: "lorem ipsum 2"},
				}
				mailRepository.EXPECT().
					Save(gomock.Any(), gomock.Eq(expectedMails)).
					Return(nil)

				// 2nd cycle
				mailRepository.EXPECT().
					Get(gomock.Any()).
					Return(nil, nil)
			},
			assert: func(t *testing.T, err error, logs []string) {
				assert.NoError(t, err)
				assert.Len(t, logs, 1)
				assert.Equal(t, `{"level":"error","error":"foo error","mail_id":"mail id 1","message":"mail error"}`, logs[0])
			},
		},
		{
			name: "all mails error",
			setup: func(mailRepository *mock.MockmailRepository, mailer *mock.Mockmailer) {
				gotMails := []models.Mail{
					{ID: "mail id 1", Address: "foo@example.com", Subject: "mail 1", Content: "lorem ipsum 1"},
					{ID: "mail id 2", Address: "bar@example.com", Subject: "mail 2", Content: "lorem ipsum 2"},
				}
				mailRepository.EXPECT().
					Get(gomock.Any()).
					Return(gotMails, nil)

				mailer.EXPECT().
					MailTo(gomock.Eq("foo@example.com"), gomock.Eq("mail 1"), gomock.Eq("lorem ipsum 1")).
					Return(errors.New("foo error"))

				mailer.EXPECT().
					MailTo(gomock.Eq("bar@example.com"), gomock.Eq("mail 2"), gomock.Eq("lorem ipsum 2")).
					Return(errors.New("bar error"))
			},
			assert: func(t *testing.T, err error, logs []string) {
				assert.EqualError(t, err, "process: all mails finished with error")
				assert.Len(t, logs, 2)
				assert.Equal(t, `{"level":"error","error":"foo error","mail_id":"mail id 1","message":"mail error"}`, logs[0])
				assert.Equal(t, `{"level":"error","error":"bar error","mail_id":"mail id 2","message":"mail error"}`, logs[1])
			},
		},
		{
			name: "get mails error",
			setup: func(mailRepository *mock.MockmailRepository, mailer *mock.Mockmailer) {
				mailRepository.EXPECT().
					Get(gomock.Any()).
					Return(nil, errors.New("dummy error"))
			},
			assert: func(t *testing.T, err error, logs []string) {
				assert.EqualError(t, err, "process: get mails from repository: dummy error")
				assert.Len(t, logs, 0)
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			retrier := mockretrier.NewMockRetrier(ctrl)
			retrier.EXPECT().
				Process(gomock.Any()).
				DoAndReturn(func(process func() error) error { return process() }).
				AnyTimes()

			mockMailRepository := mock.NewMockmailRepository(ctrl)
			mockMailer := mock.NewMockmailer(ctrl)
			tt.setup(mockMailRepository, mockMailer)

			config := Config{}
			logbuf := &bytes.Buffer{}
			logger := zerolog.NewLoggerWithWriter(logbuf)

			service := NewService(config, retrier, retrier, mockMailRepository, mockMailer, logger)
			err := service.Run(context.Background())

			var logs []string
			for s := bufio.NewScanner(logbuf); s.Scan(); {
				logs = append(logs, s.Text())
			}

			tt.assert(t, err, logs)
		})
	}
}
