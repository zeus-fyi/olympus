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

type CreateAwsLambdaSignerRequest struct {
	aegis_aws_auth.AuthAWS `json:"authAWS"`
	FunctionName           string `json:"functionName"`
	KeystoresLayerName     string `json:"keystoresLayerName,omitempty"`
}

func CreateBlsLambdaFunctionHandler(c echo.Context) error {
	request := new(CreateAwsLambdaSignerRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateLambdaFunctionBlsSigner(c)
}

func (a *CreateAwsLambdaSignerRequest) CreateLambdaFunctionBlsSigner(c echo.Context) error {
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	lambdaFnUrl, err := serverless_aws_automation.CreateLambdaFunction(ctx, a.AuthAWS, a.FunctionName, a.KeystoresLayerName)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("ou", ou).Msg("AwsRequest, CreateLambdaFunctionBlsSigner error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, lambdaFnUrl)
}

func CreateLambdaFunctionSecretsKeyGenHandler(c echo.Context) error {
	request := new(AwsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateLambdaFunctionSecretsKeyGen(c)
}

func (a *AwsRequest) CreateLambdaFunctionSecretsKeyGen(c echo.Context) error {
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	lambdaFnUrl, err := serverless_aws_automation.CreateLambdaFunctionSecretsKeyGen(ctx, a.AuthAWS)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("ou", ou).Msg("AwsRequest, CreateLambdaFunctionSecretsKeyGen error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, lambdaFnUrl)
}

func CreateLambdaFunctionEncZipGenHandler(c echo.Context) error {
	request := new(AwsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateLambdaFunctionEncZipGen(c)
}

func (a *AwsRequest) CreateLambdaFunctionEncZipGen(c echo.Context) error {
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	lambdaFnUrl, err := serverless_aws_automation.CreateLambdaFunctionEncryptedKeystoresZip(ctx, a.AuthAWS)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("ou", ou).Msg("AwsRequest, CreateLambdaFunctionEncZipGen error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, lambdaFnUrl)
}

func CreateLambdaFunctionDepositsGenHandler(c echo.Context) error {
	request := new(AwsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateLambdaFunctionDepositsGen(c)
}

func (a *AwsRequest) CreateLambdaFunctionDepositsGen(c echo.Context) error {
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	lambdaFnUrl, err := serverless_aws_automation.CreateLambdaFunctionDepositGen(ctx, a.AuthAWS)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("ou", ou).Msg("AwsRequest, CreateLambdaFunctionDepositsGen error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, lambdaFnUrl)
}
