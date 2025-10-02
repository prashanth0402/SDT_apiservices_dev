package email

import (
	"log"
	"time"
)

type EmailSender struct{}

// SendEmail simulates sending an email (replace with SMTP/SendGrid/SES logic)
func (e EmailSender) SendEmail(email EmailInput) error {
	time.Sleep(10 * time.Millisecond) // simulate sending
	log.Printf("[EMAIL] From: %s, To: %v, CC: %v, BCC: %v, Subject: %s",
		email.F, email.T, email.C, email.BC, email.S)
	return nil
}
