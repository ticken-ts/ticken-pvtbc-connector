package ccclient

import (
	"fmt"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
)

type Submiter struct {
	chaincode peerconnector.Chaincode
}

func NewSubmiter(pc peerconnector.PeerConnector, channelName string, chaincodeName string) (*Submiter, error) {
	if !pc.IsConnected() {
		return nil, fmt.Errorf("connection with peer is not stablished")
	}

	chaincode, err := pc.GetChaincode(channelName, chaincodeName)
	if err != nil {
		return nil, err
	}

	chaincodeSubmiter := new(Submiter)
	chaincodeSubmiter.chaincode = chaincode

	return chaincodeSubmiter, nil
}

func (submiter *Submiter) Submit(function string, args ...string) ([]byte, string, error) {
	evaluateResult, txID, err := submiter.chaincode.SubmitTx(function, args...)
	if err != nil {
		return nil, "", err
	}

	return evaluateResult, txID, nil
}

func (submiter *Submiter) SubmitAsync(function string, args ...string) ([]byte, string, error) {
	submitResult, txID, err := submiter.chaincode.SubmitTxAsync(function, args...)
	if err != nil {
		return nil, "", err
	}

	return submitResult, txID, nil
}
