package peerconnector

import "github.com/hyperledger/fabric-gateway/pkg/client"

type CoreChaincodeAPI struct {
	contract *client.Contract
}

func NewCoreChaincodeAPI(contract *client.Contract) Chaincode {
	return &CoreChaincodeAPI{contract: contract}
}

func (cc CoreChaincodeAPI) ChaincodeName() string {
	return cc.ChaincodeName()
}

func (cc CoreChaincodeAPI) SubmitTx(name string, args ...string) ([]byte, error) {
	return cc.contract.SubmitTransaction(name, args...)
}

func (cc CoreChaincodeAPI) EvaluateTx(name string, args ...string) ([]byte, error) {
	return cc.contract.EvaluateTransaction(name, args...)
}

func (cc CoreChaincodeAPI) SubmitTxAsync(name string, args ...string) ([]byte, error) {
	// for now, we are going to ignore the commit
	// that we receive to simplify the API of the
	// chaincode.
	result, _, err := cc.contract.SubmitAsync(name, client.WithArguments(args...))
	return result, err
}
