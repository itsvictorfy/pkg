package email

import (
	"fmt"
	"log/slog"
	"net/smtp"
)

type EmailConfig struct {
	FromAddr string
	ToAddr   string
	Password string //https://support.google.com/mail/answer/185833?hl=en
	SmtpHost string
}

// SendEmail sends an email using the provided email configuration and message
func (email *EmailConfig) SendEmail(message []byte) error {
	auth := smtp.PlainAuth("", email.FromAddr, email.Password, email.SmtpHost)
	err := smtp.SendMail(email.SmtpHost+":587", auth, email.FromAddr, []string{email.ToAddr}, message)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("error recieved on SendEmail:%v", err)
	}
	slog.Info("Email Sent Successfully")
	return nil
}
