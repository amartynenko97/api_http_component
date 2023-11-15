package messaging

import (
	"github.com/streadway/amqp"
)

type PublishingChannel interface {
	PublishAccountBalance(protoData []byte) error
	PublishCreateOrder(protoData []byte) error
}

type ListeningChannel interface {
	ConsumeAccountBalances() <-chan amqp.Delivery
	ConsumeCreateOrders() <-chan amqp.Delivery
}
