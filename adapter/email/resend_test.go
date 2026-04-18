package email

import "testing"

func TestNewResendEmailSenderUsesConfiguredFromEmail(t *testing.T) {
	t.Setenv("RESEND_API_KEY", "test-api-key")
	t.Setenv("RESEND_FROM_EMAIL", "support@helptrix.com")

	sender := NewResendEmailSender()

	if sender.fromEmail != "support@helptrix.com" {
		t.Fatalf("expected configured from email, got %q", sender.fromEmail)
	}
}

func TestNewResendEmailSenderDefaultsToHelptrixFromEmail(t *testing.T) {
	t.Setenv("RESEND_API_KEY", "test-api-key")
	t.Setenv("RESEND_FROM_EMAIL", "")

	sender := NewResendEmailSender()

	if sender.fromEmail != defaultFromEmail {
		t.Fatalf("expected default from email %q, got %q", defaultFromEmail, sender.fromEmail)
	}
}
