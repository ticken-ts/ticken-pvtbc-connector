package ticken_pvtbc_connector

import (
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/cccallers"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
)

type Caller struct {
	*cccallers.TickenEventCaller
	*cccallers.TickenTicketCaller

	channel string
	pc      *peerconnector.PeerConnector
}

func NewCaller(mspID string, certPath string, privatekeyPath string, peerEndpoint string, gatewayPeer string, tlsCertPath string) (*Caller, error) {
	tickenConnector := new(Caller)

	pc := peerconnector.New(mspID, certPath, privatekeyPath)

	err := pc.Connect(peerEndpoint, gatewayPeer, tlsCertPath)
	if err != nil {
		return nil, err
	}

	tickenConnector.pc = pc

	return tickenConnector, nil
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
	// if it need to be changed
	caller.channel = channel

	return nil
}
