package mail

import (
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/mail-service/cmd/api/config"
	"github.com/jordan-wright/email"
	"net/smtp"
)

// EmailMessenger is an interface for sending emails
type EmailMessenger interface {
	SendEmail(message EmailMessage) error
}

// EmailMessage represents an email message
type EmailMessage struct {
	Subject     string   // Subject of the email
	Content     string   // HTML content of the email
	To          []string // Primary recipients
	Cc          []string // Carbon copy recipients
	Bcc         []string // Blind carbon copy recipients
	AttachFiles []string // List of file paths to be attached to the email
}

// GmailSender is a struct implementing the EmailMessenger interface for sending emails using Gmail
type GmailSender struct {
	Name              string
	fromEmailAddress  string
	fromEmailPassword string
	smtpAuthAddress   string
	smtpServerAddress string
}

// NewGmailSender creates a new instance of GmailSender
func NewGmailSender(cfg *config.Config) GmailSender {
	return GmailSender{
		smtpAuthAddress:   cfg.SmtpAuthAddress,
		smtpServerAddress: cfg.SmtpServerAddress,
		Name:              cfg.EmailName,
		fromEmailAddress:  cfg.EmailAddress,
		fromEmailPassword: cfg.EmailPassword,
	}
}

// SendEmail sends an email using GmailSender
func (sender *GmailSender) SendEmail(message EmailMessage) error {
	// Create a new email instance
	e := email.NewEmail()

	// Set email properties
	e.From = fmt.Sprintf("%s <%s>", sender.Name, sender.fromEmailAddress)
	e.Subject = message.Subject
	e.HTML = []byte(message.Content)
	e.To = message.To
	e.Cc = message.Cc
	e.Bcc = message.Bcc

	// Attach files to the email
	for _, f := range message.AttachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}

	// Set up SMTP authentication
	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, sender.smtpAuthAddress)

	// Send the email using SMTP
	return e.Send(sender.smtpServerAddress, smtpAuth)
}
