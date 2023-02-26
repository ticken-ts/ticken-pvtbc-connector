package cccallers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	chainmodels "github.com/ticken-ts/ticken-pvtbc-connector/chain-models"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/ccclient"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/consts"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
	"strconv"
	"time"
)

type TickenEventCaller struct {
	submiter *ccclient.Submiter
	querier  *ccclient.Querier
}

func NewTickenEventCaller(pc peerconnector.PeerConnector, channelName string) (*TickenEventCaller, error) {
	submiter, err := ccclient.NewSubmiter(pc, channelName, consts.TickenEventChaincode)
	if err != nil {
		return nil, err
	}

	querier, err := ccclient.NewQuerier(pc, channelName, consts.TickenEventChaincode)
	if err != nil {
		return nil, err
	}

	caller := new(TickenEventCaller)
	caller.submiter = submiter
	caller.querier = querier

	return caller, nil
}

func (caller *TickenEventCaller) CreateEvent(eventID uuid.UUID, name string, date time.Time) (*chainmodels.Event, string, error) {
	function := consts.EventCCCreateFunction
	payload, txID, err := caller.submiter.Submit(function, eventID.String(), name, date.Format(time.RFC3339))
	if err != nil {
		return nil, "", err
	}

	event := new(chainmodels.Event)
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, txID, err
	}

	return event, txID, nil
}

func (caller *TickenEventCaller) SetEventOnSale(eventID uuid.UUID) (string, error) {
	function := consts.EventCCSetEventOnSaleFunction
	_, txID, err := caller.submiter.Submit(function, eventID.String())
	return txID, err
}

func (caller *TickenEventCaller) AddSection(eventID uuid.UUID, name string, totalTickets int, ticketPrice float64) (*chainmodels.Section, string, error) {
	payload, txID, err := caller.submiter.Submit(
		consts.EventCCAddSectionFunction,
		eventID.String(),
		name,
		strconv.Itoa(totalTickets),
		fmt.Sprintf("%f", ticketPrice),
	)
	if err != nil {
		return nil, "", err
	}

	section := new(chainmodels.Section)
	if err := json.Unmarshal(payload, &section); err != nil {
		return nil, "", err
	}

	return section, txID, nil
}

func (caller *TickenEventCaller) GetEvent(eventID uuid.UUID) (*chainmodels.Event, error) {
	eventData, err := caller.querier.Query(consts.EventCCGetEventFunction, eventID.String())
	if err != nil {
		return nil, err
	}

	event := new(chainmodels.Event)
	if err := json.Unmarshal(eventData, &event); err != nil {
		return nil, err
	}

	return event, nil
}
