package v1_ethereum_aws

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	serverless_aws_automation "github.com/zeus-fyi/zeus/builds/serverless/aws_automation"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
	signing_automation_ethereum "github.com/zeus-fyi/zeus/pkg/artemis/web3/signing_automation/ethereum"
)

type VerifyAwsLambdaSignerRequest struct {
	aegis_aws_auth.AuthAWS `json:"authAWS"`
	KeySlice               []string `json:"keySlice"`
	SecretName             string   `json:"secretName"`
	FunctionURL            string   `json:"functionURL"`
}

func VerifyLambdaFunctionHandler(c echo.Context) error {
	request := new(VerifyAwsLambdaSignerRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.VerifyLambdaFunction(c)
}

func (a *VerifyAwsLambdaSignerRequest) VerifyLambdaFunction(c echo.Context) error {
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	lambdaAccessAuth := aegis_aws_auth.AuthAWS{
		Region:    "us-west-1",
		AccessKey: a.AccessKey,
		SecretKey: a.SecretKey,
	}
	dpSlice := signing_automation_ethereum.ValidatorDepositSlice{}
	for _, key := range a.KeySlice {
		dpSlice = append(dpSlice, signing_automation_ethereum.ExtendedDepositParams{ValidatorDepositParams: signing_automation_ethereum.ValidatorDepositParams{Pubkey: key}})
	}
	err := serverless_aws_automation.VerifyLambdaSignerFromDepositDataSlice(ctx, lambdaAccessAuth, dpSlice, a.FunctionURL, a.SecretName)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("ou", ou).Msg("VerifyRequest VerifyLambdaFunction error")
		err = serverless_aws_automation.VerifyLambdaSignerFromDepositDataSlice(ctx, lambdaAccessAuth, dpSlice, a.FunctionURL, a.SecretName)
		if err != nil {
			log.Ctx(ctx).Err(err).Interface("ou", ou).Msg("VerifyRequest VerifyLambdaFunction error")
			return c.JSON(http.StatusInternalServerError, err)
		}
	}
	return c.JSON(http.StatusOK, a.KeySlice)
}
