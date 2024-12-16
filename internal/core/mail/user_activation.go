package mail

import (
	"bytes"
	"context"
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
	mailRepository mailRepository
}

func NewUserActivationMailer(mailRepository mailRepository) *UserActivationMailer {
	return &UserActivationMailer{
		mailRepository: mailRepository,
	}
}

func (s *UserActivationMailer) MailTo(ctx context.Context, address string, data UserActivationData) error {
	content := &bytes.Buffer{}
	if err := userActivationTemplate.Execute(content, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return saveMail(s.mailRepository, ctx, address, userActivationSubject, content.String())
}
