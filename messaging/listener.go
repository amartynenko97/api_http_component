package messaging

import (
	"api_http_component/constants"
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

func (l *Listener) ConsumeCreateAccountFromBalances(stopCh <-chan struct{}) <-chan amqp.Delivery {
	return l.consume(constants.CreateAccountResponseQueue, stopCh)
}

func (l *Listener) consume(queueName string, stopCh <-chan struct{}) <-chan amqp.Delivery {
	messages := make(chan amqp.Delivery)

	go func() {
		defer close(messages)

		consumer, err := l.channel.Consume(
			queueName,
			"",
			true,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			log.Fatal("Failed to register a consumer:", err)
		}

		for {
			select {
			case <-stopCh:
				return
			case delivery, ok := <-consumer:
				if !ok {
					return
				}

				select {
				case messages <- delivery:
				case <-stopCh:
					return
				}
			}
		}
	}()

	return messages
}
