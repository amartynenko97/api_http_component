package httpapi

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	//orderpb "github.com/yourgithubusername/yourrepository/order" // Импорт вашего прото-пакета
	"go_task/messaging"
)

type API struct {
	// Здесь могут быть другие зависимости
}

func NewAPI() *API {
	return &API{}
}

func (api *API) CreateOrderHandler(c *gin.Context) {
	var orderRequest orderpb.CreateOrderRequest // Использование модели из протофайла

	// Получаем JSON из запроса
	requestData, err := c.GetRawData()
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to parse request data"})
		return
	}

	// Парсим JSON в протоструктуру
	if err := json.Unmarshal(requestData, &orderRequest); err != nil {
		c.JSON(400, gin.H{"error": "Failed to parse JSON"})
		return
	}

	// Отправляем запрос в RabbitMQ
	err = messaging.PublishOrderRequest(orderRequest)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to process order"})
		return
	}

	c.JSON(200, gin.H{"message": "Order sent successfully"})
}
