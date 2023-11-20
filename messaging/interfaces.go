package messaging

import (
	"github.com/streadway/amqp"
)

type PublishingChannel interface {
	PublishCreateAccountBalancesToBalances(protoData []byte) error
	PublishCreateAccountBalancesToHttpApi(protoData []byte) error
}

type ListeningChannel interface {
	ConsumeCreateAccountBalancesFromBalances() <-chan amqp.Delivery
	ConsumeCreateAccountBalancesFromHttpApi() <-chan amqp.Delivery
}
