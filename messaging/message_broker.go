package messaging

import (
	"github.com/streadway/amqp"
	"log"
)

type RabbitMQConfig struct {
	URL string
}

type MessageBroker struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	publisher PublishingChannel
	listener  ListeningChannel
}

func NewMessageBroker(config RabbitMQConfig) (*MessageBroker, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		log.Println("Failed to connect to RabbitMQ:", err)
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Println("Failed to open a channel:", err)
		return nil, err
	}

	return &MessageBroker{
		conn:      conn,
		channel:   channel,
		publisher: NewPublisher(channel),
		listener:  NewListener(channel),
	}, nil
}

func (m *MessageBroker) GetPublishingChannel() PublishingChannel {
	return m.publisher
}

func (m *MessageBroker) GetListeningChannel() ListeningChannel {
	return m.listener
}

func (m *MessageBroker) Close() {
	if m.channel != nil {
		err := m.channel.Close()
		if err != nil {
			return
		}
	}

	if m.conn != nil {
		err := m.conn.Close()
		if err != nil {
			return
		}
	}
}
