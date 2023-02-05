package cccallers

import (
	"encoding/json"
	"github.com/google/uuid"
	chainmodels "github.com/ticken-ts/ticken-pvtbc-connector/chain-models"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/ccclient"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/consts"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
)

type TickenTicketCaller struct {
	submiter *ccclient.Submiter
	querier  *ccclient.Querier
}

func NewTickenTicketCaller(pc peerconnector.PeerConnector, channelName string) (*TickenTicketCaller, error) {
	submiter, err := ccclient.NewSubmiter(pc, channelName, consts.TickenTicketChaincode)
	if err != nil {
		return nil, err
	}

	querier, err := ccclient.NewQuerier(pc, channelName, consts.TickenTicketChaincode)
	if err != nil {
		return nil, err
	}

	caller := new(TickenTicketCaller)
	caller.submiter = submiter
	caller.querier = querier

	return caller, nil
}

func (caller *TickenTicketCaller) IssueTicket(ticketID, eventID, owner uuid.UUID, section string) (*chainmodels.Ticket, error) {
	function := consts.TicketCCIssueFunction
	data, _, err := caller.submiter.Submit(function, ticketID.String(), eventID.String(), section, owner.String())
	if err != nil {
		return nil, err
	}

	var ticket chainmodels.Ticket
	if err := json.Unmarshal(data, &ticket); err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (caller *TickenTicketCaller) GetTicket(ticketID uuid.UUID) (*chainmodels.Ticket, error) {
	function := consts.TicketCCGetTicketFunction
	data, _, err := caller.submiter.Submit(function, ticketID.String())
	if err != nil {
		return nil, err
	}

	var ticket chainmodels.Ticket
	if err := json.Unmarshal(data, &ticket); err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (caller *TickenTicketCaller) GetSectionTickets(section string) ([]*chainmodels.Ticket, error) {
	function := consts.TicketCCGetSectionTicketsFunction
	data, _, err := caller.submiter.Submit(function, section)
	if err != nil {
		return nil, err
	}

	var sectionTickets []*chainmodels.Ticket
	if err := json.Unmarshal(data, &sectionTickets); err != nil {
		return nil, err
	}

	return sectionTickets, nil
}
