package chain_models

import (
	"github.com/google/uuid"
	"time"
)

type EventStatus string

const (
	// EventStatusDraft is the status of an
	// event that is not yet published
	EventStatusDraft EventStatus = "draft"

	// EventStatusOnSale is the status of an
	// event that is published for sale
	EventStatusOnSale EventStatus = "on_sale"

	// EventStatusRunning is the status of an
	// event that is currently happening
	EventStatusRunning EventStatus = "running"

	// EventStatusFinished is the status of an
	// event that has finished
	EventStatusFinished EventStatus = "finished"
)

type Event struct {
	EventID  uuid.UUID   `json:"event_id"`
	Name     string      `json:"name"`
	Date     time.Time   `json:"date"`
	Sections []*Section  `json:"sections"`
	Status   EventStatus `json:"status"`

	// identity of the event and auditory
	MSPID             string `json:"msp_id"`
	OrganizerUsername string `json:"organizer_username"`
}

type Section struct {
	EventID      uuid.UUID `json:"event_id"`
	Name         string    `json:"name"`
	TicketPrice  float64   `json:"ticket_price"`
	TotalTickets int       `json:"total_tickets"`
	SoldTickets  int       `json:"sold_tickets"`
}
