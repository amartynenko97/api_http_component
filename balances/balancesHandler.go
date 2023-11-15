package balances

import (
	"go_task/messaging"
)

type BalancesHandler struct {
	listeningChannel messaging.ListeningChannel
}

func NewBalancesHandler(listeningChannel messaging.ListeningChannel) *BalancesHandler {
	return &BalancesHandler{
		listeningChannel: listeningChannel,
	}
}

func (b *BalancesHandler) HandleAccountBalance() {

}
