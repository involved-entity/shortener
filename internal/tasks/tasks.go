package tasks

import (
	"log"
)

type EmailPayload struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func SendVerificationEmail(payload EmailPayload) error {
	log.Printf("Sending email to %s with code: %s", payload.Email, payload.Code)

	// err := smtp.SendMail(...)
	// if err != nil { return err }

	return nil
}
