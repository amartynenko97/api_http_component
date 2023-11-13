package messaging

import (
	"github.com/streadway/amqp"
	"go_task/protofile"
	"google.golang.org/protobuf/proto"
	"log"
)

type Listener struct {
	channel *amqp.Channel
	// Другие поля, если нужно
}

func NewListener(channel *amqp.Channel) *Listener {
	return &Listener{
		channel: channel,
	}
}

func (l *Listener) ListenForAccountBalances() {
	messages, err := l.channel.Consume(
		"your_queue_name_account_balance",
		"account_balance_consumer",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal("Failed to register a consumer for account balances:", err)
	}

	for msg := range messages {
		// Преобразование бинарных данных из RabbitMQ в protobuf
		var accountBalance protofile.CreateAccountBalance
		err := proto.Unmarshal(msg.Body, &accountBalance)
		if err != nil {
			log.Println("Failed to unmarshal account balance message:", err)
			// Ваш код обработки ошибки, если это необходимо
			continue
		}

		// Логика обработки сообщения о создании баланса аккаунта
		// Ваш код здесь
		log.Printf("Received account balance message: %+v\n", accountBalance)
	}
}
