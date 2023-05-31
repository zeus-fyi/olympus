package internal_deploy

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/dynamic_secrets"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/base_request"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeployChoreographySecretsHandler(c echo.Context) error {
	request := new(base_request.InternalDeploymentActionRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
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
	_, err = zeus.K8Util.CreateSecretWithKnsIfDoesNotExist(ctx, request.Kns.CloudCtxNs, &sec, nil)
	if err != nil {
		log.Err(err).Msg("DeployChoreographySecretsHandler")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}

func DeployDynamicSecretsHandler(c echo.Context) error {
	request := new(base_request.InternalDeploymentActionRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	ctx := context.Background()
	sec, err := dynamic_secrets.LookupAndCreateSecret(ctx, request.OrgUser.OrgID, request.ClusterName, request.Kns.CloudCtxNs)
	if err != nil {
		log.Err(err).Msg("DeployDynamicSecretsHandler")
		return c.JSON(http.StatusInternalServerError, err)
	}
	_, err = zeus.K8Util.CreateSecretWithKnsIfDoesNotExist(ctx, request.Kns.CloudCtxNs, sec, nil)
	if err != nil {
		log.Err(err).Msg("DeployDynamicSecretsHandler")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}
