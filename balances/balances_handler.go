package balances

import (
	"fmt"
	"go_task/messaging"
	"golang.org/x/net/context"
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
	// ... ваш код инициализации, подключения к брокеру сообщений и т.д.

	for {
		select {
		case <-ctx.Done():
			// Контекст отменен, завершаем слушание
			return nil
		case delivery := <-h.listeningChannel.ConsumeAccountBalances():
			// Обработка сообщения delivery
			if err := h.processAccountBalance(delivery.Body); err != nil {
				// В случае ошибки в обработке сообщения, возвращаем ошибку
				return err
			}
		}
	}
}

func (h *BalancesHandler) processAccountBalance(protoData []byte) error {
	// Логика обработки сообщения
	// В случае ошибки:
	return fmt.Errorf("error processing account balance: %s", errorMessage)
}
