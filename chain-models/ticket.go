package chain_models

type Ticket struct {
	TicketID string `json:"ticket_id"`
	EventID  string `json:"event_id"`
	Owner    string `json:"owner"`
	Section  string `json:"section"`
	Status   string `json:"status"`
}
