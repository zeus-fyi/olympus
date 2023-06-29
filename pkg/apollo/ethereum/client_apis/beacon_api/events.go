package beacon_api

import (
	"context"
	"fmt"
	"time"

	"github.com/r3labs/sse/v2"
)

// topics
// head, block, attestation, voluntary_exit, finalized_checkpoint, chain_reorg, contribution_and_proof.

const topicsPath = "eth/v1/events?topics="

func beaconSubscriptionSSE(ctx context.Context, nodeEndpointURL, topic string) error {
	url := nodeEndpointURL + "/" + topicsPath + topic
	client := sse.NewClient(url)

	err := client.SubscribeRaw(func(msg *sse.Event) {
		// Got some data!
		fmt.Println(string(msg.Data))
	})
	time.Sleep(5 * time.Second)

	return err
}
