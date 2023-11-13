package messaging

import (
	"github.com/streadway/amqp"
	"go_task/protofile"
	"log"
)

func (p *Publisher) PublishAccountBalance(accountBalance *protofile.CreateAccountBalance) error {
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
			ContentType: "application/json",
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

func (p *Publisher) PublishOrder(orderBody []byte) error {
	err := p.channel.Publish(
		"your_exchange_name",
		"order_routing_key",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/octet-stream", // Установка типа контента на application/octet-stream
			Body:        orderBody,
		},
	)

	if err != nil {
		log.Println("Failed to publish order message:", err)
		return err
	}

	log.Println("Order message published successfully")
	return nil
}
