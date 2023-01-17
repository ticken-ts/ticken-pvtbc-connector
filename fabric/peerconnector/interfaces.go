package peerconnector

import (
	"context"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type PeerConnector interface {
	IsConnected() bool
	Connect(peerEndpoint string, gatewayPeer string, tlsCertPath string) error
	ConnectWithRawTlsCert(peerEndpoint string, gatewayPeer string, tlsCert []byte) error
	GetChaincode(channelName string, chaincodeName string) (Chaincode, error)

	GetChaincodeNotificationsChannel(
		ctx context.Context, channelName string, chaincodeName string,
	) (<-chan *client.ChaincodeEvent, error)
}

type Chaincode interface {
	ChaincodeName() string
	SubmitTx(name string, args ...string) ([]byte, error)
	EvaluateTx(name string, args ...string) ([]byte, error)
	SubmitTxAsync(name string, args ...string) ([]byte, error)
}
