package cccallers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	chain_models "github.com/ticken-ts/ticken-pvtbc-connector/chain-models"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/ccclient"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
	"strconv"
	"time"
)

const TickenEventChaincode = "cc-event"

const (
	EventCCGetFunction        = "Get"
	EventCCCreateFunction     = "Create"
	EventCCAddSectionFunction = "AddSection"
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

func (caller *TickenEventCaller) CreateAsync(eventID uuid.UUID, name string, date time.Time) error {
	_, _, err := caller.submiter.SubmitAsync(
		EventCCCreateFunction,
		eventID.String(),
		name,
		date.Format(time.RFC3339),
	)

	if err != nil {
		return err
	}

	return nil
}

func (caller *TickenEventCaller) AddSectionAsync(eventID uuid.UUID, name string, totalTickets int, ticketPrice float64) error {
	_, _, err := caller.submiter.SubmitAsync(
		EventCCAddSectionFunction,
		eventID,
		name,
		strconv.Itoa(totalTickets),
		fmt.Sprintf("%f", ticketPrice),
	)

	if err != nil {
		return err
	}

	return nil
}

func (caller *TickenEventCaller) GetEvent(eventID uuid.UUID) (*chain_models.Event, error) {
	eventData, err := caller.querier.Query(EventCCGetFunction, eventID)
	if err != nil {
		return nil, err
	}

	event := new(chain_models.Event)

	err = json.Unmarshal(eventData, &event)
	if err != nil {
		return nil, err
	}

	return event, nil
}
