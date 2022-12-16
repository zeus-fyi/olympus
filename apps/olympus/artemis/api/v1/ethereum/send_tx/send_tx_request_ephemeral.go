package artemis_eth_txs

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	artemis_ethereum_transcations "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/transcations"
)

func (t *EthereumSendEtherRequest) SendEtherEphemeral(c echo.Context) error {
	ctx := context.Background()
	err := artemis_ethereum_transcations.ArtemisEthereumEphemeralTxBroadcastWorker.ExecuteArtemisSendEthTxWorkflow(ctx, t.SendEtherPayload)
	ou := c.Get("orgUser").(org_users.OrgUser)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("SendEtherEphemeral, ExecuteArtemisSendEthTxWorkflow error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	resp := Response{Message: "send eth tx in progress"}
	return c.JSON(http.StatusAccepted, resp)
}

func (t *EthereumSendSignedTxRequest) SendEphemeralSignedTx(c echo.Context) error {
	ctx := context.Background()
	err := artemis_ethereum_transcations.ArtemisEthereumEphemeralTxBroadcastWorker.ExecuteArtemisSendSignedTxWorkflow(ctx, &t.Transaction)
	ou := c.Get("orgUser").(org_users.OrgUser)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("SendEphemeralSignedTx, ExecuteArtemisSendSignedTxWorkflow error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	resp := Response{Message: "tx in progress"}
	return c.JSON(http.StatusAccepted, resp)
}
