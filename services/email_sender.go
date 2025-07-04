package services

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

type EmailSender struct {
}

func (emailSender *EmailSender) SendEmail(from string, to []string, message []byte) error {
	// Load the environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// Load configuration from environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	email := os.Getenv("SENDER_EMAIL")
	password := os.Getenv("PASSWORD")

	// Authentication
	auth := smtp.PlainAuth("", email, password, smtpHost)

	// Send email
	return smtp.SendMail(fmt.Sprintf("%s:%s", smtpHost, smtpPort), auth, email, to, message)
}
