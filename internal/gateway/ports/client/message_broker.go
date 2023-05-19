package client

type Publisher interface {
	Publish(msg []byte) error
}
