package cccallers

import (
	"encoding/json"
	"fmt"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/ccclient"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
	"github.com/ticken-ts/ticken-pvtbc-connector/onchain-models"
)

const TickenTicketChaincode = "ticken-ticket"

const (
	TicketCCIssueFunction = "Issue"
	TicketCCSignFunction  = "Sign"
	TicketCCScanFunction  = "Scan"
)

type TickenTicketCaller struct {
	submiter *ccclient.Submiter
	querier  *ccclient.Querier
}

func NewTickenTicketCaller(pc *peerconnector.PeerConnector, channelName string) (*TickenTicketCaller, error) {
	submiter, err := ccclient.NewSubmiter(pc, channelName, TickenTicketChaincode)
	if err != nil {
		return nil, err
	}

	querier, err := ccclient.NewQuerier(pc, channelName, TickenTicketChaincode)
	if err != nil {
		return nil, err
	}

	caller := new(TickenTicketCaller)
	caller.submiter = submiter
	caller.querier = querier

	return caller, nil
}

func (caller *TickenTicketCaller) IssueTicket(ticketID string, eventID string, section string, owner string) (*onchain_models.Ticket, error) {
	data, err := caller.submiter.Submit(TicketCCIssueFunction, ticketID, eventID, section, owner)
	if err != nil {
		return nil, err
	}

	ticket := new(onchain_models.Ticket)
	err = json.Unmarshal(data, &ticket)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (caller *TickenTicketCaller) SignTicket(ticketID string, eventID string, signer string, signature []byte) (*onchain_models.Ticket, error) {
	hexSignature := fmt.Sprintf("%x", signature)

	data, err := caller.submiter.Submit(TicketCCSignFunction, ticketID, eventID, signer, hexSignature)
	if err != nil {
		return nil, err
	}

	ticket := new(onchain_models.Ticket)
	err = json.Unmarshal(data, &ticket)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (caller *TickenTicketCaller) ScanTicket(ticketID string, eventID string, owner string) (*onchain_models.Ticket, error) {
	data, err := caller.submiter.Submit(TicketCCScanFunction, ticketID, eventID, owner)
	if err != nil {
		return nil, err
	}

	ticket := new(onchain_models.Ticket)
	err = json.Unmarshal(data, &ticket)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}
