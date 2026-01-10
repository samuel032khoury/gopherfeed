package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"time"

	gomail "gopkg.in/mail.v2"
)

const (
	FromName           = "GopherFeed"
	maxRetries         = 3
	UserInviteTemplate = "user_invitation.gtpl"
)

//go:embed "templates"
var FS embed.FS

type MailtrapClient struct {
	fromEmail string
	dialer    *gomail.Dialer
}

func NewMailtrap(fromEmail, host, username, password string, port int) (*MailtrapClient, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("Mailtrap credentials are not set")
	}

	dialer := gomail.NewDialer(host, port, username, password)

	return &MailtrapClient{
		fromEmail: fromEmail,
		dialer:    dialer,
	}, nil
}

func (mt *MailtrapClient) Send(templateFile, username, email string, data any) error {
	// Create a new message
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", mt.fromEmail)
	message.SetHeader("To", email)

	// Template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(subject, "subject", data); err != nil {
		return err
	}

	body := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(body, "body", data); err != nil {
		return err
	}

	// Set subject and body
	message.SetHeader("Subject", subject.String())
	message.SetBody("text/html", body.String())

	// Send with retry logic
	for i := range maxRetries {
		err := mt.dialer.DialAndSend(message)
		if err != nil {
			log.Printf("Failed to send email to %v, attempt %d of %d: %v", email, i+1, maxRetries, err.Error())
			log.Println(err)
			// Exponential backoff
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		log.Printf("Email sent successfully to %v", email)
		return nil
	}

	return fmt.Errorf("failed to send email to %v after %d attempts", email, maxRetries)
}
