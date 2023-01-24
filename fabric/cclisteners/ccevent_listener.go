package cclisteners

import (
	"context"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/ccclient"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/consts"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
)

var validEventCCNotificationTypes = []string{
	consts.EventCreatedNotification,
	consts.SectionAddedNotification,
}

type TickenEventListener struct {
	listener *ccclient.Listener
	callback func(notification *ccclient.CCNotification)
}

func NewTickenEventListener(pc peerconnector.PeerConnector, channel string) (*TickenEventListener, error) {
	eventListener := new(TickenEventListener)

	listener, err := ccclient.NewListener(pc, channel, consts.TickenEventChaincode)
	if err != nil {
		return nil, err
	}

	eventListener.listener = listener
	eventListener.callback = nil

	return eventListener, nil

}

func (eventListener *TickenEventListener) ListenCCEvent(ctx context.Context, callback func(notification *ccclient.CCNotification)) {
	eventListener.callback = callback

	internalCallback := func(notification *ccclient.CCNotification) {
		if eventNotificationTypeIsValid(notification.Type) {
			eventListener.callback(notification)
		}
	}

	eventListener.listener.Listen(ctx, internalCallback)
}

func eventNotificationTypeIsValid(typeToCheck string) bool {
	for _, notificationType := range validEventCCNotificationTypes {
		if typeToCheck == notificationType {
			return true
		}
	}
	return false
}
