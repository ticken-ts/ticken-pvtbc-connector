package peerconnector

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	chainmodels "github.com/ticken-ts/ticken-pvtbc-connector/chain-models"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/consts"
	"math/big"
	"strconv"
	"time"
)

const (
	eventElementName  = "event"
	ticketElementName = "ticket"
)

type DevChaincodeAPI struct {
	name    string
	channel string

	// these are the context identity
	// which is requesting the chaincode
	// to submit/evaluate a transaction
	ctxMSPID             string
	ctxOrganizerUsername string

	// the first key is the element type. e.g: ticket/event
	// the second key represents the id of the element
	storedElements           map[string]map[uuid.UUID][]byte
	fakeNotificationsChannel chan *client.ChaincodeEvent
}

func (cc DevChaincodeAPI) ChaincodeName() string {
	return cc.name
}

func (cc DevChaincodeAPI) SubmitTx(name string, args ...string) ([]byte, string, error) {
	var payload []byte
	var err error
	var elemKey uuid.UUID
	var notification *client.ChaincodeEvent
	var elementName string

	switch cc.ChaincodeName() {
	case consts.TickenTicketChaincode:
		elementName = ticketElementName
		payload, elemKey, notification, err = cc.handleTicketCCAPI(name, args...)
	case consts.TickenEventChaincode:
		elementName = eventElementName
		payload, elemKey, notification, err = cc.handleEventCCAPI(name, args...)
	default:
		return nil, "", fmt.Errorf("chaincode %s not exists", cc.ChaincodeName())
	}

	if err != nil {
		return nil, "", err
	}

	cc.storedElements[elementName][elemKey] = payload
	if notification != nil {
		// avoid blocking when sending notification on the channel
		go func() { cc.fakeNotificationsChannel <- notification }()
	}

	fakeTxID := uuid.New().String()
	return payload, fakeTxID, nil
}

func (cc DevChaincodeAPI) EvaluateTx(name string, args ...string) ([]byte, error) {
	var payload []byte
	var err error

	switch cc.ChaincodeName() {
	case consts.TickenTicketChaincode:
		payload, _, _, err = cc.handleTicketCCAPI(name, args...)
	case consts.TickenEventChaincode:
		payload, _, _, err = cc.handleEventCCAPI(name, args...)
	default:
		return nil, fmt.Errorf("chaincode %s not exists", cc.ChaincodeName())
	}

	return payload, err
}

func (cc DevChaincodeAPI) SubmitTxAsync(name string, args ...string) ([]byte, string, error) {
	return cc.SubmitTx(name, args...)
}

// handleEventCCAPI decides based on the function name, how to mock the
// operation, parsing the args correct and generating a similar output
// based on what the real ccevent chaincode should respond
func (cc DevChaincodeAPI) handleEventCCAPI(name string, args ...string) ([]byte, uuid.UUID, *client.ChaincodeEvent, error) {
	switch name {
	case consts.EventCCCreateFunction:
		return cc.handleEventCCCreateTx(args...)
	case consts.EventCCAddSectionFunction:
		return cc.handleEventCCAddSectionTx(args...)
	case consts.EventCCGetEventFunction:
		return cc.handleEventCCGetEventTx(args...)
	case consts.EventCCSetEventOnSaleFunction:
		return cc.handleEventCCSetEventOnSaleTx(args...)

	default:
		return nil, uuid.Nil, nil, fmt.Errorf("function not found")
	}
}

// handleTicketCCAPI decides based on the function name, how to mock the
// operation, parsing the args correct and generating a similar output
// based on what the real ccticket chaincode should respond
func (cc DevChaincodeAPI) handleTicketCCAPI(name string, args ...string) ([]byte, uuid.UUID, *client.ChaincodeEvent, error) {
	switch name {
	case consts.TicketCCIssueFunction:
		return cc.handleTicketCCIssueTx(args...)
	case consts.TicketCCGetTicketFunction:
		return cc.handleTicketCCGetTicketTx(args...)
	case consts.TicketCCGetSectionTicketsFunction:
		return cc.handleTicketCCGetSectionTicketsTx(args...)

	default:
		return nil, uuid.Nil, nil, fmt.Errorf("function not found")
	}
}

func (cc DevChaincodeAPI) handleTicketCCIssueTx(args ...string) ([]byte, uuid.UUID, *client.ChaincodeEvent, error) {
	if len(args) != 5 {
		return nil, uuid.Nil, nil, fmt.Errorf("wrong arg numbers: expected %d, obtained %d", 5, len(args))
	}

	ticketID, _ := uuid.Parse(args[0])
	eventID, _ := uuid.Parse(args[1])
	section := args[2]
	ownerID, _ := uuid.Parse(args[3])

	// Create TokenID as big.Int from args[4]
	tokenID, _ := big.NewInt(0).SetString(args[4], 10)

	ticket := &chainmodels.Ticket{
		TicketID: ticketID,
		EventID:  eventID,
		OwnerID:  ownerID,
		Section:  section,
		Status:   "issued",
		TokenID:  tokenID,
	}

	ticketBytes, err := json.Marshal(ticket)
	if err != nil {
		return nil, uuid.Nil, nil, err
	}

	return ticketBytes, ticketID, nil, nil
}

func (cc DevChaincodeAPI) handleTicketCCGetTicketTx(args ...string) ([]byte, uuid.UUID, *client.ChaincodeEvent, error) {
	if len(args) != 1 {
		return nil, uuid.Nil, nil, fmt.Errorf("wrong arg numbers: expected %d, obtained %d", 1, len(args))
	}

	ticketID, _ := uuid.Parse(args[0])

	ticketBytes, exists := cc.storedElements[ticketElementName][ticketID]
	if !exists {
		return nil, uuid.Nil, nil, fmt.Errorf("ticket with id %s doest not exist", ticketID)
	}

	return ticketBytes, ticketID, nil, nil
}

func (cc DevChaincodeAPI) handleTicketCCGetSectionTicketsTx(args ...string) ([]byte, uuid.UUID, *client.ChaincodeEvent, error) {
	if len(args) != 2 {
		return nil, uuid.Nil, nil, fmt.Errorf("wrong arg numbers: expected %d, obtained %d", 2, len(args))
	}

	eventID, _ := uuid.Parse(args[0])
	section := args[1]

	var sectionTickets [][]byte
	var ticket chainmodels.Ticket
	for _, ticketBytes := range cc.storedElements[ticketElementName] {
		if err := json.Unmarshal(ticketBytes, &ticket); err != nil {
			return nil, uuid.Nil, nil, err
		}
		if ticket.EventID == eventID && ticket.Section == section {
			sectionTickets = append(sectionTickets, ticketBytes)
		}
	}

	sectionTicketsSerialized, _ := json.Marshal(sectionTickets)

	return sectionTicketsSerialized, uuid.Nil, nil, nil
}

func (cc DevChaincodeAPI) handleEventCCGetEventTx(args ...string) ([]byte, uuid.UUID, *client.ChaincodeEvent, error) {
	if len(args) != 1 {
		return nil, uuid.Nil, nil, fmt.Errorf("wrong arg numbers: expected %d, obtained %d", 1, len(args))
	}

	eventID, _ := uuid.Parse(args[0])

	eventBytes, exists := cc.storedElements[eventElementName][eventID]
	if !exists {
		return nil, uuid.Nil, nil, fmt.Errorf("event with id %s doest not exist", eventID)
	}

	return eventBytes, eventID, nil, nil
}

func (cc DevChaincodeAPI) handleEventCCCreateTx(args ...string) ([]byte, uuid.UUID, *client.ChaincodeEvent, error) {
	if len(args) != 3 {
		return nil, uuid.Nil, nil, fmt.Errorf("wrong arg numbers: expected %d, obtained %d", 3, len(args))
	}

	eventID, _ := uuid.Parse(args[0])
	name := args[1]
	date, _ := time.Parse(time.RFC3339, args[2])

	ccEvent := chainmodels.Event{
		EventID:  eventID,
		Name:     name,
		Date:     date,
		Sections: make([]*chainmodels.Section, 0),

		OnSale: false,

		MSPID:             cc.ctxMSPID,
		OrganizerUsername: cc.ctxOrganizerUsername,
	}

	eventBytes, _ := json.Marshal(ccEvent)

	notification := notificationFrom(eventBytes, consts.EventCreatedNotification, cc.name)

	return eventBytes, eventID, notification, nil
}

func (cc DevChaincodeAPI) handleEventCCSetEventOnSaleTx(args ...string) ([]byte, uuid.UUID, *client.ChaincodeEvent, error) {
	if len(args) != 1 {
		return nil, uuid.Nil, nil, fmt.Errorf("wrong arg numbers: expected %d, obtained %d", 1, len(args))
	}

	eventID, _ := uuid.Parse(args[0])

	eventBytes, exists := cc.storedElements[eventElementName][eventID]
	if !exists {
		return nil, uuid.Nil, nil, fmt.Errorf("event with id %s not exists", eventID)
	}

	var event chainmodels.Event
	err := json.Unmarshal(eventBytes, &event)
	if err != nil {
		return nil, uuid.Nil, nil, err
	}

	event.OnSale = true

	eventModifidBytes, err := json.Marshal(event)
	if err != nil {
		return nil, uuid.Nil, nil, err
	}

	return eventModifidBytes, eventID, nil, nil
}

func (cc DevChaincodeAPI) handleEventCCAddSectionTx(args ...string) ([]byte, uuid.UUID, *client.ChaincodeEvent, error) {
	if len(args) != 4 {
		return nil, uuid.Nil, nil, fmt.Errorf("wrong arg numbers: expected %d, obtained %d", 4, len(args))
	}

	eventID, _ := uuid.Parse(args[0])
	name := args[1]
	totalTickets, _ := strconv.Atoi(args[2])
	ticketPrice, _ := strconv.ParseFloat(args[3], 64)

	eventBytes, exists := cc.storedElements[eventElementName][eventID]

	if !exists {
		return nil, uuid.Nil, nil, fmt.Errorf("event with id %s not exists", eventID)
	}

	var event chainmodels.Event
	err := json.Unmarshal(eventBytes, &event)
	if err != nil {
		return nil, uuid.Nil, nil, err
	}

	section := chainmodels.Section{
		EventID:      eventID,
		Name:         name,
		TicketPrice:  ticketPrice,
		TotalTickets: totalTickets,
		SoldTickets:  0,
	}

	event.Sections = append(event.Sections, &section)

	eventBytes, _ = json.Marshal(event)
	sectionBytes, _ := json.Marshal(section)

	notification := notificationFrom(sectionBytes, consts.SectionAddedNotification, cc.name)

	return eventBytes, eventID, notification, nil
}

// notificationFrom generates a fake notification to send in
// the fake notification channel. This mocked notification is
// equal to the notification received using the core connector,
// the only difference are the txID and blockNum values
func notificationFrom(payload []byte, notificationName string, chaincodeName string) *client.ChaincodeEvent {
	return &client.ChaincodeEvent{
		BlockNumber:   1,
		TransactionID: uuid.New().String(),
		ChaincodeName: chaincodeName,
		EventName:     notificationName,
		Payload:       payload,
	}
}
