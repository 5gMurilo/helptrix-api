package emailinterfaces

type IEmailSender interface {
	Send(to, subject, htmlBody string) error
}
