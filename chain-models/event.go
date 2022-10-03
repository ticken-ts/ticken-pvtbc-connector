package chain_models

import (
	"container/list"
	"time"
)

type Event struct {
	EventID        string    `json:"event_id"`
	Name           string    `json:"name"`
	Date           time.Time `json:"date"`
	Sections       list.List `json:"sections"`
	OrganizationID string    `json:"organization_id"`
}
