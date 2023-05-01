package client

type EmailPublisher interface {
	Publish(msg []byte) error
}

type FeedbackPublisher interface {
	Publish(msg []byte) error
}
