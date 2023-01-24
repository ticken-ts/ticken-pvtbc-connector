package chain_models

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	EventID  uuid.UUID  `json:"event_id"`
	Name     string     `json:"name"`
	Date     time.Time  `json:"date"`
	Sections []*Section `json:"sections"`

	// indicates if the event is currently
	// available to sell tickets
	OnSale bool `json:"on_sale"`

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
