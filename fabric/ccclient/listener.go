package ccclient

import (
	"context"
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
)

type CCNotification struct {
	BlockNum uint64
	TxID     string
	Type     string
	Payload  []byte
}

type Listener struct {
	pc            peerconnector.PeerConnector
	chaincode     string
	channel       string
	notifications <-chan *client.ChaincodeEvent
}

func NewListener(pc peerconnector.PeerConnector, channelName string, chaincodeName string) (*Listener, error) {
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

	listener.notifications = nil
	listener.channel = channelName
	listener.chaincode = chaincodeName

	return listener, nil
}

func (listener *Listener) Listen(ctx context.Context, callback func(notification *CCNotification)) {
	notificationsChannel, err := listener.pc.GetChaincodeNotificationsChannel(
		ctx,
		listener.channel,
		listener.chaincode,
	)

	if err != nil {
		panic(err)
	}

	listener.notifications = notificationsChannel

	go func() {
		for notification := range notificationsChannel {

			ccnotification := &CCNotification{
				Type:     notification.EventName,
				TxID:     notification.TransactionID,
				BlockNum: notification.BlockNumber,
				Payload:  notification.Payload,
			}

			callback(ccnotification)
		}
	}()
}
