package beacon_api

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

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
	Params  Params `json:"params"`
}

type Params struct {
	Result       Result `json:"result"`
	Subscription string `json:"subscription"`
}
type Result struct {
	Difficulty       string `json:"difficulty"`
	ExtraData        string `json:"extraData"`
	GasLimit         string `json:"gasLimit"`
	GasUsed          string `json:"gasUsed"`
	LogsBloom        string `json:"logsBloom"`
	Miner            string `json:"miner"`
	Nonce            string `json:"nonce"`
	Number           string `json:"number"`
	ParentHash       string `json:"parentHash"`
	ReceiptRoot      string `json:"receiptRoot"`
	Sha3Uncles       string `json:"sha3Uncles"`
	StateRoot        string `json:"stateRoot"`
	Timestamp        string `json:"timestamp"`
	TransactionsRoot string `json:"transactionsRoot"`
}

func TriggerWorkflowOnNewBlockHeaderEvent(ctx context.Context, wsAddr string, timestampChan chan<- time.Time) {
	for {
		wfCtx, cancel := context.WithCancel(ctx)
		err := establishWebsocketConnection(wfCtx, wsAddr, timestampChan)
		if err != nil {
			log.Err(err).Msg("Error in WebSocket connection or processing, retrying in 1 second...")
			cancel()
			time.Sleep(1 * time.Second)
			continue
		}
		cancel()
	}
}

func establishWebsocketConnection(ctx context.Context, wsAddr string, timestampChan chan<- time.Time) error {
	// Connect to WebSocket server.
	ws, _, err := websocket.Dial(ctx, wsAddr, &websocket.DialOptions{})
	if err != nil {
		return err
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
		return err
	}

	for {
		// Read and process messages from the WebSocket server.
		err = processMessages(ctx, ws, timestampChan)
		if err != nil {
			return err
		}
	}
}

func processMessages(ctx context.Context, ws *websocket.Conn, timestampChan chan<- time.Time) error {
	var msg Message
	err := wsjson.Read(ctx, ws, &msg)
	if err != nil {
		return err
	}
	// Print the received message.
	_, rerr := json.MarshalIndent(msg, "", "  ")
	if rerr != nil {
		return rerr
	}
	if msg.Params.Result.Timestamp == "" {
		return nil
	}
	t, terr := hexToTime(msg.Params.Result.Timestamp)
	if terr != nil {
		return terr
	}
	timestampChan <- t
	log.Info().Msg(fmt.Sprintf("New block header event received at %s", t))
	return nil
}

func hexToTime(hexStr string) (time.Time, error) {
	// strip the '0x' prefix
	cleanHex := hexStr[2:]

	// convert hex to int64 (base 16)
	sec, err := strconv.ParseInt(cleanHex, 16, 64)
	if err != nil {
		return time.Time{}, err
	}

	// convert seconds to time
	return time.Unix(sec, 0), nil
}
