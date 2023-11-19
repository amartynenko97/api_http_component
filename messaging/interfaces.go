package messaging

import (
	"github.com/streadway/amqp"
)

type PublishingChannel interface {
	PublishCreateAccountBalances(protoData []byte) error
	//PublishGetAccountBalances(protoData []byte) error
}

type ListeningChannel interface {
	ConsumeCreateAccountBalances() <-chan amqp.Delivery
	//ConsumeGetAccountBalances() <-chan amqp.Delivery
}
