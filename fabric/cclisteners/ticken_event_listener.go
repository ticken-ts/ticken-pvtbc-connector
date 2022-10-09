package cclisteners

import (
	"context"
	"encoding/json"
	"fmt"
	chain_models "github.com/ticken-ts/ticken-pvtbc-connector/chain-models"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/ccclient"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
)

const TickenEventChaincode = "ticken-event"

type TickenEventListener struct {
	listener  *ccclient.Listener
	callbacks map[string]func(event *chain_models.Event)
}

func NewTickenEventListener(pc *peerconnector.PeerConnector, channel string) (*TickenEventListener, error) {
	eventListener := new(TickenEventListener)
	listener, err := ccclient.NewListener(pc, channel, TickenEventChaincode)
	if err != nil {
		return nil, err
	}

	eventListener.listener = listener
	eventListener.callbacks = make(map[string]func(event *chain_models.Event))

	return eventListener, nil

}

func (eventListener *TickenEventListener) ListenNewEvents(ctx context.Context, callback func(event *chain_models.Event)) error {

	_, exists := eventListener.callbacks["create"]
	if exists {
		return fmt.Errorf("already listening to this event")
	}

	eventListener.callbacks["create"] = callback

	internalCallback := func(payload []byte) {
		event := new(chain_models.Event)
		err := json.Unmarshal(payload, event)
		if err != nil {
			panic(err)
		}

		callback(event)
	}

	eventListener.listener.Listen(ctx, "create", internalCallback)
	return nil
}

func (eventListener *TickenEventListener) ListenEventModifications(ctx context.Context, callback func(event *chain_models.Event)) error {

	_, exists := eventListener.callbacks["eventModified"]
	if exists {
		return fmt.Errorf("already listening to this event")
	}

	eventListener.callbacks["eventModified"] = callback

	internalCallback := func(payload []byte) {
		event := new(chain_models.Event)
		err := json.Unmarshal(payload, event)
		if err != nil {
			panic(err)
		}

		callback(event)
	}

	eventListener.listener.Listen(ctx, "eventModified", internalCallback)
	return nil
}
