package signup

import (
	"bytes"
	"fmt"
	"net/url"
	"text/template"
)

const (
	activationMailSubject = "User activation"
)

var (
	activationMailTemplate = template.Must(template.ParseFiles("activation_mail.html"))
)

type activationMailData struct {
	ActivationURL string
}

func buildActivationMailContent(activationURL url.URL, token string) (string, error) {
	query := activationURL.Query()
	query.Set("token", token)
	activationURL.RawQuery = query.Encode()

	content := &bytes.Buffer{}
	contentData := activationMailData{
		ActivationURL: activationURL.String(),
	}

	if err := activationMailTemplate.Execute(content, contentData); err != nil {
		return "", fmt.Errorf("execute activation mail template: %w", err)
	}

	return content.String(), nil
}
