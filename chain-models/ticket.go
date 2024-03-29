package chain_models

import (
	"github.com/google/uuid"
)

type Ticket struct {
	TicketID uuid.UUID `json:"ticket_id"`
	Status   string    `json:"status"`

	EventID uuid.UUID `json:"event_id"`
	Section string    `json:"section"`

	TokenID      string `json:"token_id"`
	ContractAddr string `json:"contract_addr"`

	// represents the owner id in the
	// web service database
	OwnerID uuid.UUID `json:"owner"`
}
