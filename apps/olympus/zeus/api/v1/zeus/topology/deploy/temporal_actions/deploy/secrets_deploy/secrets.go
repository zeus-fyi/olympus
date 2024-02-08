package internal_secrets_deploy

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
)

type InternalSecretsCopyFromTo struct {
	SecretNames []string              `json:"secretNames"`
	FromKns     kns.TopologyKubeCtxNs `json:"fromKns"`
	ToKns       kns.TopologyKubeCtxNs `json:"toKns"`
}

func DeploySecretsHandlerWrapper(k autok8s_core.K8Util) func(c echo.Context) error {
	return func(c echo.Context) error {
		request := new(InternalSecretsCopyFromTo)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		ctx := context.Background()
		for _, sec := range request.SecretNames {
			_, err := k.CopySecretToAnotherKns(ctx, request.FromKns.CloudCtxNs, request.ToKns.CloudCtxNs, sec, nil)
			if err != nil {
				log.Err(err).Msg("DeploySecretsHandler")
				return c.JSON(http.StatusInternalServerError, err)
			}
		}
		return c.JSON(http.StatusOK, nil)
	}
}
