package tasks

import (
	"bytes"
	"errors"
	"log"
	"net/smtp"
	"text/template"
)

func SendVerificationEmail(configEmail string, configPassword string) func(email string, code string) error {
	return func(email string, code string) error {
		log.Printf("Sending email to %s with code: %v", email, code)

		smtpHost := "smtp.gmail.com"
		smtpPort := "587"
		smtpUsername := configEmail
		smtpPassword := configPassword

		auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

		emailTemplate := `From: {{.From}}
		To: {{.To}}
		Subject: Shortener Email Verification

		Hello,

		Please use the following code to verify your email address: {{.Code}}

		If you didn't request this, please ignore this email.

		Best regards,
		Shortener Team
	`

		data := struct {
			From string
			To   string
			Code string
		}{
			From: smtpUsername,
			To:   email,
			Code: code,
		}

		var body bytes.Buffer
		t := template.Must(template.New("email").Parse(emailTemplate))
		if err := t.Execute(&body, data); err != nil {
			log.Printf("error generating email: %v", err)
			return errors.New("error generating email")
		}

		err := smtp.SendMail(
			smtpHost+":"+smtpPort,
			auth,
			smtpUsername,
			[]string{email},
			body.Bytes(),
		)
		if err != nil {
			log.Printf("error sending email: %v", err)
			return errors.New("error sending email")
		}

		return nil
	}
}
