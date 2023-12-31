package api_http_component

import (
	"api_http_component/config"
	"api_http_component/httpapi"
	"api_http_component/logger"
	"api_http_component/messaging"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

	logger := logger.SetupLogger()

	rabbitMQConfig := messaging.RabbitMQConfig{
		URL: config.RabbitConfig,
	}

	messageBroker, err := messaging.NewMessageBroker(rabbitMQConfig)
	if err != nil {
		logger.Info("Failed to initialize MessageBroker")
		return
	}
	defer messageBroker.Close()

	router := gin.Default()
	httpHandler := httpapi.NewHTTPHandler(logger, messageBroker.GetPublishingChannel(), messageBroker.GetListeningChannel())
	httpHandler.RegisterRoutes(router)
	listenerReady := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		httpHandler.StartListener(ctx, listenerReady)
	}()

	<-listenerReady

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := router.Run(":8080")
		if err != nil {
			cancel()
		}
	}()
	wg.Wait()
}
