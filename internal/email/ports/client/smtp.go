package client

type SMTP interface {
	SendFilePDF(to []string, subject, body string, data []byte) error
}
