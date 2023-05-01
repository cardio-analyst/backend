package service

type Email interface {
	ListenToEmailMessages() error
}
