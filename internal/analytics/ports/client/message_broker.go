package client

type FeedbackConsumer interface {
	Consume(handler func(data []byte) error) error
}
