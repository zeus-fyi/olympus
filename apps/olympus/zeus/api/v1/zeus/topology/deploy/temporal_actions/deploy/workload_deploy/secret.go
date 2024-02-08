package internal_deploy

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/dynamic_secrets"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/base_request"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeployChoreographySecretsHandlerWrapper(k autok8s_core.K8Util) func(c echo.Context) error {
	return func(c echo.Context) error {
		k8CfgInterface := c.Get("k8Cfg")
		if k8CfgInterface != nil {
			k8Cfg, ok := k8CfgInterface.(autok8s_core.K8Util) // Ensure the type assertion is correct
			if ok {
				k = k8Cfg
			}
		}
		// Attempt to retrieve the InternalDeploymentActionRequest from the context
		requestInterface := c.Get("internalDeploymentActionRequest")
		if requestInterface == nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal deployment action request not found in context"})
		}
		request, ok := requestInterface.(*base_request.InternalDeploymentActionRequest) // Ensure the type assertion is correct
		if !ok {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid request type"})
		}
		ctx := context.Background()
		m := make(map[string]string)

		key, err := auth.FetchUserAuthToken(ctx, request.OrgUser)
		if err != nil {
			log.Err(err).Msg("DeployChoreographySecretsHandler")
			return c.JSON(http.StatusInternalServerError, err)
		}
		m["bearer"] = key.PublicKey
		m["cloud-provider"] = request.Kns.CloudProvider
		m["ctx"] = request.Kns.Context
		m["ns"] = request.Kns.Namespace
		m["region"] = request.Kns.Region
		sec := v1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "choreography",
				Namespace: request.Kns.Namespace,
			},
			StringData: m,
			Type:       "Opaque",
		}
		_, err = k.CreateSecretWithKnsIfDoesNotExist(ctx, request.Kns.CloudCtxNs, &sec, nil)
		if err != nil {
			log.Err(err).Msg("DeployChoreographySecretsHandler")
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, nil)
	}
}

func DeployDynamicSecretsHandlerWrapper(k autok8s_core.K8Util) func(c echo.Context) error {
	return func(c echo.Context) error {
		request := new(base_request.InternalDeploymentActionRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		ctx := context.Background()
		secrets, err := dynamic_secrets.LookupAndCreateSecrets(ctx, request.OrgUser.OrgID, request.Kns.SecretRef, request.Kns.CloudCtxNs)
		if err != nil {
			log.Err(err).Msg("DeployDynamicSecretsHandler")
			return c.JSON(http.StatusInternalServerError, err)
		}
		for _, sec := range secrets {
			_, err = k.CreateSecretWithKnsIfDoesNotExist(ctx, request.Kns.CloudCtxNs, sec, nil)
			if err != nil {
				log.Err(err).Msg("DeployDynamicSecretsHandler")
				return c.JSON(http.StatusInternalServerError, err)
			}
		}
		return c.JSON(http.StatusOK, nil)
	}
}
