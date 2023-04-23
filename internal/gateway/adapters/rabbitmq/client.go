package rabbitmq

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
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

func (c *Client) connect() (err error) {
	if c.uri == "" {
		return errors.New("no RMQ hostname provided")
	}

	log.Debugf("RMQ: connecting to host %q", c.uri)

	c.conn, err = amqp.Dial(c.uri)
	if err != nil {
		return fmt.Errorf("failed to connect to RMQ: %w", err)
	}
	c.ch, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}

	go func() {
		e := <-c.conn.NotifyClose(make(chan *amqp.Error))
		if e != nil {
			log.Errorf("RMQ: disconnected with error: %v", e.Reason)
			c.connected = false
		}
	}()
	c.connected = true

	log.Debug("RMQ: successfully connected to host")

	return nil
}

func (c *Client) createQueue(name string) error {
	if !c.connected {
		if err := c.connect(); err != nil {
			return err
		}
	}

	_, err := c.ch.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}

func (c *Client) Connect() error {
	if err := c.createQueue(c.queue); err != nil {
		log.Errorf("RMQ: creating queue %q: %v", c.queue, err)
		return err
	}

	return nil
}

func (c *Client) Publish(msg []byte) error {
	if !c.connected {
		if err := c.connect(); err != nil {
			return err
		}
	}

	err := c.ch.Publish(
		c.exchange,
		c.routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
	return err
}

func (c *Client) Close() (err error) {
	if c.ch != nil {
		if err = c.ch.Close(); err != nil {
			log.Warn(err)
		}
	}

	if c.conn != nil {
		if err = c.conn.Close(); err != nil {
			log.Warn(err)
		}
	}

	c.connected = false
	return
}
