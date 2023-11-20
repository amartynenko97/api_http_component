package balances

import (
	"encoding/json"
	"go_task/constants"
	"go_task/messaging"
	"golang.org/x/net/context"
)

type BalancesHandler struct {
	publishingChannel messaging.PublishingChannel
	listeningChannel  messaging.ListeningChannel
}

func NewBalancesHandler(publishingChannel messaging.PublishingChannel, listeningChannel messaging.ListeningChannel) *BalancesHandler {
	return &BalancesHandler{
		publishingChannel: publishingChannel,
		listeningChannel:  listeningChannel,
	}
}

func (h *BalancesHandler) StartListener(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case delivery := <-h.listeningChannel.ConsumeCreateAccountBalancesFromHttpApi():
			if err := h.processAccountBalance(delivery.Body); err != nil {
				return err
			}
		}
	}
}

func (h *BalancesHandler) processAccountBalance(protoData []byte) error {

	errorResponse := constants.ErrorResponse{
		Error: string(constants.NoSuchCurrency),
	}

	errorJSON, err := json.Marshal(errorResponse)
	if err != nil {
		return err
	}

	err = h.publishingChannel.PublishCreateAccountBalancesToHttpApi(errorJSON)
	if err != nil {
		return err
	}

	return nil
}
