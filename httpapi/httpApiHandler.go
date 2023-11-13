package httpapi

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go_task/messaging"
	"net/http"
)

type HTTPHandler struct {
	publisher *messaging.Publisher
	// Другие зависимости, если необходимо
}

func NewHTTPHandler(publisher *messaging.Publisher) *HTTPHandler {
	return &HTTPHandler{
		publisher: publisher,
	}
}

func (h *HTTPHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/createAccountBalance", h.CreateAccountBalance)
}

func (h *HTTPHandler) CreateAccountBalance(c *gin.Context) {
	var request structure.CreateAccountBalance

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Преобразование структуры в бинарный формат протобуфа
	protoData, err := json.Marshal(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal account balance message"})
		return
	}

	// Отправка сообщения в RabbitMQ
	err = h.publisher.PublishAccountBalance(protoData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish account balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account balance created successfully"})
}
