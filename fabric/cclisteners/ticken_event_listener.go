package cclisteners

import (
	"context"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/ccclient"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/config"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
)

type TickenEventListener struct {
	listener *ccclient.Listener
	callback func(notification *CCEventNotification)
}

type CCEventNotificationType string

const (
	EventCreatedNotification CCEventNotificationType = "event-created"
	SectionAddedNotification CCEventNotificationType = "section-added"
)

type CCEventNotification struct {
	BlockNum uint64
	TxID     string
	Type     CCEventNotificationType
	Payload  []byte
}

func NewTickenEventListener(pc *peerconnector.PeerConnector, channel string) (*TickenEventListener, error) {
	eventListener := new(TickenEventListener)

	listener, err := ccclient.NewListener(pc, channel, config.TickenEventChaincode)
	if err != nil {
		return nil, err
	}

	eventListener.listener = listener
	eventListener.callback = nil

	return eventListener, nil

}

func (eventListener *TickenEventListener) ListenCCEvent(ctx context.Context, callback func(notification *CCEventNotification)) {
	eventListener.callback = callback

	internalCallback := func(notification *ccclient.CCNotification) {
		notificationType := stringToEventNotificationType(notification.Type)

		// if we can not identify the notification type,
		// we just are going to skip the notification
		// processing
		if len(notificationType) == 0 {
			return
		}

		eventNotification := &CCEventNotification{
			Type:     notificationType,
			TxID:     notification.TxID,
			BlockNum: notification.BlockNum,
			Payload:  notification.Payload,
		}

		eventListener.callback(eventNotification)
	}

	eventListener.listener.Listen(ctx, internalCallback)
}

func stringToEventNotificationType(s string) CCEventNotificationType {
	if s == string(EventCreatedNotification) {
		return EventCreatedNotification
	}

	if s == string(SectionAddedNotification) {
		return SectionAddedNotification
	}

	return ""
}
