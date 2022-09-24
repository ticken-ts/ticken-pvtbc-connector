package ccclient

import (
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
)

type Submiter struct {
	chaincode *client.Contract
}

func NewSubmiter(pc *peerconnector.PeerConnector, channelName string, chaincodeName string) (*Submiter, error) {
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

func (submiter *Submiter) Submit(function string, args ...string) ([]byte, error) {
	evaluateResult, err := submiter.chaincode.SubmitTransaction(function, args...)
	if err != nil {
		return nil, err
	}

	return evaluateResult, nil
}

func (submiter *Submiter) SubmitAsync(function string, args ...string) ([]byte, *client.Commit, error) {
	submitResult, commit, err := submiter.chaincode.SubmitAsync(function, client.WithArguments(args...))
	if err != nil {
		return nil, nil, err
	}

	return submitResult, commit, nil
}
