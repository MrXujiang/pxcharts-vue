package mail

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"mvtable/internal/pkg/config"
	"sync"
)

type SmtpSender struct {
	Host     string
	Port     int
	FromName string
	FromMail string
	UserName string
	AuthCode string
}

var (
	globalSender *SmtpSender
	once         sync.Once
)

// Init 设置全局SMTP发送器
func Init(sender *config.MailConfig) {
	once.Do(func() {
		globalSender = &SmtpSender{
			Host:     sender.Host,
			Port:     sender.Port,
			FromName: sender.FromName,
			FromMail: sender.FromMail,
			UserName: sender.UserName,
			AuthCode: sender.AuthCode,
		}
	})
}

func SendPlanText(to, subject, body string) error {
	return globalSender.SendPlainText(to, subject, body)
}

func (a *SmtpSender) SendPlainText(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", a.FromMail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(a.Host, a.Port, a.FromMail, a.AuthCode)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %s", err.Error())
	}

	return nil
}
