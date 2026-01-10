package mailer

import "embed"

const (
	FromName           = "GopherFeed"
	maxRetries         = 3
	UserInviteTemplate = "user_invitation.gtpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any) error
}
