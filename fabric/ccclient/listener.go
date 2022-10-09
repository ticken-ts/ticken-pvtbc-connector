package ccclient

import (
	"context"
	"fmt"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
)

type Listener struct {
	pc        *peerconnector.PeerConnector
	chaincode string
	channel   string
}

func NewListener(pc *peerconnector.PeerConnector, channelName string, chaincodeName string) (*Listener, error) {
	if !pc.IsConnected() {
		return nil, fmt.Errorf("connection with peer is not stablished")
	}

	// just to check if chaincode exists
	_, err := pc.GetChaincode(channelName, chaincodeName)
	if err != nil {
		return nil, err
	}

	listener := new(Listener)
	listener.chaincode = chaincodeName

	return listener, nil
}

func (listener *Listener) Listen(ctx context.Context, eventType string, callback func([]byte)) {
	events, err := listener.pc.GetChaincodeEvents(ctx, listener.channel, listener.chaincode)
	if err != nil {
		panic(err)
	}

	go func() {
		for event := range events {
			if event.EventName == eventType {
				callback(event.Payload)
			}
		}
	}()
}
