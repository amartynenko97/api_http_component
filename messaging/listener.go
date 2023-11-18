package messaging

import (
	"github.com/streadway/amqp"
	"log"
)

type Listener struct {
	channel *amqp.Channel
}

func NewListener(channel *amqp.Channel) *Listener {
	return &Listener{
		channel: channel,
	}
}

func (l *Listener) ConsumeAccountBalances() <-chan amqp.Delivery {
	return l.consume("account_balance_queue_name", "account_balance_consumer")
}

func (l *Listener) ConsumeCreateOrders() <-chan amqp.Delivery {
	return l.consume("create_order_queue_name", "create_order_consumer")
}

func (l *Listener) consume(queueName, consumerName string) <-chan amqp.Delivery {
	messages, err := l.channel.Consume(
		queueName,
		consumerName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal("Failed to register a consumer:", err)
	}

	return messages
}
