package peerconnector

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// these variables should be global because it must be shared between
// all the instances of the dev connector to simulate the transactions
// comes from the same connection
var (
	storedElements          = make(map[string]map[uuid.UUID][]byte)
	fakeNotificationChannel = make(chan *client.ChaincodeEvent)
)

type DevPeerConnector struct {
	isConnected bool

	// mocked organization identities
	mspID             string
	organizerUsername string

	// internal mocked storage
	storedElements           map[string]map[uuid.UUID][]byte
	fakeNotificationsChannel chan *client.ChaincodeEvent
}

func NewDev(mspID string, organizerUsername string) PeerConnector {
	for _, elementName := range []string{eventElementName, ticketElementName} {
		_, elementStorageExists := storedElements[elementName]
		if !elementStorageExists {
			storedElements[elementName] = make(map[uuid.UUID][]byte)
		}
	}

	return &DevPeerConnector{
		isConnected: false,

		mspID:             mspID,
		organizerUsername: organizerUsername,

		storedElements:           storedElements,
		fakeNotificationsChannel: fakeNotificationChannel,
	}
}

func (hfc *DevPeerConnector) IsConnected() bool {
	return hfc.isConnected
}

func (hfc *DevPeerConnector) Connect(_ string, _ string, _ string) error {
	if hfc.IsConnected() {
		return fmt.Errorf("gateway is already connected")
	}
	hfc.isConnected = true
	return nil
}

func (hfc *DevPeerConnector) ConnectWithRawTlsCert(_ string, _ string, _ []byte) error {
	if hfc.IsConnected() {
		return fmt.Errorf("gateway is already connected")
	}
	hfc.isConnected = true
	return nil
}

func (hfc *DevPeerConnector) GetChaincode(channelName string, chaincodeName string) (Chaincode, error) {
	devChaincode := &DevChaincodeAPI{
		name:    chaincodeName,
		channel: channelName,

		ctxMSPID:             hfc.mspID,
		ctxOrganizerUsername: hfc.organizerUsername,

		storedElements:           hfc.storedElements,
		fakeNotificationsChannel: hfc.fakeNotificationsChannel,
	}

	return devChaincode, nil
}

func (hfc *DevPeerConnector) GetChaincodeNotificationsChannel(
	_ context.Context, _ string, _ string,
) (<-chan *client.ChaincodeEvent, error) {
	return hfc.fakeNotificationsChannel, nil
}
