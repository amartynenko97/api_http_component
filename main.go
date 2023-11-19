package go_task

import (
	"github.com/gin-gonic/gin"
	"go_task/balances"
	"go_task/constants"
	"go_task/httpapi"
	"go_task/messaging"
	"golang.org/x/net/context"
	"log"
	"net/http"
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

	errorChannel := make(chan error, 1)

	go func() {
		defer wg.Done()
		if err := balancesHandler.StartListener(ctx); err != nil {
			errorChannel <- err
			cancel() // Отмена контекста при ошибке
		}
	}()

	select {
	case err := <-errorChannel:
		switch err := err.(type) {
		case *constants.CustomError:
			log.Printf("Custom error")
			c.JSON(http.StatusBadRequest, gin.H{"error": string(err.Type)})
		default:
			log.Printf("Error in one of the components")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
	case <-ctx.Done():
	}

	go func() {
		err := router.Run(":8080")
		if err != nil {
			log.Fatal("Failed to start the server:", err)
			cancel()
		}
	}()
	wg.Wait()
}
