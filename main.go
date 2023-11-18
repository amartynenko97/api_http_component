package go_task

import (
	"github.com/gin-gonic/gin"
	"go_task/balances"
	"go_task/httpapi"
	"go_task/messaging"
	"golang.org/x/net/context"
	"log"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

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

	httpHandler := httpapi.NewHTTPHandler(messageBroker.GetPublishingChannel())

	balancesHandler := balances.NewBalancesHandler(messageBroker.GetListeningChannel())

	httpHandler.RegisterRoutes(router)

	wg.Add(1)
	go func() {
		defer wg.Done()
		balancesHandler.StartListener(ctx)
	}()

	go func() {
		err := router.Run(":8080")
		if err != nil {
			log.Fatal("Failed to start the server:", err)
			cancel()
		}
	}()
	wg.Wait()
}
