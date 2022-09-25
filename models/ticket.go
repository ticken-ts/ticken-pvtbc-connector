package models

type Ticket struct {
	TicketID string `json:"ticket_id"`
	EventID  string `json:"event_id"`
	Owner    string `json:"owner"`
	Status   string `json:"status"`
}
