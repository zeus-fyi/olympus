package config_maps

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

type ConfigMapActionRequest struct {
	kns.TopologyKubeCtxNs
	Action        string
	ConfigMapName string
	Keys          KeySwap
	FilterOpts    *string_utils.FilterOpts
}

type KeySwap struct {
	KeyOne string
	KeyTwo string
}

func (cfm *ConfigMapActionRequest) KeySwap(c echo.Context) error {
	ctx := context.Background()
	cm, err := zeus.K8Util.ConfigMapKeySwap(ctx, cfm.CloudCtxNs, cfm.FilterOpts, cfm.ConfigMapName, cfm.Keys.KeyOne, cfm.Keys.KeyTwo)
	if err != nil {
		log.Err(err).Msg("ConfigmapKeySwap")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, cm)
}
