package artemis_ethereum_validator_service

import (
	"net/http"

	"github.com/labstack/echo/v4"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
)

type EthereumValidatorServiceGroupCloudCtxNsRequest struct {
	hestia_autogen_bases.ValidatorsServiceOrgGroupsCloudCtxNsSlice
}

func EthereumEphemeryValidatorHandler(c echo.Context) error {
	request := new(EthereumValidatorServiceGroupCloudCtxNsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.OrchestrateEphemeryValidatorPlacement(c)
}

func (v *EthereumValidatorServiceGroupCloudCtxNsRequest) OrchestrateEphemeryValidatorPlacement(c echo.Context) error {
	//ctx := context.Background()
	//ou := c.Get("orgUser").(org_users.OrgUser)

	//if err != nil {
	//	log.Err(err).Interface("orgUser", ou).Msg("SendEphemeralSignedTx, ExecuteArtemisSendSignedTxWorkflow error")
	//	return c.JSON(http.StatusBadRequest, nil)
	//}
	return c.JSON(http.StatusAccepted, v.ValidatorsServiceOrgGroupsCloudCtxNsSlice)
}
