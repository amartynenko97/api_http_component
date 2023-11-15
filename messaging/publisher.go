package messaging

import (
	"github.com/streadway/amqp"
	"log"
)

type Publisher struct {
	channel *amqp.Channel
}

func NewPublisher(channel *amqp.Channel) *Publisher {
	return &Publisher{
		channel: channel,
	}
}

func (p *Publisher) PublishAccountBalance(protoData []byte) error {
	return p.publish(protoData, "account_balance_routing_key")
}

func (p *Publisher) PublishCreateOrder(protoData []byte) error {
	return p.publish(protoData, "create_order_routing_key")
}

func (p *Publisher) publish(protoData []byte, routingKey string) error {
	err := p.channel.Publish(
		"your_exchange_name",
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        protoData,
		},
	)

	if err != nil {
		log.Println("Failed to publish message:", err)
		return err
	}

	log.Printf("Message published successfully with routing key: %s\n", routingKey)
	return nil
}
