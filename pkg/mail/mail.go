package mail

import (
	"context"
	"gopkg.in/gomail.v2"
	"sync"
	"time"
)

var (
	globalSender *SmtpSender
	once         sync.Once
)

func SetSender(sender *SmtpSender) {
	once.Do(func() {
		globalSender = sender
	})
}

func Send(ctx context.Context, to []string, cc []string, bcc []string, subject string, body string, file ...string) error {
	return globalSender.Send(ctx, to, cc, bcc, subject, body, file...)
}

func SendTo(ctx context.Context, to string, subject string, body string, file ...string) error {
	return globalSender.SendTo(ctx, to, subject, body, file...)
}

type SmtpSender struct {
	SmtpHost string
	Port     int
	FromName string
	FromMail string
	UserName string
	AuthCode string
}

func (s *SmtpSender) Send(ctx context.Context, to []string, cc []string, bcc []string, subject string, body string, file ...string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(s.FromMail, s.FromName))
	m.SetHeader("To", to...)
	m.SetHeader("Cc", cc...)
	m.SetHeader("Bcc", bcc...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	for _, f := range file {
		m.Attach(f)
	}

	d := gomail.NewDialer(s.SmtpHost, s.Port, s.UserName, s.AuthCode)
	return d.DialAndSend(m)
}

func (s *SmtpSender) SendTo(ctx context.Context, to string, subject string, body string, file ...string) error {
	var err error
	for i := 0; i < 3; i++ {
		err = s.Send(ctx, []string{to}, nil, nil, subject, body, file...)
		if err == nil {
			time.Sleep(time.Millisecond * 500)
			continue
		}
		err = nil
		break
	}
	return err
}
