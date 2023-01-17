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

func (caller *TickenEventCaller) CreateEventAsync(eventID uuid.UUID, name string, date time.Time) error {
	_, err := caller.submiter.SubmitAsync(
		consts.EventCCCreateFunction,
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
	_, err := caller.submiter.SubmitAsync(
		consts.EventCCAddSectionFunction,
		eventID.String(),
		name,
		strconv.Itoa(totalTickets),
		fmt.Sprintf("%f", ticketPrice),
	)

	if err != nil {
		return err
	}

	return nil
}

func (caller *TickenEventCaller) GetEvent(eventID uuid.UUID) (*chainmodels.Event, error) {
	eventData, err := caller.querier.Query(consts.EventCCGetFunction, eventID.String())
	if err != nil {
		return nil, err
	}

	event := new(chainmodels.Event)

	err = json.Unmarshal(eventData, &event)
	if err != nil {
		return nil, err
	}

	return event, nil
}
