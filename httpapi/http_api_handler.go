package httpapi

import (
	"api_http_component/constants"
	"api_http_component/messaging"
	"api_http_component/protofile"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"net/http"
)

type HTTPHandler struct {
	publishingChannel messaging.PublishingChannel
	listeningChannel  messaging.ListeningChannel
	createAccountCB   func(ctx context.Context, protoData []byte)
}

func NewHTTPHandler(publishingChannel messaging.PublishingChannel, listeningChannel messaging.ListeningChannel) *HTTPHandler {
	return &HTTPHandler{
		publishingChannel: publishingChannel,
		listeningChannel:  listeningChannel,
	}
}

func (h *HTTPHandler) SetCreateAccountCallback(cb func(ctx context.Context, protoData []byte)) {
	h.createAccountCB = cb
}

func (h *HTTPHandler) StartListener(ctx context.Context, ready chan<- struct{}) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case delivery := <-h.listeningChannel.ConsumeCreateAccountFromBalances():
				go h.createAccountCB(ctx, delivery.Body)
			}
		}
	}()
	close(ready)
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

	go func(protoData []byte, c *gin.Context) {
		if err := h.processCreateAccount(protoData); err != nil {
			h.handleListenerError(err, c)
		}
	}(protoData, c)

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
