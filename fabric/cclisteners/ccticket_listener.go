package cclisteners

import (
	"context"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/ccclient"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/consts"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
)

var validTicketCCNotificationTypes = []string{""}

type TickenTicketListener struct {
	listener *ccclient.Listener
	callback func(notification *ccclient.CCNotification)
}

func NewTickenTicketListener(pc peerconnector.PeerConnector, channel string) (*TickenTicketListener, error) {
	ticketListener := new(TickenTicketListener)

	listener, err := ccclient.NewListener(pc, channel, consts.TickenTicketChaincode)
	if err != nil {
		return nil, err
	}

	ticketListener.listener = listener
	ticketListener.callback = nil

	return ticketListener, nil
}

func (ticketListener *TickenTicketListener) ListenCCTicket(ctx context.Context, callback func(notification *ccclient.CCNotification)) {
	ticketListener.callback = callback

	internalCallback := func(notification *ccclient.CCNotification) {
		if ticketCCNotificationTypeIsValid(notification.Type) {
			ticketListener.callback(notification)
		}
	}

	ticketListener.listener.Listen(ctx, internalCallback)
}

func ticketCCNotificationTypeIsValid(typeToCheck string) bool {
	for _, notificationType := range validTicketCCNotificationTypes {
		if typeToCheck == notificationType {
			return true
		}
	}
	return false
}
