package messaging

import (
	"github.com/streadway/amqp"
)

type PublishingChannel interface {
	PublishCreateAccountToBalances(protoData []byte) error
}

type ListeningChannel interface {
	ConsumeCreateAccountFromBalances(stopCh <-chan struct{}) <-chan amqp.Delivery
}
