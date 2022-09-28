package cclisteners

import "github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"

type TickenTicketListener struct {
}

func NewTickenTicketListener(pc *peerconnector.PeerConnector, channel string) (*TickenTicketListener, error) {
	return new(TickenTicketListener), nil
}
