package chain_models

import (
	"time"
)

type Event struct {
	EventID  string     `json:"event_id"`
	Name     string     `json:"name"`
	Date     time.Time  `json:"date"`
	Sections []*Section `json:"sections"`

	// identity of the event and auditory
	MSPID             string `json:"msp_id"`
	OrganizerUsername string `json:"organizer_username"`
}

type Section struct {
	Name         string  `json:"name"`
	TicketPrice  float64 `json:"ticket_price"`
	TotalTickets int     `json:"total_tickets"`
	SoldTickets  int     `json:"sold_tickets"`
}
