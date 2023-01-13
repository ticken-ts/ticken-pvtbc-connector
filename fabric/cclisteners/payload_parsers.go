package cclisteners

import (
	"encoding/json"
	"time"
)

// *************** Event Payload's Parsing ************** //

type PVTBCSectionDTO struct {
	EventID      string  `json:"event_id"`
	Name         string  `json:"name"`
	TicketPrice  float64 `json:"ticket_price"`
	TotalTickets int     `json:"total_tickets"`
	SoldTickets  int     `json:"sold_tickets"`
}

func ParseEventCreatedNotification(payload []byte) (*PVTBCSectionDTO, error) {
	sectionDTO := new(PVTBCSectionDTO)
	err := json.Unmarshal(payload, sectionDTO)
	return sectionDTO, err
}

// ****************************************************** //
//
//
//
// ************* Section Payload's Parsing ************** //

type PVTBCEventDTO struct {
	EventID           string    `json:"event_id"`
	Name              string    `json:"name"`
	Date              time.Time `json:"date"`
	MSPID             string    `json:"msp_id"`
	OrganizerUsername string    `json:"organizer_username"`
}

func ParseSectionAddedNotification(payload []byte) (*PVTBCEventDTO, error) {
	eventDTO := new(PVTBCEventDTO)
	err := json.Unmarshal(payload, eventDTO)
	return eventDTO, err
}

// ****************************************************** //
