package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"
)

type MailInterface interface {
	SendMail(to string, from string, subject string, buffer *bytes.Buffer) error
}

type MailService struct {
	mailClient *smtp.Client
}

func NewMailService() *MailService {

	return nil

	smtpAddr := os.Getenv("SMTP_ADDR")
	smtpServerName := os.Getenv("SMTP_SERVER_NAME")
	c, err := smtp.Dial(smtpAddr)
	if err != nil {
		log.Fatal(err)
	}
	c.StartTLS(&tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServerName,
	})

	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	err = c.Auth(auth)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	m := new(MailService)
	m.mailClient = c
	return m
}

func (self *MailService) SendMail(to string, from string, subject string, buffer *bytes.Buffer) error {
	_ = self.mailClient.Mail(from)
	_ = self.mailClient.Rcpt(to)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	s := fmt.Sprintf("Subject: %s\n", subject)
	wc, err := self.mailClient.Data()
	if err != nil {
		return err
	}
	defer wc.Close()
	if _, err = wc.Write([]byte(s + mime)); err != nil {
		return err
	}
	if _, err = buffer.WriteTo(wc); err != nil {
		return err
	}

	return nil
}
