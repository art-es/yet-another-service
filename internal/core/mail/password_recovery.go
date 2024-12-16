package mail

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"html/template"
)

const passwordRecoverySubject = "Password recovery"

var (
	//go:embed password_recovery_template.html
	passwordRecoveryTemplateData []byte
	passwordRecoveryTemplate     = template.Must(template.New("").Parse(string(passwordRecoveryTemplateData)))
)

type PasswordRecoveryData struct {
	RecoveryURL string
}

type PasswordRecoveryMailer struct {
	mailRepository mailRepository
}

func NewPasswordRecoveryMailer(mailRepository mailRepository) *PasswordRecoveryMailer {
	return &PasswordRecoveryMailer{
		mailRepository: mailRepository,
	}
}

func (s *PasswordRecoveryMailer) MailTo(ctx context.Context, address string, data PasswordRecoveryData) error {
	content := &bytes.Buffer{}
	if err := passwordRecoveryTemplate.Execute(content, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return saveMail(s.mailRepository, ctx, address, passwordRecoverySubject, content.String())
}
