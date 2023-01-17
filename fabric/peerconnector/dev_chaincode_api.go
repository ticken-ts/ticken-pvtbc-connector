package peerconnector

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	chainmodels "github.com/ticken-ts/ticken-pvtbc-connector/chain-models"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/consts"
	"strconv"
	"time"
)

type DevChaincodeAPI struct {
	name    string
	channel string

	// these are the context identity
	// which is requesting the chaincode
	// to submit/evaluate a transaction
	ctxMSPID             string
	ctxOrganizerUsername string

	storedElements           map[uuid.UUID][]byte
	fakeNotificationsChannel chan *client.ChaincodeEvent
}

func (cc DevChaincodeAPI) ChaincodeName() string {
	return cc.name
}

func (cc DevChaincodeAPI) SubmitTx(name string, args ...string) ([]byte, error) {
	var payload []byte
	var err error
	var elemKey uuid.UUID
	var notification string

	switch cc.ChaincodeName() {
	case consts.TickenTicketChaincode:
		return nil, nil
	case consts.TickenEventChaincode:
		payload, elemKey, notification, err = cc.handleEventCCAPI(name, args...)
	default:
		return nil, fmt.Errorf("chaincode %s not exists", cc.ChaincodeName())
	}

	if err != nil {
		return nil, err
	}

	cc.storedElements[elemKey] = payload
	cc.fakeNotificationsChannel <- notificationFrom(payload, notification, cc.name)
	return payload, nil
}

func (cc DevChaincodeAPI) EvaluateTx(name string, args ...string) ([]byte, error) {
	var payload []byte
	var err error

	switch cc.ChaincodeName() {
	case consts.TickenTicketChaincode:
		return nil, nil
	case consts.TickenEventChaincode:
		payload, _, _, err = cc.handleEventCCAPI(name, args...)
	default:
		return nil, fmt.Errorf("chaincode %s not exists", cc.ChaincodeName())
	}

	return payload, err
}

func (cc DevChaincodeAPI) SubmitTxAsync(name string, args ...string) ([]byte, error) {
	return cc.SubmitTx(name, args...)
}

// handleEventCCAPI decides based on the function name, how to mock the
// operation, parsing the args correct and generating a similar output
// based on what the real chaincode should respond
func (cc DevChaincodeAPI) handleEventCCAPI(name string, args ...string) ([]byte, uuid.UUID, string, error) {
	switch name {
	case consts.EventCCCreateFunction:
		return cc.handleEventCCCreateTx(args...)
	case consts.EventCCAddSectionFunction:
		return cc.handleEventCCAddSectionTx(args...)
	case consts.EventCCGetFunction:
		return cc.handleEventCCAddSectionTx(args...)
	default:
		return nil, uuid.Nil, "", fmt.Errorf("function not found")
	}
}

func (cc DevChaincodeAPI) handleEventCCGetTx(args ...string) ([]byte, uuid.UUID, string, error) {
	if len(args) != 1 {
		return nil, uuid.Nil, "", fmt.Errorf("wrong arg numbers: expected %d, obtained %d", 1, len(args))
	}

	eventID, _ := uuid.Parse(args[0])

	eventBytes, exists := cc.storedElements[eventID]
	if !exists {
		return nil, uuid.Nil, "", fmt.Errorf("event with id %s doest not exist", eventID)
	}

	return eventBytes, eventID, "", nil
}

func (cc DevChaincodeAPI) handleEventCCCreateTx(args ...string) ([]byte, uuid.UUID, string, error) {
	if len(args) != 3 {
		return nil, uuid.Nil, "", fmt.Errorf("wrong arg numbers: expected %d, obtained %d", 3, len(args))
	}

	eventID, _ := uuid.Parse(args[0])
	name := args[1]
	date, _ := time.Parse(time.RFC3339, args[2])

	ccEvent := chainmodels.Event{
		EventID:  eventID,
		Name:     name,
		Date:     date,
		Sections: make([]*chainmodels.Section, 0),

		MSPID:             cc.ctxMSPID,
		OrganizerUsername: cc.ctxOrganizerUsername,
	}

	eventBytes, _ := json.Marshal(ccEvent)

	return eventBytes, eventID, consts.EventCreatedNotification, nil
}

func (cc DevChaincodeAPI) handleEventCCAddSectionTx(args ...string) ([]byte, uuid.UUID, string, error) {
	if len(args) != 4 {
		return nil, uuid.Nil, "", fmt.Errorf("wrong arg numbers: expected %d, obtained %d", 4, len(args))
	}

	eventID, _ := uuid.Parse(args[0])
	name := args[1]
	totalTickets, _ := strconv.Atoi(args[2])
	ticketPrice, _ := strconv.ParseFloat(args[3], 64)

	eventBytes, exists := cc.storedElements[eventID]

	if !exists {
		return nil, uuid.Nil, "", fmt.Errorf("event with id %s not exists", eventID)
	}

	var event chainmodels.Event
	err := json.Unmarshal(eventBytes, &event)
	if err != nil {
		return nil, uuid.Nil, "", err
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

	return eventBytes, eventID, consts.SectionAddedNotification, nil
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
