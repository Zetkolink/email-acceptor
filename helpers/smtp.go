package helpers

import (
	"email-sender/models"
	"net/smtp"
	"strings"
)

type SmtpHandler struct {
	host string
	auth smtp.Auth
}

func InitSmtp(cfg Config) *SmtpHandler {
	auth := smtp.PlainAuth(
		cfg.Smtp.Identity,
		cfg.Smtp.Username,
		cfg.Smtp.Password,
		cfg.Smtp.Hostname,
	)

	return &SmtpHandler{
		auth: auth,
		host: cfg.Smtp.Hostname + ":" + cfg.Smtp.Port,
	}
}

func (s SmtpHandler) buildBody(msg *models.MessageRequest) []byte {
	body := "From:" + msg.Sender + "\n" + "To:" + strings.Join(msg.To, ",") + "\n"
	if msg.Subject != nil {
		body += "Subject:" + *msg.Subject + "\n"
	}
	body += msg.Message

	return []byte(body)
}

func (s SmtpHandler) SendMail(msg *models.MessageRequest) error {
	err := smtp.SendMail(
		s.host,
		s.auth,
		msg.Sender,
		msg.To,
		s.buildBody(msg),
	)
	if err != nil {
		return err
	}

	return nil
}
