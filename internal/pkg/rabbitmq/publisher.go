package rabbitmq

import "github.com/streadway/amqp"

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
