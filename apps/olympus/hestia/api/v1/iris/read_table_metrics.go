package hestia_iris_v1_routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
	hestia_billing "github.com/zeus-fyi/olympus/hestia/web/billing"
	resty_base "github.com/zeus-fyi/zeus/zeus/z_client/base"
)

type ReadMetricsRequest struct{}

func ReadTableMetricsRequestHandler(c echo.Context) error {
	request := new(ReadMetricsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.ReadTableStats(c)
}

func (r *ReadMetricsRequest) ReadTableStats(c echo.Context) error {
	tblName := c.Param("groupName")
	token, ok := c.Get("bearer").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	resp, err := GetTableDetails(c.Request().Context(), token, tblName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, resp)
}

func GetTableDetails(ctx context.Context, token, tblName string) (iris_redis.TableMetricsSummary, error) {
	planUsageDetails := iris_redis.TableMetricsSummary{}
	rc := resty_base.GetBaseRestyClient(hestia_billing.IrisApiUrl, token)
	endpoint := fmt.Sprintf("/v1/table/%s/metrics", tblName)
	resp, err := rc.R().SetResult(&planUsageDetails).Get(endpoint)
	if err != nil {
		log.Err(err).Msg("GetPlan: IrisPlatformSetupCacheUpdateRequest")
		return planUsageDetails, err
	}
	if resp.StatusCode() >= 400 {
		log.Err(err).Msg("GetPlan: IrisPlatformSetupCacheUpdateRequest")
		return planUsageDetails, err
	}
	return planUsageDetails, err
}
