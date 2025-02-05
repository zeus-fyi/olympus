package v1_ethereum_aws

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	serverless_aws_automation "github.com/zeus-fyi/zeus/builds/serverless/aws_automation"
)

func CreateServerlessInternalUserHandler(c echo.Context) error {
	request := new(AwsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateInternalServerlessUser(c)
}

func (a *AwsRequest) CreateInternalServerlessUser(c echo.Context) error {
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	err := serverless_aws_automation.InternalUserRolePolicySetupForLambdaDeployment(ctx, a.AuthAWS)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("ou", ou).Msg("AwsRequest, CreateInternalServerlessUser error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}

func CreateServerlessExternalUserHandler(c echo.Context) error {
	request := new(AwsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateExternalServerlessUser(c)
}

func (a *AwsRequest) CreateExternalServerlessUser(c echo.Context) error {
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	err := serverless_aws_automation.ExternalUserRolePolicySetupForLambdaDeployment(ctx, a.AuthAWS)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("ou", ou).Msg("AwsRequest, CreateExternalServerlessUser error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}
