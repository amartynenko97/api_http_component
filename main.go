package go_task

import (
	"go_task/httpapi"
	"go_task/messaging"
)

func main() {
	// Start RabbitMQ listener to receive messages
	go messaging.StartListeningForOrders()

	// Start HTTP API server
	httpapi.StartServer()
}
