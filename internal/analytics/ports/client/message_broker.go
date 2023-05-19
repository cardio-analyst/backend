package client

type Consumer interface {
	Consume(handler func(data []byte) error) error
}
