package v1_ethereum_aws

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	serverless_aws_automation "github.com/zeus-fyi/zeus/builds/serverless/aws_automation"
	aws_aegis_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

func CreateLambdaFunctionHandler(c echo.Context) error {
	request := new(AwsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateKeystores(c)
}

func (a *AwsRequest) CreateLambdaFunction(c echo.Context) error {
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	lambdaFnUrl, err := serverless_aws_automation.CreateLambdaFunction(ctx, a.AuthAWS)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("ou", ou).Msg("AwsRequest, CreateLambdaFunction error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, lambdaFnUrl)
}

func VerifyLambdaFunctionHandler(c echo.Context) error {
	request := new(AwsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.VerifyLambdaFunction(c)
}

func (a *AwsRequest) VerifyLambdaFunction(c echo.Context) error {
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	lambdaAccessAuth := aws_aegis_auth.AuthAWS{
		Region:    "us-west-1",
		AccessKey: a.AccessKey,
		SecretKey: a.SecretKey,
	}
	// TODO
	ageEncryptionSecretName := ""
	serviceURL := ""
	err := serverless_aws_automation.VerifyLambdaSigner(ctx, lambdaAccessAuth, filepaths.Path{}, serviceURL, ageEncryptionSecretName)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("ou", ou).Msg("VerifyRequest VerifyLambdaFunction error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}
