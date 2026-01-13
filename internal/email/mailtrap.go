package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"time"

	"go.uber.org/zap"
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
	logger    *zap.SugaredLogger
}

func NewMailtrap(fromEmail, host, username, password string, port int, log *zap.SugaredLogger) (*MailtrapClient, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("mailtrap credentials are not set")
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
