package peerconnector

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	chainmodels "github.com/ticken-ts/ticken-pvtbc-connector/chain-models"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/consts"
	"testing"
	"time"
)

var (
	testMSPID             = "tickenMSP"
	testOrganizerUsername = "joey.tribbiani"
)

func TestSubmitCreateEventTransaction(t *testing.T) {
	devConnector := NewDevConnector(testMSPID, testOrganizerUsername)

	err := devConnector.Connect("", "", "")
	if err != nil {
		t.Errorf("error connecting: %s", err.Error())
	}

	cc, err := devConnector.GetChaincode("ticken-channel", consts.TickenEventChaincode)
	if err != nil {
		t.Errorf("error getting chaincode: %s", err.Error())
	}

	eventID := uuid.New()
	name := "test-event"
	date := time.Date(2023, 11, 10, 0, 0, 0, 0, time.UTC)

	txPayload, err := cc.SubmitTx(
		consts.EventCCCreateFunction,
		eventID.String(),
		name,
		date.Format(time.RFC3339),
	)
	if err != nil {
		t.Errorf("error submitting transaction: %s", err.Error())
	}

	var ccEvent chainmodels.Event
	if err = json.Unmarshal(txPayload, &ccEvent); err != nil {
		t.Errorf("error unmarshaling tx payload: %s", err)
	}

	assert.Equal(t, ccEvent.EventID, eventID, "event id are not equal")
	assert.Equal(t, ccEvent.Name, name, "name are not equal")
	assert.Equal(t, ccEvent.Date, date, "date are not equal")
	assert.Equal(t, ccEvent.MSPID, testMSPID, "MSP ID are not equal")
	assert.Equal(t, ccEvent.OrganizerUsername, testOrganizerUsername, "organizer are not equal")
}
