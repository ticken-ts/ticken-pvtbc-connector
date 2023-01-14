package cclisteners

import (
	"context"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/ccclient"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/config"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
)

type TickenTicketListener struct {
	listener *ccclient.Listener
	callback func(notification *CCTicketNotification)
}

type CCTicketNotificationType string

type CCTicketNotification struct {
	BlockNum uint64
	TxID     string
	Type     CCTicketNotificationType
	Payload  []byte
}

func NewTickenTicketListener(pc *peerconnector.PeerConnector, channel string) (*TickenTicketListener, error) {
	ticketListener := new(TickenTicketListener)

	listener, err := ccclient.NewListener(pc, channel, config.TickenTicketChaincode)
	if err != nil {
		return nil, err
	}

	ticketListener.listener = listener
	ticketListener.callback = nil

	return ticketListener, nil
}

func (ticketListener *TickenTicketListener) ListenCCTicket(ctx context.Context, callback func(notification *CCTicketNotification)) {
	ticketListener.callback = callback

	internalCallback := func(notification *ccclient.CCNotification) {
		notificationType := stringToTicketNotificationType(notification.Type)

		// if we can not identify the notification type,
		// we just are going to skip the notification
		// processing
		if len(notificationType) == 0 {
			return
		}

		eventNotification := &CCTicketNotification{
			Type:     notificationType,
			TxID:     notification.TxID,
			BlockNum: notification.BlockNum,
			Payload:  notification.Payload,
		}

		ticketListener.callback(eventNotification)
	}

	ticketListener.listener.Listen(ctx, internalCallback)
}

func stringToTicketNotificationType(_ string) CCTicketNotificationType {
	return ""
}
