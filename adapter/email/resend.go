package email

import (
	"errors"
	"os"

	"github.com/resend/resend-go/v2"
)

const defaultFromEmail = "no-reply@helptrix.com"

type ResendEmailSender struct {
	client    *resend.Client
	fromEmail string
}

func NewResendEmailSender() *ResendEmailSender {
	apiKey := os.Getenv("RESEND_API_KEY")
	fromEmail := os.Getenv("RESEND_FROM_EMAIL")
	if fromEmail == "" {
		fromEmail = defaultFromEmail
	}

	client := resend.NewClient(apiKey)
	return &ResendEmailSender{
		client:    client,
		fromEmail: fromEmail,
	}
}

func (s *ResendEmailSender) Send(to, subject, htmlBody string) error {
	params := &resend.SendEmailRequest{
		From:    s.fromEmail,
		To:      []string{to},
		Subject: subject,
		Html:    htmlBody,
	}
	_, err := s.client.Emails.Send(params)
	if err != nil {
		return errors.New("error sending email via resend: " + err.Error())
	}
	return nil
}
