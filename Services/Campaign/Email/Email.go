package email

import (
	"crypto/tls"
	"log"

	"gopkg.in/mail.v2"
)

type EmailInput struct {
	F  string //From
	T  string // To
	S  string //Subject
	B  string //Body
	BC string //BCC
	C  string //CC
}

func SendMail(pEmailDetails EmailInput, pClientId string) {

	// Create a new mail message
	email := mail.NewMessage()

	// Set the "From" field
	email.SetHeader("From", email.FormatAddress("profilematcher2024@gmail.com", "ProfileMatcher"))

	// Set the "To" field
	email.SetHeader("To", pEmailDetails.T)

	// Set the "Cc" field, if provided
	if pEmailDetails.C != "" {
		email.SetHeader("Cc", pEmailDetails.C)
	}
	if pEmailDetails.BC != "" {
		email.SetHeader("Bcc", pEmailDetails.BC)
	}

	// Set the "Subject" and "Body"
	email.SetHeader("Subject", pEmailDetails.S)
	email.SetBody("text/html", pEmailDetails.B)

	// Set up SMTP server configuration (Gmail in this case)
	d := mail.NewDialer("smtp.gmail.com", 587, "profilematcher2024@gmail.com", "jmhn bsvq sqbz uvms")

	// For security, use app-specific passwords instead of your Gmail password. Follow Gmail's instructions for creating app passwords.

	// Skip TLS verification (optional, for testing purposes)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Log the details for debugging
	log.Println("Sending email to:", pEmailDetails.T)
	if pEmailDetails.C != "" {
		log.Println("CC:", pEmailDetails.C)
	}
	log.Println("Subject:", pEmailDetails.S)
	log.Println("email", email)
	// Send the email and handle errors
	if err := d.DialAndSend(email); err != nil {
		log.Println("Error sending email:", err)
	}
}
