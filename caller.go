package ticken_pvtbc_connector

import (
	"fmt"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/cccallers"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
)

type Caller struct {
	*cccallers.TickenEventCaller
	*cccallers.TickenTicketCaller

	channel string
	pc      peerconnector.PeerConnector
}

func NewCaller(pc peerconnector.PeerConnector) (*Caller, error) {
	if pc == nil {
		return nil, fmt.Errorf("peer connection is nil")
	}

	if !pc.IsConnected() {
		return nil, fmt.Errorf("peer connection is not stablished")
	}

	return &Caller{pc: pc}, nil
}

func (caller *Caller) SetChannel(channel string) error {
	if caller.channel == channel {
		// optimization to avoid changing channel
		// most of the peers will share the same
		// channel, so this optimization is useful
		// to avoid checking outside
		return nil
	}

	tickenEventCaller, err := cccallers.NewTickenEventCaller(caller.pc, channel)
	if err != nil {
		return err
	}

	tickenTicketCaller, err := cccallers.NewTickenTicketCaller(caller.pc, channel)
	if err != nil {
		return err
	}

	caller.TickenEventCaller = tickenEventCaller
	caller.TickenTicketCaller = tickenTicketCaller

	// update channel to keep reference
	// if it needs to be changed
	caller.channel = channel

	return nil
}
