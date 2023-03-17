package artemis_ethereum_validator_service

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	signing_automation_ethereum "github.com/zeus-fyi/zeus/pkg/artemis/signing_automation/ethereum"
)

type DepositEthereumValidatorsService struct {
	Network               string                                        `json:"network"`
	ValidatorDepositSlice []signing_automation_ethereum.DepositDataJSON `json:"validatorDepositSlice"`
}

func CreateEthereumValidatorsHandler(c echo.Context) error {
	request := new(DepositEthereumValidatorsService)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.DepositValidators(c)
}

func (v *DepositEthereumValidatorsService) DepositValidators(c echo.Context) error {
	//ctx := context.Background()
	//ou := c.Get("orgUser").(org_users.OrgUser)

	//if err != nil {
	//	log.Err(err).Interface("orgUser", ou).Msg("SendEphemeralSignedTx, ExecuteArtemisSendSignedTxWorkflow error")
	//	return c.JSON(http.StatusBadRequest, nil)
	//}
	// TODO

	switch strings.ToLower(v.Network) {
	case "mainnet":
		return c.JSON(http.StatusNotImplemented, nil)
	case "goerli":
		return c.JSON(http.StatusNotImplemented, nil)
	case "ephemery":
		return c.JSON(http.StatusAccepted, nil)
	default:
		return c.JSON(http.StatusBadRequest, nil)
	}
}
