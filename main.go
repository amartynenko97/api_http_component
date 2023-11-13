package go_task

import (
	"github.com/gin-gonic/gin"
	"go_task/httpapi"
	"go_task/messaging"
	"log"
)

func main() {
	// Инициализация маршрутизатора Gin
	router := gin.Default()

	// Инициализация и настройка обработчика HTTP API
	httpHandler := httpapi.NewHTTPHandler()

	// Инициализация и настройка мессенджера (publisher и listener)
	messageBroker := messaging.NewMessageBroker()
	publisher := messaging.NewPublisher(messageBroker.GetChannel())
	listener := messaging.NewListener(messageBroker.GetChannel())

	// Передача экземпляра publisher в httpHandler
	httpHandler.SetPublisher(publisher)

	// Регистрация маршрутов
	httpHandler.RegisterRoutes(router)

	// Запуск HTTP-сервера
	err := router.Run(":8080")
	if err != nil {
		log.Fatal("Failed to start the server: ", err)
	}
}
