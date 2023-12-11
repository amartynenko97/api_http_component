package httpapi

import (
	"api_http_component/constants"
	"api_http_component/messaging"
	"api_http_component/protofile"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
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
					h.messageQueue <- Message{QueueType: constants.QueueTypeCreateAccount, Delivery: delivery}
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
		case constants.QueueTypeCreateAccount:
			balanceErrorMessage := &protofile.BalanceErrorMessage{}
			err := proto.Unmarshal(delivery.Delivery.Body, balanceErrorMessage)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deserialize balance error message"})
				return
			}

			switch balanceErrorMessage.GetErrorCode() {
			case protofile.BalancesErrorCodes_BALANCE_ERROR_CODE_INTERNAL:
				c.JSON(http.StatusInternalServerError, gin.H{"error": balanceErrorMessage.Message})
			default:
				c.JSON(http.StatusOK, gin.H{"message": "Account created successfully"})
			}
		}
	}
}
