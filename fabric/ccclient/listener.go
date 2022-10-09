package ccclient

import (
	"context"
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
)

const initialEventsSize = 5

type Listener struct {
	pc           *peerconnector.PeerConnector
	chaincode    string
	channel      string
	eventsChanns []<-chan *client.ChaincodeEvent
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
	listener.pc = pc
	listener.channel = channelName
	listener.chaincode = chaincodeName
	listener.eventsChanns = []<-chan *client.ChaincodeEvent{}

	return listener, nil
}

func (listener *Listener) Listen(ctx context.Context, eventType string, callback func([]byte)) {
	eventsChann, err := listener.pc.GetChaincodeEvents(ctx, listener.channel, listener.chaincode)
	if err != nil {
		panic(err)
	}

	listener.eventsChanns = append(listener.eventsChanns, eventsChann)

	go func() {
		for event := range eventsChann {
			if event.EventName == eventType {
				callback(event.Payload)
			}
		}
	}()
}
