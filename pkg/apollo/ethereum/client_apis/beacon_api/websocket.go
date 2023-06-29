package beacon_api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// Request represents the request structure to be sent to the WebSocket server.
type Request struct {
	ID      int      `json:"id"`
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
}

// Message represents the incoming messages structure from the WebSocket server after subscription.
type Message struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Result       map[string]interface{} `json:"result"`
		Subscription string                 `json:"subscription"`
	} `json:"params"`
}

func SubscribeToEvent(ctx context.Context, wsAddr string) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// Connect to WebSocket server.
	ws, _, err := websocket.Dial(ctx, wsAddr, &websocket.DialOptions{})
	if err != nil {
		log.Err(err).Msg("failed to connect to WebSocket server")
	}
	defer ws.Close(websocket.StatusInternalError, "failed to close conn to WebSocket server")
	// Create subscription request.
	request := Request{
		ID:      1,
		Jsonrpc: "2.0",
		Method:  "eth_subscribe",
		Params:  []string{"newHeads"},
	}
	// Send the subscription request to the WebSocket server.
	err = wsjson.Write(ctx, ws, request)
	if err != nil {
		log.Err(err).Msg("Failed to send subscription request to the WebSocket server")
	}
	for {
		// Read messages from the WebSocket server.
		var msg Message
		err = wsjson.Read(ctx, ws, &msg)
		if err != nil {
			log.Err(err).Msg("Failed to read message from the WebSocket server")
		}
		// Print the received message.
		result, rerr := json.MarshalIndent(msg, "", "  ")
		if rerr != nil {
			log.Err(err).Msg("Failed to parse the received message")
		}
		fmt.Println(string(result))
	}
}
