package v1_iris

import (
	"context"
	"net/http"
	"strings"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
)

// Step 1: Declare a channel outside the handler. This could be a global or part of a struct.
var dataChannel = make(chan string, 100)

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
					dataChannel <- v.(string) // Assuming the Redis message can be directly converted to []byte
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
	ou := org_users.OrgUser{}
	ouc := c.Get("orgUser")
	if ouc != nil {
		ouser, ok := ouc.(org_users.OrgUser)
		if ok && ouser.OrgID > 0 {
			ou = ouser
		} else {
			return c.JSON(http.StatusUnauthorized, Response{Message: "user not found"})
		}
	}
	plan := ""
	if c.Get("servicePlan") != nil {
		sp, ok := c.Get("servicePlan").(string)
		if ok {
			plan = MempoolPlan(sp)
		}
	}
	if plan == "" {
		return c.JSON(http.StatusBadRequest, Response{Message: "no service plan found, signup on QuickNode Marketplace to get started"})
	}

	go func(orgID int) {
		defer conn.Close()
		for {
			select {
			// Step 2: Modify the WebSocket handler to listen to the channel in the goroutine.
			case data := <-dataChannel:
				// Step 3: Write data to the WebSocket when data is received from the channel.
				err := wsutil.WriteServerMessage(conn, ws.OpBinary, []byte(data))
				if err != nil {
					log.Err(err).Msg("mempoolWebSocketHandler: wsutil.WriteServerMessage")
					// Handle error
					return
				}
				go func(orgID int, bodyBytes []byte) {
					ps := iris_usage_meters.NewPayloadSizeMeter(bodyBytes)
					err = iris_redis.IrisRedisClient.IncrementResponseUsageRateMeter(context.Background(), ou.OrgID, ps)
					if err != nil {
						log.Err(err).Interface("ou", ou).Msg("mempoolWebSocketHandler: IncrementResponseUsageRateMeter")
					}
				}(orgID, []byte(data))
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
	}(ou.OrgID)
	return nil
}

func MempoolPlan(plan string) string {
	switch strings.ToLower(plan) {
	case "enterprise", "standard", "performance", "lite", "discovery", "discover":
		return plan
	default:
		return ""
	}
}
