package rabbitmq

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

func (c *Client) Consume(handler func(data []byte) error) error {
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
			err = handler(msg.Body)
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
