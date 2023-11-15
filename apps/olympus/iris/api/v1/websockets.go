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
			messages, err := iris_redis.IrisRedisClient.Stream(context.Background(), iris_redis.EthMempoolStreamName, "0")
			if err != nil {
				log.Err(err).Msg("error reading redis stream")
				return
			}

			for _, msg := range messages {
				for _, event := range msg.Messages {
					for _, v := range event.Values {
						payload, ok := v.(string)
						if !ok {
							continue
						}
						err = wsutil.WriteServerMessage(conn, ws.OpBinary, []byte(payload))
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
						}(orgID, []byte(payload))
					}
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
