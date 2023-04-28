package rabbitmq

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Client struct {
	uri string

	exchange   string
	routingKey string
	queue      string
	handler    func(data []byte) error

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

	MessagesHandler func(data []byte) error
}

func NewClient(opts ClientOptions) *Client {
	return &Client{
		uri:        fmt.Sprintf("amqp://%v:%v@%v:%v/", opts.User, opts.Password, opts.Host, opts.Port),
		exchange:   opts.ExchangeName,
		routingKey: opts.RoutingKey,
		queue:      opts.QueueName,
		handler:    opts.MessagesHandler,
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

func (c *Client) createExchange(name string) error {
	if !c.connected {
		if err := c.connect(); err != nil {
			return err
		}
	}

	return c.ch.ExchangeDeclare(
		name,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
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

func (c *Client) bindQueue(name, routingKey, exchange string) error {
	if !c.connected {
		if err := c.connect(); err != nil {
			return err
		}
	}

	return c.ch.QueueBind(
		name,
		routingKey,
		exchange,
		false,
		nil,
	)
}

func (c *Client) Connect() error {
	if err := c.createExchange(c.exchange); err != nil {
		log.Errorf("RMQ: creating exchange %q: %v", c.exchange, err)
		return err
	}

	if err := c.createQueue(c.queue); err != nil {
		log.Errorf("RMQ: creating queue %q: %v", c.queue, err)
		return err
	}

	if err := c.bindQueue(c.queue, c.routingKey, c.exchange); err != nil {
		log.Errorf("RMQ: binding queue %q to exchange %q with key %q: %v", c.queue, c.exchange, c.routingKey, err)
		return err
	}

	return nil
}

func (c *Client) Consume() error {
	if c.queue == "" {
		return errors.New("no RMQ queue provided")
	}

	for {
		if !c.connected {
			if err := c.connect(); err != nil {
				log.Errorf("RMQ: %v", err)
				time.Sleep(time.Second * 5)
				continue
			}
		}

		var err error
		c.messages, err = c.ch.Consume(
			c.queue,
			"",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Errorf("RMQ: failed to consume queue: %v", err)
			time.Sleep(time.Second * 5)
			continue
		}

		log.Debug("RMQ: consuming messages")

		for msg := range c.messages {
			err = c.handler(msg.Body)
			if err != nil {
				log.Errorf("RMQ: %v", err)
				_ = msg.Nack(false, true)
				continue
			}
			err = msg.Ack(false)
			if err != nil {
				return err
			}
		}
	}
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
