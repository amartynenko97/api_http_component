package httpapi

import (
	"github.com/gin-gonic/gin"
	"go_task/messaging"
	"go_task/protofile"
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
	router.POST("/createOrder", h.CreateOrder)
	router.POST("/createAccountBalance", h.CreateAccountBalance)
}

func (h *HTTPHandler) CreateOrder(c *gin.Context) {
	var request protofile.CreateOrderRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Отправка сообщения в RabbitMQ
	err := h.publisher.PublishOrder(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order created successfully"})
}

func (h *HTTPHandler) CreateAccountBalance(c *gin.Context) {
	var request protofile.CreateAccountBalance

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Отправка сообщения в RabbitMQ
	err := h.publisher.PublishAccountBalance(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish account balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account balance created successfully"})
}
