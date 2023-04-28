package service

type FeedbackService interface {
	MessagesHandler() func(data []byte) error
}
