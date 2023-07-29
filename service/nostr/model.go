package pstr

import (
	"context"
	"fmt"
	"log"

	"github.com/lnconsole/photobolt/env"
	"github.com/nbd-wtf/go-nostr"
)

var (
	Mainch          = make(chan nostr.Event)
	KindJobFeedback = 65000
	KindJobResult   = 65001
	RelayUrl        string
	relay           *nostr.Relay
	pk              string
)

func Init(
	relayUrl string,
	sk string,
) error {
	var err error
	RelayUrl = relayUrl

	pk, err = nostr.GetPublicKey(sk)
	log.Printf("pk: %s", pk)
	if err != nil {
		return err
	}

	// connect to relay
	relay, err = nostr.RelayConnect(context.Background(), RelayUrl)
	if err != nil {
		return err
	}
	log.Printf("connected to: %s", relay)
	go func() {
		for notice := range relay.Notices {
			log.Printf("(%s) notice: %s", relay.URL, notice)
		}
	}()
	return nil
}

func Subscribe(filters nostr.Filters) *nostr.Subscription {
	sub := relay.Subscribe(context.Background(), filters)
	go func() {
		for evt := range sub.Events {
			Mainch <- *evt
		}
	}()
	return sub
}

func Publish(ctx context.Context, evt nostr.Event) (*nostr.Event, error) {
	if evt.PubKey == "" {
		evt.PubKey = pk
	}

	if evt.Sig == "" {
		err := evt.Sign(env.PhotoBolt.NostrPrivateKey)
		if err != nil {
			return nil, fmt.Errorf("error signing event: %w", err)
		}
	}
	// log.Printf("%v", evt)
	sRelay, err := nostr.RelayConnect(ctx, RelayUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %s", RelayUrl, err)
	}
	log.Printf("posting to: %s,kind: %d,id: %s, %s", RelayUrl, evt.Kind, evt.ID, sRelay.Publish(ctx, evt))
	sRelay.Close()

	return &evt, nil
}
