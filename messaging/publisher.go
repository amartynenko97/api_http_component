package messaging

import (
	"github.com/streadway/amqp"
	"log"
)

type Publisher struct {
	channel *amqp.Channel
	// Other fields, if any
}

// NewPublisher creates a new instance of Publisher
func NewPublisher(channel *amqp.Channel) *Publisher {
	return &Publisher{
		channel: channel,
	}
}

// PublishAccountBalance publishes an account balance message to RabbitMQ
func (p *Publisher) PublishAccountBalance(protoData []byte) error {
	err := p.channel.Publish(
		"your_exchange_name",
		"account_balance_routing_key",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        protoData, // Update to use protoData directly
		},
	)

	if err != nil {
		log.Println("Failed to publish account balance message:", err)
		return err
	}

	log.Println("Account balance message published successfully")
	return nil
}
