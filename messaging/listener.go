package messaging

import (
	"github.com/streadway/amqp"
	"log"
)

type Listener struct {
	channel           *amqp.Channel
	accountBalances   <-chan amqp.Delivery
	createOrderEvents <-chan amqp.Delivery
}

func NewListener(channel *amqp.Channel) *Listener {
	return &Listener{
		channel: channel,
	}
}

func (l *Listener) StartListening() {
	// Запустите цикл для прослушивания обеих очередей
	go func() {
		for {
			select {
			case accountBalance := <-l.accountBalances:
				// Обработка сообщения для балансов
				// accountBalance.Body содержит тело сообщения
				// ...

			case createOrder := <-l.createOrderEvents:
				// Обработка сообщения для создания заказов
				// createOrder.Body содержит тело сообщения
				// ...
			}
		}
	}()
}

func (l *Listener) ConsumeAccountBalances() <-chan amqp.Delivery {
	return l.consume("account_balance_queue_name", "account_balance_consumer")
}

func (l *Listener) ConsumeCreateOrders() <-chan amqp.Delivery {
	return l.consume("create_order_queue_name", "create_order_consumer")
}

func (l *Listener) consume(queueName, consumerName string) <-chan amqp.Delivery {
	messages, err := l.channel.Consume(
		queueName,
		consumerName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal("Failed to register a consumer:", err)
	}

	return messages
}
