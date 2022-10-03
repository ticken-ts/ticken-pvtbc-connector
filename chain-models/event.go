package chain_models

import (
	"time"
)

type Section struct {
	Name             string `json:"name"`
	TotalTickets     int    `json:"total_tickets"`
	RemainingTickets int    `json:"remaining_tickets"`
}

type Event struct {
	EventID        string    `json:"event_id"`
	Name           string    `json:"name"`
	Date           time.Time `json:"date"`
	Sections       []Section `json:"sections"`
	OrganizationID string    `json:"organization_id"`
}
