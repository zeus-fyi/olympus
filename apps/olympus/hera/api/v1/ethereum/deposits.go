package hera_ethereum_validator_deposits

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/gochain/web3/web3_actions"
)

type ValidatorDepositDataRequest struct {
	web3_actions.SendTxPayload
}
type Response struct {
	Message string `json:"message"`
}

func (t *ValidatorDepositDataRequest) GenerateDepositData(c echo.Context) error {
	//ctx := context.Background()
	//ou := c.Get("orgUser").(org_users.OrgUser)
	//if err != nil {
	//	log.Err(err).Interface("orgUser", ou).Msg("SendEtherGoerliTx, ExecuteArtemisSendEthTxWorkflow error")
	//	return c.JSON(http.StatusBadRequest, nil)
	//}
	resp := Response{Message: "send eth tx in progress"}
	return c.JSON(http.StatusAccepted, resp)
}
