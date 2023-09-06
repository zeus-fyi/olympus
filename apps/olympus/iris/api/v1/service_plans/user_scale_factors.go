package iris_service_plans

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
)

const (
	UpdateTableLatencyScaleFactor = "/table/:groupName/scale/latency"
	UpdateTableErrorScaleFactor   = "/table/:groupName/scale/error"
	UpdateTableDecayScaleFactor   = "/table/:groupName/scale/decay"
)

type ScaleFactorsUpdateRequest struct {
	Score float64 `json:"score"`
}

func UpdateLatencyScaleFactorRequestHandler(c echo.Context) error {
	log.Info().Msg("Iris: UpdateLatencyScaleFactorRequestHandler")
	request := new(ScaleFactorsUpdateRequest)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("Iris: UpdateLatencyScaleFactorRequestHandler")
		return err
	}
	return request.UpdateLatencyScaleFactor(c)
}

func (p *ScaleFactorsUpdateRequest) UpdateLatencyScaleFactor(c echo.Context) error {
	ou := org_users.OrgUser{}
	ouc := c.Get("orgUser")
	if ouc != nil {
		ouser, aok := ouc.(org_users.OrgUser)
		if aok {
			ou = ouser
		}
	}
	tblName := c.Param("groupName")
	if len(tblName) == 0 {
		return c.JSON(http.StatusBadRequest, "table name is required")
	}
	err := iris_redis.IrisRedisClient.SetTableLatencyScaleFactor(context.Background(), ou.OrgID, tblName, p.Score)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("UpdateLatencyScaleFactor: SetTableLatencyScaleFactor error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, "ok")
}

func UpdateErrorScaleFactorRequestHandler(c echo.Context) error {
	log.Info().Msg("Iris: UpdateErrorScaleFactorRequestHandler")
	request := new(ScaleFactorsUpdateRequest)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("Iris: UpdateErrorScaleFactorRequestHandler")
		return err
	}
	return request.UpdateErrorScaleFactor(c)
}

func (p *ScaleFactorsUpdateRequest) UpdateErrorScaleFactor(c echo.Context) error {
	ou := org_users.OrgUser{}
	ouc := c.Get("orgUser")
	if ouc != nil {
		ouser, aok := ouc.(org_users.OrgUser)
		if aok {
			ou = ouser
		}
	}
	tblName := c.Param("groupName")
	if len(tblName) == 0 {
		return c.JSON(http.StatusBadRequest, "table name is required")
	}
	err := iris_redis.IrisRedisClient.SetTableErrorScaleFactor(context.Background(), ou.OrgID, tblName, p.Score)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("UpdateErrorScaleFactor: SetTableErrorScaleFactor error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, "ok")
}

func UpdateDecayScaleFactorRequestHandler(c echo.Context) error {
	log.Info().Msg("Iris: UpdateDecayScaleFactorRequestHandler")
	request := new(ScaleFactorsUpdateRequest)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("Iris: UpdateDecayScaleFactorRequestHandler")
		return err
	}
	return request.UpdateDecayScaleFactor(c)
}

func (p *ScaleFactorsUpdateRequest) UpdateDecayScaleFactor(c echo.Context) error {
	ou := org_users.OrgUser{}
	ouc := c.Get("orgUser")
	if ouc != nil {
		ouser, aok := ouc.(org_users.OrgUser)
		if aok {
			ou = ouser
		}
	}
	tblName := c.Param("groupName")
	if len(tblName) == 0 {
		return c.JSON(http.StatusBadRequest, "table name is required")
	}
	err := iris_redis.IrisRedisClient.SetTableDecayScaleFactor(context.Background(), ou.OrgID, tblName, p.Score)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("UpdateDecayScaleFactor: SetTableDecayScaleFactor error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, "ok")
}
