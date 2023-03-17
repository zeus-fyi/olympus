package v1_ethereum_aws

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	serverless_aws_automation "github.com/zeus-fyi/zeus/builds/serverless/aws_automation"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
)

type AwsExternalUserAccessRequest struct {
	aegis_aws_auth.AuthAWS   `json:"authAWS"`
	ExternalUserName         string `json:"externalUserName"`
	ExternalAccessSecretName string `json:"externalAccessSecretName"`
}

func CreateServerlessExternalUserAuthHandler(c echo.Context) error {
	request := new(AwsExternalUserAccessRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateOrFetchExternalServerlessUserAuth(c)
}

func (a *AwsExternalUserAccessRequest) CreateOrFetchExternalServerlessUserAuth(c echo.Context) error {
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	externalAccessKeys, err := serverless_aws_automation.GetExternalAccessKeySecretIfExists(ctx, a.AuthAWS, a.ExternalAccessSecretName)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("ou", ou).Msg("AwsRequest, CreateExternalServerlessUserAuth error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	if externalAccessKeys.AccessKey != "" && externalAccessKeys.SecretKey != "" {
		return c.JSON(http.StatusOK, externalAccessKeys)
	}
	externalAccessKeys, err = serverless_aws_automation.CreateExternalLambdaUserAccessKeys(ctx, a.AuthAWS)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("ou", ou).Msg("AwsRequest, CreateExternalServerlessUserAuth error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	err = serverless_aws_automation.AddExternalAccessKeysInAWSSecretManager(ctx, a.AuthAWS, a.ExternalAccessSecretName, externalAccessKeys)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("ou", ou).Msg("AwsRequest, CreateExternalServerlessUserAuth error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, externalAccessKeys)
}
