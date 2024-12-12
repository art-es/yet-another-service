package smtp

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"
)

var mailTemplate = template.Must(template.ParseFiles("mail_template.txt"))

type mailTemplateData struct {
	Subject string
	Content string
}

type Config struct {
	Host     string
	Port     int
	Identity string
	Username string
	Password string
}

type Service struct {
	address string
	auth    smtp.Auth
	from    string
}

func NewService(config Config) *Service {
	address := config.Host
	if config.Port != 0 {
		address = fmt.Sprintf("%s:%d", config.Host, config.Port)
	}

	return &Service{
		address: address,
		auth:    smtp.PlainAuth(config.Identity, config.Username, config.Password, config.Host),
		from:    config.Username,
	}
}

func (s *Service) MailTo(address, subject, content string) error {
	body := &bytes.Buffer{}
	data := mailTemplateData{
		Subject: subject,
		Content: content,
	}

	if err := mailTemplate.Execute(body, data); err != nil {
		return fmt.Errorf("execute mail template: %w", err)
	}

	if err := smtp.SendMail(s.address, s.auth, s.from, []string{address}, body.Bytes()); err != nil {
		return fmt.Errorf("send mail: %w", err)
	}

	return nil
}
