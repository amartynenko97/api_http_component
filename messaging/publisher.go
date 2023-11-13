package messaging

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"go_task/protofile"
	"log"
)

type Publisher struct {
	channel *amqp.Channel
	// Другие поля, если нужно
}

func NewPublisher(channel *amqp.Channel) *Publisher {
	return &Publisher{
		channel: channel,
	}
}

func (p *Publisher) PublishAccountBalance(accountBalance *protofile.CreateAccountBalance) error {
	// Преобразование структуры в бинарный формат протобуфа
	messageBody, err := json.Marshal(accountBalance)
	if err != nil {
		log.Println("Failed to marshal account balance message:", err)
		return err
	}

	err = p.channel.Publish(
		"your_exchange_name",
		"account_balance_routing_key",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        messageBody,
		},
	)

	if err != nil {
		log.Println("Failed to publish account balance message:", err)
		return err
	}

	log.Println("Account balance message published successfully")
	return nil
}
