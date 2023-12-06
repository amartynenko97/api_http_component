package messaging

import (
	"github.com/streadway/amqp"
)

type PublishingChannel interface {
	PublishCreateAccountToBalances(protoData []byte) error
}

type ListeningChannel interface {
	ConsumeCreateAccountFromBalances() <-chan amqp.Delivery
}
