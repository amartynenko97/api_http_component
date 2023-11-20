package httpapi

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go_task/constants"
	"go_task/messaging"
	"go_task/protofile"
	"golang.org/x/net/context"
	"net/http"
)

type HTTPHandler struct {
	publishingChannel messaging.PublishingChannel
	listeningChannel  messaging.ListeningChannel
}

func NewHTTPHandler(publishingChannel messaging.PublishingChannel, listeningChannel messaging.ListeningChannel) *HTTPHandler {
	return &HTTPHandler{
		publishingChannel: publishingChannel,
		listeningChannel:  listeningChannel,
	}
}

func (h *HTTPHandler) StartListener(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case delivery := <-h.listeningChannel.ConsumeCreateAccountBalancesFromHttpApi():
				go func(protoData []byte, c *gin.Context) {
					if err := h.processCreateAccountBalance(protoData); err != nil {
						h.handleListenerError(err, c)
					}
				}(delivery.Body, delivery.Context)
			}
		}
	}()
}

func (h *HTTPHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/createAccountBalance", h.CreateAccountBalance)
	//router.GET("/createOrder", h.GetAccountBalance)
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

	err = h.publishingChannel.PublishCreateAccountBalancesToBalances(protoData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish account balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account balance created successfully"})
}

func (h *HTTPHandler) handleListenerError(err error, c *gin.Context) {
	c.JSON(http.StatusBadRequest, constants.ErrorResponse{Error: err.Error()})
}

func (h *HTTPHandler) processCreateAccountBalance(protoData []byte) error {
	var errorResponse constants.ErrorResponse
	if err := json.Unmarshal(protoData, &errorResponse); err != nil {
		return err
	}
	return nil
}

//func (h *HTTPHandler) GetAccountBalance(c *gin.Context) {
//	var request protofile.CreateOrderRequest
//
//	if err := c.ShouldBindJSON(&request); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	protoData, err := json.Marshal(&request)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal create order message"})
//		return
//	}
//
//	err = h.publishingChannel.PublishGetAccountBalances(protoData)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish create order"})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "Create order successfully"})
//}
