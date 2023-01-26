package peerconnector

import (
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
)

type CoreChaincodeAPI struct {
	contract *client.Contract
}

func NewCoreChaincodeAPI(contract *client.Contract) Chaincode {
	return &CoreChaincodeAPI{contract: contract}
}

func (cc CoreChaincodeAPI) ChaincodeName() string {
	return cc.contract.ContractName()
}

func (cc CoreChaincodeAPI) SubmitTx(name string, args ...string) ([]byte, string, error) {
	result, commit, err := cc.contract.SubmitAsync(name, client.WithArguments(args...))
	if err != nil {
		return result, "", err
	}

	status, err := commit.Status()
	if err != nil {
		return result, "", err
	}

	if !status.Successful {
		return nil, "", txError(status.TransactionID, status.Code)
	}

	return result, commit.TransactionID(), nil
}

func (cc CoreChaincodeAPI) EvaluateTx(name string, args ...string) ([]byte, error) {
	return cc.contract.EvaluateTransaction(name, args...)
}

func (cc CoreChaincodeAPI) SubmitTxAsync(name string, args ...string) ([]byte, string, error) {
	result, commit, err := cc.contract.SubmitAsync(name, client.WithArguments(args...))
	return result, commit.TransactionID(), err
}

func txError(transactionID string, code peer.TxValidationCode) error {
	return fmt.Errorf(
		"transaction %s failed to commit with status code %d (%s)",
		transactionID, int32(code), peer.TxValidationCode_name[int32(code)])
}
