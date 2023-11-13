package httpapi

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"go_task/protofile"
	"net/http"
)

type HTTPHandler struct {
	publishingChannel chan<- amqp.Publishing
	listeningChannel  <-chan amqp.Delivery
}

func NewHTTPHandler(publishingChannel chan<- amqp.Publishing, listeningChannel <-chan amqp.Delivery) *HTTPHandler {
	return &HTTPHandler{
		publishingChannel: publishingChannel,
		listeningChannel:  listeningChannel,
	}
}

func (h *HTTPHandler) SetPublishingChannel(channel chan<- amqp.Publishing) {
	h.publishingChannel = channel
}

func (h *HTTPHandler) SetListeningChannel(channel <-chan amqp.Delivery) {
	h.listeningChannel = channel
}

func (h *HTTPHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/createAccountBalance", h.CreateAccountBalance)
}

func (h *HTTPHandler) CreateAccountBalance(c *gin.Context) {
	var request protofile.CreateOrderRequest

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
	err = h.publishingChannel.PublishAccountBalance(protoData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish account balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account balance created successfully"})
}
