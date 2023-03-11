package v1_ethereum_aws

import (
	"net/http"

	"github.com/labstack/echo/v4"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
)

type EthereumRequest struct {
	aegis_aws_auth.AuthAWS `json:"authAWS"`
}

func GenerateValidatorsHandler(c echo.Context) error {
	request := new(EthereumRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GenerateValidators(c)
}

func (e *EthereumRequest) GenerateValidators(c echo.Context) error {
	//ctx := context.Background()
	//ou := c.Get("orgUser").(org_users.OrgUser)
	//err := serverless_aws_automation.CreateLambdaFunctionKeystoresLayer(ctx, e.AuthAWS, filepaths.Path{})
	//if err != nil {
	//	log.Ctx(ctx).Err(err).Interface("ou", ou).Msg("AwsRequest, CreateInternalServerlessUser error")
	//	return c.JSON(http.StatusInternalServerError, err)
	//}
	return c.JSON(http.StatusOK, nil)
}
