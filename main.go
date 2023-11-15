package go_task

import (
	"github.com/gin-gonic/gin"
	"go_task/httpapi"
	"go_task/messaging"
	"log"
)

func main() {
	rabbitMQConfig := messaging.RabbitMQConfig{
		URL: "amqp://guest:guest@localhost:5672/",
	}

	messageBroker, err := messaging.NewMessageBroker(rabbitMQConfig)
	if err != nil {
		log.Fatal("Failed to initialize MessageBroker:", err)
		return
	}
	defer messageBroker.Close()

	router := gin.Default()

	httpHandler := httpapi.NewHTTPHandler(messageBroker.GetPublishingChannel(), messageBroker.GetListeningChannel())

	httpHandler.StartListener()

	httpHandler.RegisterRoutes(router)

	err = router.Run(":8080")
	if err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
