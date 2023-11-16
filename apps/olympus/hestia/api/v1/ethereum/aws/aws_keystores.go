package v1_ethereum_aws

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
	aws_lambda "github.com/zeus-fyi/zeus/pkg/cloud/aws/lambda"
)

type CreateAwsLambdaKeystoreLayerRequest struct {
	aegis_aws_auth.AuthAWS `json:"authAWS"`
	KeystoresLayerName     string `json:"keystoresLayerName,omitempty"`
}

func CreateServerlessKeystoresLayerHandler(c echo.Context) error {
	request := new(CreateAwsLambdaKeystoreLayerRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	request.KeystoresLayerName = c.FormValue("keystoresLayerName")
	authJSON := c.FormValue("authAWS")
	err := json.Unmarshal([]byte(authJSON), &request.AuthAWS)
	if err != nil {
		ou := c.Get("orgUser").(org_users.OrgUser)
		log.Err(err).Interface("ou", ou).Msg("CreateKeystoresLayer: error unmarshalling authAWS from form value")
		return c.JSON(http.StatusBadRequest, err)
	}
	return request.CreateKeystoresLayer(c)
}

func (a *CreateAwsLambdaKeystoreLayerRequest) CreateKeystoresLayer(c echo.Context) error {
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	zipFile, err := c.FormFile("keystoresZip")
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("CreateKeystoresLayer: error retrieving keystoresZip from request payload")
		return c.JSON(http.StatusBadRequest, err)
	}
	fi, err := zipFile.Open()
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("CreateKeystoresLayer: error retrieving keystoresZip from request payload")
		return c.JSON(http.StatusBadRequest, err)
	}
	defer fi.Close()
	zipBytes := new(bytes.Buffer)
	_, err = io.Copy(zipBytes, fi)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("CreateKeystoresLayer: error reading keystoresZip file")
		return c.JSON(http.StatusInternalServerError, err)
	}
	lm, err := aws_lambda.InitLambdaClient(ctx, a.AuthAWS)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("AwsRequest, CreateKeystoresLayer error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	ly, err := lm.CreateServerlessBLSLambdaFnKeystoreLayer(ctx, a.KeystoresLayerName, zipBytes.Bytes())
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("AwsRequest, CreateKeystoresLayer error")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, ly.Version)
}
