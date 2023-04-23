package smtp

import (
	"crypto/tls"
	"fmt"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"

	"github.com/cardio-analyst/backend/internal/email/ports/client"
)

var _ client.SMTP = (*Client)(nil)

type Client struct {
	smtpServer *mail.SMTPServer
	smtpClient *mail.SMTPClient
	username   string
}

func NewClient(host string, port int, username, password string) *Client {
	smtpServer := mail.NewSMTPClient()

	smtpServer.Host = host
	smtpServer.Port = port
	smtpServer.Username = username
	smtpServer.Password = password
	smtpServer.Encryption = mail.EncryptionSSLTLS
	smtpServer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &Client{
		smtpServer: smtpServer,
		username:   username,
	}
}

func (c *Client) Connect() error {
	smtpClient, err := c.smtpServer.Connect()
	if err != nil {
		return err
	}

	c.smtpClient = smtpClient
	return nil
}

func (c *Client) SendFilePDF(to []string, subject, body string, fileData []byte) error {
	emailMsg := mail.NewMSG()

	emailMsg.SetFrom(c.username)
	emailMsg.AddTo(to...)

	emailMsg.SetSubject(subject)
	emailMsg.SetBody(mail.TextHTML, body)

	emailMsg.Attach(&mail.File{
		Name:     fmt.Sprintf("%v.pdf", time.Now().Format("2006_01_02_15_04_05")),
		Data:     fileData,
		MimeType: "application/pdf",
		Inline:   true,
	})

	if emailMsg.Error != nil {
		return fmt.Errorf("sending email through SMTP: %v", emailMsg.Error)
	}

	if err := emailMsg.Send(c.smtpClient); err != nil {
		return fmt.Errorf("sending email through SMTP: %v", err)
	}

	return nil
}

func (c *Client) Close() error {
	return c.smtpClient.Close()
}
