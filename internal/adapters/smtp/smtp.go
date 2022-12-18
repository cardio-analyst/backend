package smtp

import (
	"crypto/tls"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	mail "github.com/xhit/go-simple-mail/v2"

	"github.com/cardio-analyst/backend/internal/config"
	"github.com/cardio-analyst/backend/internal/ports/smtp"
)

var _ smtp.Client = (*Client)(nil)

type Client struct {
	smtpClient *mail.SMTPClient
	username   string
}

func NewClient(cfg config.SMTPConfig) (*Client, error) {
	smtpServer := mail.NewSMTPClient()

	smtpServer.Host = cfg.Host
	smtpServer.Port = cfg.Port
	smtpServer.Username = cfg.Username
	smtpServer.Password = cfg.Password
	smtpServer.Encryption = mail.EncryptionSSLTLS
	smtpServer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	smtpClient, err := smtpServer.Connect()
	if err != nil {
		return nil, err
	}

	return &Client{
		smtpClient: smtpClient,
		username:   cfg.Username,
	}, nil
}

func (c *Client) SendFile(to []string, subject, body, filePath string) error {
	emailMsg := mail.NewMSG()

	emailMsg.SetFrom(c.username)
	emailMsg.AddTo(to...)

	emailMsg.SetSubject(subject)
	emailMsg.SetBody(mail.TextHTML, body)

	emailMsg.Attach(&mail.File{
		FilePath: filePath,
		Name:     fmt.Sprintf("%v.pdf", time.Now().Format("2006_01_02_15_04_05")),
		Inline:   true,
	})

	if emailMsg.Error != nil {
		log.Warn("failed to send email through smtp: %v", emailMsg.Error)
		return emailMsg.Error
	}

	return emailMsg.Send(c.smtpClient)
}

func (c *Client) Close() error {
	return c.smtpClient.Close()
}
