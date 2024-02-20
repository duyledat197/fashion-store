package email

import "context"

type EmailData struct {
	From    string
	To      string
	Content string
}

type Provider interface {
	SendMail(ctx context.Context, data *EmailData) error
}
