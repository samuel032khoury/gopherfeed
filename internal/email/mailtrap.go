package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"time"

	"github.com/samuel032khoury/gopherfeed/internal/logger"
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
	logger    logger.Logger
}

func NewMailtrap(fromEmail, host, username, password string, port int, log logger.Logger) (*MailtrapClient, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("Mailtrap credentials are not set")
	}

	if log == nil {
		log = logger.NewNoopLogger()
	}

	dialer := gomail.NewDialer(host, port, username, password)

	return &MailtrapClient{
		fromEmail: fromEmail,
		dialer:    dialer,
		logger:    log,
	}, nil
}

func (mt *MailtrapClient) Send(to string, templatePath string, data any) error {
	// Create a new message
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", mt.fromEmail)
	message.SetHeader("To", to)

	// Template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+templatePath)
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
			mt.logger.Warnw("Failed to send email",
				"to", to,
				"attempt", i+1,
				"max_retries", maxRetries,
				"error", err,
			)
			// Exponential backoff
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		return nil
	}

	return fmt.Errorf("failed to send email to %v after %d attempts", to, maxRetries)
}
