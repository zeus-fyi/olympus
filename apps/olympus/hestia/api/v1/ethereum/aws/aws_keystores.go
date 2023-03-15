package v1_ethereum_aws

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	aws_lambda "github.com/zeus-fyi/zeus/pkg/cloud/aws/lambda"
)

func CreateServerlessKeystoresLayerHandler(c echo.Context) error {
	request := new(CreateAwsLambdaSignerRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateKeystoresLayer(c)
}

func (a *CreateAwsLambdaSignerRequest) CreateKeystoresLayer(c echo.Context) error {
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	lm, err := aws_lambda.InitLambdaClient(ctx, a.AuthAWS)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("ou", ou).Msg("AwsRequest, CreateKeystoresLayer error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	ly, err := lm.CreateServerlessBLSLambdaFnKeystoreLayer(ctx, a.FunctionName, nil)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("ou", ou).Msg("AwsRequest, CreateKeystoresLayer error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, ly.Version)
}
