package rabbitmq

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

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
