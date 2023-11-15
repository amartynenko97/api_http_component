package httpapi

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go_task/messaging"
	"go_task/protofile"
	"net/http"
)

type HTTPHandler struct {
	publishingChannel messaging.PublishingChannel
}

func NewHTTPHandler(publishingChannel messaging.PublishingChannel) *HTTPHandler {
	return &HTTPHandler{
		publishingChannel: publishingChannel,
	}
}

func (h *HTTPHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/createAccountBalance", h.CreateAccountBalance)
	router.POST("/createOrder", h.CreateOrder)
}

func (h *HTTPHandler) CreateAccountBalance(c *gin.Context) {
	var request protofile.CreateOrderRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	protoData, err := json.Marshal(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal account balance message"})
		return
	}

	err = h.publishingChannel.PublishAccountBalance(protoData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish account balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account balance created successfully"})
}

func (h *HTTPHandler) CreateOrder(c *gin.Context) {
	var request protofile.CreateOrderRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	protoData, err := json.Marshal(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal create order message"})
		return
	}

	err = h.publishingChannel.PublishCreateOrder(protoData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish create order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Create order successfully"})
}
