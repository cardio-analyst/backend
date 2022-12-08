package smtp

type Client interface {
	SendFile(to []string, subject, body, filePath string) error
}
