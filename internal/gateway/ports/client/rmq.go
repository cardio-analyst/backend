package client

type RabbitMQPublisher interface {
	Publish(msg []byte) error
}
