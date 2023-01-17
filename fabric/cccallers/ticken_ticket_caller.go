package cccallers

import (
	"encoding/json"
	"fmt"
	"github.com/ticken-ts/ticken-pvtbc-connector/chain-models"
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

func (caller *TickenTicketCaller) IssueTicket(ticketID string, eventID string, section string, owner string) (*chain_models.Ticket, error) {
	data, err := caller.submiter.Submit(consts.TicketCCIssueFunction, ticketID, eventID, section, owner)
	if err != nil {
		return nil, err
	}

	ticket := new(chain_models.Ticket)
	err = json.Unmarshal(data, &ticket)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (caller *TickenTicketCaller) SignTicket(ticketID string, eventID string, signer string, signature []byte) (*chain_models.Ticket, error) {
	hexSignature := fmt.Sprintf("%x", signature)

	data, err := caller.submiter.Submit(consts.TicketCCSignFunction, ticketID, eventID, signer, hexSignature)
	if err != nil {
		return nil, err
	}

	ticket := new(chain_models.Ticket)
	err = json.Unmarshal(data, &ticket)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (caller *TickenTicketCaller) ScanTicket(ticketID string, eventID string, owner string) (*chain_models.Ticket, error) {
	data, err := caller.submiter.Submit(consts.TicketCCScanFunction, ticketID, eventID, owner)
	if err != nil {
		return nil, err
	}

	ticket := new(chain_models.Ticket)
	err = json.Unmarshal(data, &ticket)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}
