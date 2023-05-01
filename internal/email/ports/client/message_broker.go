package client

type EmailConsumer interface {
	Consume(handler func(data []byte) error) error
}
