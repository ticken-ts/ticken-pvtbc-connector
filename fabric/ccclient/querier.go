package ccclient

import (
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
)

type Querier struct {
	chaincode *client.Contract
}

func NewQuerier(pc *peerconnector.PeerConnector, channelName string, chaincodeName string) (*Querier, error) {
	if !pc.IsConnected() {
		return nil, fmt.Errorf("connection with peer is not stablished")
	}

	chaincode, err := pc.GetChaincode(channelName, chaincodeName)
	if err != nil {
		return nil, err
	}

	chaincodeQuerier := new(Querier)
	chaincodeQuerier.chaincode = chaincode

	return chaincodeQuerier, nil
}

func (querier *Querier) Query(function string, args ...string) ([]byte, error) {
	evaluateResult, err := querier.chaincode.EvaluateTransaction(function, args...)

	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	return evaluateResult, nil
}
