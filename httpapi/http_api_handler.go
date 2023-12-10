package httpapi

import (
	"api_http_component/constants"
	"api_http_component/messaging"
	"api_http_component/protofile"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
	"net/http"
)

type HTTPHandler struct {
	publishingChannel messaging.PublishingChannel
	listeningChannel  messaging.ListeningChannel
	messageQueue      chan Message
}

type Message struct {
	QueueType string
	Delivery  amqp.Delivery
}

func NewHTTPHandler(publishingChannel messaging.PublishingChannel, listeningChannel messaging.ListeningChannel) *HTTPHandler {
	return &HTTPHandler{
		publishingChannel: publishingChannel,
		listeningChannel:  listeningChannel,
		messageQueue:      make(chan Message),
	}
}

func (h *HTTPHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/createAccountBalance", h.CreateAccountBalances)
}

func (h *HTTPHandler) StartListener(ctx context.Context, ready chan<- struct{}) {
	go func() {
		defer close(ready)

		stopCh := make(chan struct{})
		defer close(stopCh)

		createAccountMessages := h.listeningChannel.ConsumeCreateAccountFromBalances(stopCh)

		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case delivery, ok := <-createAccountMessages:
					if !ok {
						return
					}
					h.messageQueue <- Message{QueueType: "CreateAccount", Delivery: delivery}
				}
			}
		}()

		select {
		case <-ctx.Done():
			close(stopCh)
			return
		}
	}()
}

func (h *HTTPHandler) CreateAccountBalances(c *gin.Context) {
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

	err = h.publishingChannel.PublishCreateAccountToBalances(protoData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish account balance"})
		return
	}

	select {
	case delivery := <-h.messageQueue:
		switch delivery.QueueType {
		case "CreateAccount":

		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Account created successfully"})
}

func (h *HTTPHandler) handleListenerError(err error, c *gin.Context) {
	c.JSON(http.StatusBadRequest, constants.ErrorResponse{Error: err.Error()})
}

func (h *HTTPHandler) processCreateAccount(protoData []byte) error {
	var errorResponse constants.ErrorResponse
	if err := json.Unmarshal(protoData, &errorResponse); err != nil {
		return err
	}
	return nil
}
