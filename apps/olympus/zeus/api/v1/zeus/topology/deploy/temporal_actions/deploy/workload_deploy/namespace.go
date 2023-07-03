package internal_deploy

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/base_request"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

func DeployNamespaceHandler(c echo.Context) error {
	ctx := context.Background()
	request := new(base_request.InternalDeploymentActionRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	log.Debug().Interface("kns", request.Kns).Msg("DeployNamespaceHandler: CreateNamespaceIfDoesNotExist")
	_, err := zeus.K8Util.CreateNamespaceIfDoesNotExist(ctx, request.Kns.CloudCtxNs)
	if err != nil {
		log.Err(err).Msg("DeployNamespaceHandler")
		return c.JSON(http.StatusInternalServerError, err)
	}
	if request.Kns.CloudCtxNs.Context == "zeusfyi" && request.Kns.CloudCtxNs.CloudProvider == "ovh" {
		fromKns := zeus_common_types.CloudCtxNs{
			CloudProvider: "do",
			Region:        "sfo3",
			Context:       "do-sfo3-dev-do-sfo3-zeus",
			Namespace:     "zeus",
			Alias:         "zeus",
			Env:           "",
		}
		_, err = zeus.K8Util.CopySecretToAnotherKns(ctx, fromKns, request.Kns.CloudCtxNs, "zeus-fyi-ext", nil)
		if err != nil {
			log.Err(err).Msg("DeploySecretsHandler")
			return c.JSON(http.StatusInternalServerError, err)
		}
		namespace := request.Kns.CloudCtxNs.Namespace
		switch namespace {
		case "artemis", "hardhat", "zeus", "iris", "hestia", "hera", "aegis", "poseidon", "ephemeral-staking", "goerli-staking", "olympus", "mainnet-staking", "tyche":
			_, err = zeus.K8Util.CopySecretToAnotherKns(ctx, fromKns, request.Kns.CloudCtxNs, "age-auth", nil)
			if err != nil {
				log.Err(err).Msg("DeploySecretsHandler")
				return c.JSON(http.StatusInternalServerError, err)
			}
			_, err = zeus.K8Util.CopySecretToAnotherKns(ctx, fromKns, request.Kns.CloudCtxNs, "spaces-auth", nil)
			if err != nil {
				log.Err(err).Msg("DeploySecretsHandler")
				return c.JSON(http.StatusInternalServerError, err)
			}
			_, err = zeus.K8Util.CopySecretToAnotherKns(ctx, fromKns, request.Kns.CloudCtxNs, "spaces-key", nil)
			if err != nil {
				log.Err(err).Msg("DeploySecretsHandler")
				return c.JSON(http.StatusInternalServerError, err)
			}
		default:
			_, err = zeus.K8Util.CopySecretToAnotherKns(ctx, fromKns, request.Kns.CloudCtxNs, "age-auth", nil)
			if err != nil {
				log.Err(err).Msg("DeploySecretsHandler")
				return c.JSON(http.StatusInternalServerError, err)
			}
			_, err = zeus.K8Util.CopySecretToAnotherKns(ctx, fromKns, request.Kns.CloudCtxNs, "spaces-auth", nil)
			if err != nil {
				log.Err(err).Msg("DeploySecretsHandler")
				return c.JSON(http.StatusInternalServerError, err)
			}
			_, err = zeus.K8Util.CopySecretToAnotherKns(ctx, fromKns, request.Kns.CloudCtxNs, "spaces-key", nil)
			if err != nil {
				log.Err(err).Msg("DeploySecretsHandler")
				return c.JSON(http.StatusInternalServerError, err)
			}
		}
	}
	return c.JSON(http.StatusOK, nil)
}
