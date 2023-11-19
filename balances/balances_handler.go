package balances

import (
	"github.com/gin-gonic/gin"
	"go_task/messaging"
	"golang.org/x/net/context"
	"net/http"
)

type BalancesHandler struct {
	listeningChannel messaging.ListeningChannel
}

func NewBalancesHandler(listeningChannel messaging.ListeningChannel) *BalancesHandler {
	return &BalancesHandler{
		listeningChannel: listeningChannel,
	}
}

func (h *BalancesHandler) StartListener(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case delivery := <-h.listeningChannel.ConsumeCreateAccountBalances():
			if err := h.processAccountBalance(delivery.Body); err != nil {
				errorHandler.HandleError(http.StatusInternalServerError, gin.H{"error": "Failed to process account balance"})
				return err
			}
		}
	}
}

func (h *BalancesHandler) processAccountBalance(protoData []byte) error {
	// Логика обработки сообщения
	//if /* какое-то условие для ошибки */ {
	//	return &constants.CustomError{Type: constants.AccountNotHaveBalance}
	//}
	// ...
	return nil
}
