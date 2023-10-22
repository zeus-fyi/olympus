package v1_iris

import (
	"context"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
)

// Step 1: Declare a channel outside the handler. This could be a global or part of a struct.
var dataChannel = make(chan []byte, 100)

// SendDataToWebSocket function that pushes data to the channel. Call this when you want to send data to the WebSocket.
func SendDataToWebSocket() {
	for {
		messages, err := iris_redis.IrisRedisClient.Stream(context.Background(), iris_redis.EthMempoolStreamName, "0")
		if err != nil {
			log.Err(err).Msg("error reading redis stream")
			return
		}
		for _, msg := range messages {
			for _, event := range msg.Messages {
				for _, v := range event.Values {
					dataChannel <- v.([]byte) // Assuming the Redis message can be directly converted to []byte
				}
			}
		}
	}
}

func mempoolWebSocketHandler(c echo.Context) error {
	conn, _, _, cerr := ws.UpgradeHTTP(c.Request(), c.Response().Writer)
	if cerr != nil {
		log.Err(cerr).Msg("mempoolWebSocketHandler: ws.UpgradeHTTP")
		return cerr
	}
	go func() {
		defer conn.Close()
		for {
			select {
			// Step 2: Modify the WebSocket handler to listen to the channel in the goroutine.
			case data := <-dataChannel:
				// Step 3: Write data to the WebSocket when data is received from the channel.
				err := wsutil.WriteServerMessage(conn, ws.OpBinary, data) // Assuming data is binary.
				if err != nil {
					log.Err(err).Msg("mempoolWebSocketHandler: wsutil.WriteServerMessage")
					// Handle error
					return
				}
			default:
				msg, op, err := wsutil.ReadClientData(conn)
				if err != nil {
					log.Err(err).Msg("mempoolWebSocketHandler: wsutil.ReadClientData")
					// Handle error, e.g., log it
					return
				}
				err = wsutil.WriteServerMessage(conn, op, msg)
				if err != nil {
					log.Err(err).Msg("mempoolWebSocketHandler: wsutil.WriteServerMessage")
					// Handle error
					return
				}
			}
		}
	}()
	return nil
}
