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
	publishingChannel chan amqp.Publishing
	listeningChannel  chan amqp.Delivery
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
		publishingChannel: make(chan amqp.Publishing),
		listeningChannel:  make(chan amqp.Delivery),
	}

	return &MessageBroker{
		conn:    conn,
		channel: channel,
		handler: handler,
	}, nil
}

func (m *MessageBroker) GetPublishingChannel() chan<- amqp.Publishing {
	return m.handler.publishingChannel
}

func (m *MessageBroker) GetListeningChannel() <-chan amqp.Delivery {
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
