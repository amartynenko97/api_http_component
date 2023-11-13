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
	channel chan amqp.Publishing
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
		channel: make(chan amqp.Publishing),
	}

	return &MessageBroker{
		conn:    conn,
		channel: channel,
		handler: handler,
	}, nil
}

func (m *MessageBroker) GetChannel() chan<- amqp.Publishing {
	return m.handler.channel
}

func (m *MessageBroker) Close() {
	if m.channel != nil {
		m.channel.Close()
	}

	if m.conn != nil {
		m.conn.Close()
	}
}
