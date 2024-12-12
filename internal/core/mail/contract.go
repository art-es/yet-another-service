//go:generate mockgen -source=contract.go -destination=mock/contract.go -package=mock
package mail

type mailer interface {
	MailTo(address, subject, content string) error
}
