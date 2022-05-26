package services

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/sethetter/go-web-starter/pkg/config"
)

type IEmailService interface {
	SendEmail(string, string, string, string) error
}

type EmailService struct {
	Conf *config.EmailConfig
}

func (svc *EmailService) SendEmail(from, to, subject, body string) error {
	msg := fmt.Sprintf(
		"From: %s <%s>\nTo: %s\nSubject: %s\nContent-Type: text/html; charset=UTF-8\n\n%s",
		from,
		svc.Conf.FromEmail,
		to,
		subject,
		body,
	)

	host := strings.Split(svc.Conf.SMTPHost, ":")[0]
	auth := smtp.PlainAuth("", svc.Conf.SMTPUsername, svc.Conf.SMTPPassword, host)
	return smtp.SendMail(svc.Conf.SMTPHost, auth, svc.Conf.FromEmail, []string{to}, []byte(msg))
}
