package main

import (
	"bytes"
	"html/template"

	"github.com/resendlabs/resend-go"
)

func (app *application) sendEmail(tmplPath string, data interface{}, recipient, subject string) error {
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	var emailContent bytes.Buffer
	err = tmpl.Execute(&emailContent, data)
	if err != nil {
		return err
	}

	mail := app.mailer.MailClient()

	emailResponseChannel := make(chan error)
	params := &resend.SendEmailRequest{
		From:    "Air BnB <info@air-bnb-cl.one>",
		To:      []string{recipient},
		Subject: subject,
		Html:    emailContent.String(),
	}

	go func() {
		_, err := mail.Emails.Send(params)
		if err != nil {
			emailResponseChannel <- err
		}
		close(emailResponseChannel)
	}()
	err, ok := <-emailResponseChannel
	if !ok {
		if err != nil {
			return err
		}
	}

	return nil
}
