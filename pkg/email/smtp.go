package email

import (
	"context"
	"net/smtp"
)

type SmtpAccount struct {
	Host     string
	Username string
	Password string
}

type Smtp struct {
	account *SmtpAccount
}

func NewSmtp(account *SmtpAccount) Provider {
	return &Smtp{account}
}

func (s *Smtp) SendMail(ctx context.Context, data *EmailData) error {
	auth := smtp.PlainAuth("", s.account.Username, s.account.Password, s.account.Host)
	to := []string{data.To}

	if err := smtp.SendMail("smtp.gmail.com:587", auth, data.From, to, []byte(data.Content)); err != nil {
		return err
	}

	return nil
}
