package client

type SMTP interface {
	SendFile(to []string, subject, body, filePath string) error
}
