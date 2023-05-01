package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

type Client struct {
	uri string

	exchange   string
	routingKey string
	queue      string

	conn      *amqp.Connection
	ch        *amqp.Channel
	messages  <-chan amqp.Delivery
	connected bool
}

type ClientOptions struct {
	User     string
	Password string
	Host     string
	Port     int

	ExchangeName string
	RoutingKey   string
	QueueName    string
}

func NewClient(opts ClientOptions) *Client {
	return &Client{
		uri:        fmt.Sprintf("amqp://%v:%v@%v:%v/", opts.User, opts.Password, opts.Host, opts.Port),
		exchange:   opts.ExchangeName,
		routingKey: opts.RoutingKey,
		queue:      opts.QueueName,
	}
}
