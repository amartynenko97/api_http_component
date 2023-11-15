package messaging

import (
	"github.com/streadway/amqp"
	"log"
)

type RabbitMQConfig struct {
	URL string
}

type MessageBroker struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	handler *Handler
}

type Handler struct {
	publishingChannel PublishingChannel
	listeningChannel  ListeningChannel
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

	handler := &Handler{
		publishingChannel: NewPublisher(channel),
		listeningChannel:  NewListener(channel),
	}

	return &MessageBroker{
		conn:    conn,
		channel: channel,
		handler: handler,
	}, nil
}

func (m *MessageBroker) GetPublishingChannel() PublishingChannel {
	return m.handler.publishingChannel
}

func (m *MessageBroker) GetListeningChannel() ListeningChannel {
	return m.handler.listeningChannel
}

func (m *MessageBroker) Close() {
	if m.channel != nil {
		m.channel.Close()
	}

	if m.conn != nil {
		m.conn.Close()
	}
}
