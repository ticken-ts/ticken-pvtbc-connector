package ticken_pvtbc_connector

import (
	"context"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/cclisteners"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
)

type Listener struct {
	*cclisteners.TickenEventListener
	*cclisteners.TickenTicketListener

	pc      *peerconnector.PeerConnector
	channel string
}

func NewListener(mspID string, certPath string, privatekeyPath string, peerEndpoint string, gatewayPeer string, tlsCertPath string) (*Listener, error) {
	listener := new(Listener)

	pc := peerconnector.New(mspID, certPath, privatekeyPath)

	err := pc.Connect(peerEndpoint, gatewayPeer, tlsCertPath)
	if err != nil {
		return nil, err
	}

	listener.pc = pc

	return listener, nil
}

func (listener *Listener) SetChannel(ctx context.Context, channel string) error {
	if listener.channel == channel {
		// optimization to avoid changing channel
		// most of the peers will share the same
		// channel, so this optimization is useful
		// to avoid checking outside
		return nil
	}

	eventListener, err := cclisteners.NewTickenEventListener(ctx, listener.pc, channel)
	if err != nil {
		return err
	}

	ticketListener, err := cclisteners.NewTickenTicketListener(listener.pc, channel)
	if err != nil {
		return err
	}

	listener.TickenTicketListener = ticketListener
	listener.TickenEventListener = eventListener

	// update channel to keep reference
	// if it needs to be changed
	listener.channel = channel

	return nil
}
