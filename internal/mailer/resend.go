package mailer

import "github.com/resendlabs/resend-go"

type Mailer struct {
	ApiKey string
}

func NewMailer(apiKey string) Mailer {
	return Mailer{ApiKey: apiKey}
}

func (m *Mailer) MailClient() *resend.Client {
	return resend.NewClient(m.ApiKey)
}
