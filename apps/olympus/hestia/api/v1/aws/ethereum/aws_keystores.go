package v1_aws_ethereum_automation

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	serverless_aws_automation "github.com/zeus-fyi/zeus/builds/serverless/aws_automation"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

func CreateServerlessKeystoresHandler(c echo.Context) error {
	request := new(AwsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateKeystores(c)
}

func (a *AwsRequest) CreateKeystores(c echo.Context) error {
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	err := serverless_aws_automation.CreateLambdaFunctionKeystoresLayer(ctx, a.AuthAWS, filepaths.Path{})
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("ou", ou).Msg("AwsRequest, CreateKeystores error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}
