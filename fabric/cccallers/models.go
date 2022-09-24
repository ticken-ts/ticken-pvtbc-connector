package cccallers

import (
	"container/list"
	"time"
)

type Event struct {
	EventID  string    `json:"event_id"`
	Name     string    `json:"name"`
	Date     time.Time `json:"date"`
	Sections list.List `json:"sections"`
}

type Ticket struct {
	TicketID string `json:"ticket_id"`
	EventID  string `json:"event_id"`
	Owner    string `json:"owner"`
	Status   string `json:"status"`
}
