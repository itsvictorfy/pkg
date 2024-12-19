package email

import (
	"testing"
)

func TestSendEmailWithSmtp(t *testing.T) {
	emailConf := EmailConfig{
		FromAddr: "someone",
		ToAddr:   "someone",
		SmtpHost: "smtp.gmail.com",
		Password: "PASSWORD", //https://support.google.com/mail/answer/185833?hl=en
	}
	message := []byte("Subject: Test mail\r\n\r\nEmail body\r\n")
	if err := emailConf.SendEmail(message); err != nil {
		t.Error(err.Error())
	}
}
