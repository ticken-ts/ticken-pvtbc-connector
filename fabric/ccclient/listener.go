package ccclient

import (
	"context"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
)

type Listener struct {
	events <-chan *client.ChaincodeEvent
}

func NewListener(ctx context.Context, pc *peerconnector.PeerConnector, channelName string, chaincodeName string) (*Listener, error) {
	listener := new(Listener)

	events, err := pc.GetChaincodeEvents(ctx, channelName, chaincodeName)
	if err != nil {
		return nil, err
	}

	listener.events = events

	return listener, nil
}

func (listener *Listener) Listen(eventType string, callback func([]byte)) {
	go func() {
		for event := range listener.events {
			if event.EventName == eventType {
				callback(event.Payload)
			}
		}
	}()
}
