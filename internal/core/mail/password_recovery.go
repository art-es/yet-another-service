package mail

import (
	"bytes"
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
	mailer mailer
}

func NewPasswordRecoveryMailer(mailer mailer) *PasswordRecoveryMailer {
	return &PasswordRecoveryMailer{
		mailer: mailer,
	}
}

func (s *PasswordRecoveryMailer) MailTo(address string, data PasswordRecoveryData) error {
	content := &bytes.Buffer{}
	if err := passwordRecoveryTemplate.Execute(content, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	if err := s.mailer.MailTo(address, passwordRecoverySubject, content.String()); err != nil {
		return fmt.Errorf("mail: %w", err)
	}

	return nil
}
