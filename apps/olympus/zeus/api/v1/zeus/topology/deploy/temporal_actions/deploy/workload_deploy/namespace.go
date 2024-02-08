package internal_deploy

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/dynamic_secrets"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/base_request"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

func DeployNamespaceHandlerWrapper(k autok8s_core.K8Util) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := context.Background()
		request := new(base_request.InternalDeploymentActionRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		log.Debug().Interface("kns", request.Kns).Msg("DeployNamespaceHandler: CreateNamespaceIfDoesNotExist")
		_, err := k.CreateNamespaceIfDoesNotExist(ctx, request.Kns.CloudCtxNs)
		if err != nil {
			log.Err(err).Msg("DeployNamespaceHandler")
			return c.JSON(http.StatusInternalServerError, err)
		}

		if strings.HasPrefix(request.Kns.CloudCtxNs.Namespace, "sui") {
			sec := dynamic_secrets.GetS3SecretSui(ctx, request.Kns.CloudCtxNs)
			_, serr := k.CreateSecretWithKnsIfDoesNotExist(ctx, request.Kns.CloudCtxNs, &sec, nil)
			if serr != nil {
				log.Err(serr).Msg("DeployNamespaceHandler")
				return c.JSON(http.StatusInternalServerError, serr)
			}
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
			_, err = k.CopySecretToAnotherKns(ctx, fromKns, request.Kns.CloudCtxNs, "zeus-fyi-ext", nil)
			if err != nil {
				log.Err(err).Msg("DeploySecretsHandler")
				return c.JSON(http.StatusInternalServerError, err)
			}
			namespace := request.Kns.CloudCtxNs.Namespace
			switch namespace {
			case "artemis", "hardhat", "zeus", "iris", "hestia", "hera", "aegis", "poseidon", "ephemeral-staking", "goerli-staking", "olympus", "mainnet-staking", "tyche", "keydb", "pandora":
				_, err = k.CopySecretToAnotherKns(ctx, fromKns, request.Kns.CloudCtxNs, "age-auth", nil)
				if err != nil {
					log.Err(err).Msg("DeploySecretsHandler")
					return c.JSON(http.StatusInternalServerError, err)
				}
				_, err = k.CopySecretToAnotherKns(ctx, fromKns, request.Kns.CloudCtxNs, "spaces-auth", nil)
				if err != nil {
					log.Err(err).Msg("DeploySecretsHandler")
					return c.JSON(http.StatusInternalServerError, err)
				}
				_, err = k.CopySecretToAnotherKns(ctx, fromKns, request.Kns.CloudCtxNs, "spaces-key", nil)
				if err != nil {
					log.Err(err).Msg("DeploySecretsHandler")
					return c.JSON(http.StatusInternalServerError, err)
				}
			default:
				_, err = k.CopySecretToAnotherKns(ctx, fromKns, request.Kns.CloudCtxNs, "age-auth", nil)
				if err != nil {
					log.Err(err).Msg("DeploySecretsHandler")
					return c.JSON(http.StatusInternalServerError, err)
				}
				_, err = k.CopySecretToAnotherKns(ctx, fromKns, request.Kns.CloudCtxNs, "spaces-auth", nil)
				if err != nil {
					log.Err(err).Msg("DeploySecretsHandler")
					return c.JSON(http.StatusInternalServerError, err)
				}
				_, err = k.CopySecretToAnotherKns(ctx, fromKns, request.Kns.CloudCtxNs, "spaces-key", nil)
				if err != nil {
					log.Err(err).Msg("DeploySecretsHandler")
					return c.JSON(http.StatusInternalServerError, err)
				}
			}
		}
		return c.JSON(http.StatusOK, nil)
	}
}
