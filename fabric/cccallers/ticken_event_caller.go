package cccallers

import (
	"encoding/json"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/ccclient"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
	"github.com/ticken-ts/ticken-pvtbc-connector/onchain-models"
)

const TickenEventChaincode = "ticken-event"

const (
	EventCCGetFunction = "Get"
)

type TickenEventCaller struct {
	submiter *ccclient.Submiter
	querier  *ccclient.Querier
}

func NewTickenEventCaller(pc *peerconnector.PeerConnector, channelName string) (*TickenEventCaller, error) {
	submiter, err := ccclient.NewSubmiter(pc, channelName, TickenEventChaincode)
	if err != nil {
		return nil, err
	}

	querier, err := ccclient.NewQuerier(pc, channelName, TickenEventChaincode)
	if err != nil {
		return nil, err
	}

	caller := new(TickenEventCaller)
	caller.submiter = submiter
	caller.querier = querier

	return caller, nil
}

func (caller *TickenEventCaller) GetEvent(eventID string) (*onchain_models.Event, error) {
	eventData, err := caller.querier.Query(EventCCGetFunction, eventID)
	if err != nil {
		return nil, err
	}

	event := new(onchain_models.Event)

	err = json.Unmarshal(eventData, &event)
	if err != nil {
		return nil, err
	}

	return event, nil
}
