package mail

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
)

const userActivationSubject = "User activation"

var (
	//go:embed user_activation_template.html
	userActivationTemplateData []byte
	userActivationTemplate     = template.Must(template.New("").Parse(string(userActivationTemplateData)))
)

type UserActivationData struct {
	ActivationURL string
}

type UserActivationMailer struct {
	mailer mailer
}

func NewUserActivationMailer(mailer mailer) *UserActivationMailer {
	return &UserActivationMailer{
		mailer: mailer,
	}
}

func (s *UserActivationMailer) MailTo(address string, data UserActivationData) error {
	content := &bytes.Buffer{}
	if err := userActivationTemplate.Execute(content, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	if err := s.mailer.MailTo(address, userActivationSubject, content.String()); err != nil {
		return fmt.Errorf("mail: %w", err)
	}

	return nil
}
