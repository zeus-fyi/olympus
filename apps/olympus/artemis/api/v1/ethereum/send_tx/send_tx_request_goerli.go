package artemis_eth_txs

import (
	"context"
	"math/big"
	"net/http"

	"github.com/gochain/gochain/v4/core/types"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/web3_actions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	artemis_ethereum_transcations "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/transcations"
)

type EthereumSendSignedTxRequest struct {
	*types.Transaction
}

type EthereumSendEtherRequest struct {
	web3_actions.SendEtherPayload
}

type GasPriceLimits struct {
	GasPrice *big.Int
	GasLimit uint64
}

type Response struct {
	Message string `json:"message"`
}

func (t *EthereumSendEtherRequest) SendEtherGoerliTx(c echo.Context) error {
	ctx := context.Background()
	err := artemis_ethereum_transcations.ArtemisEthereumGoerliTxBroadcastWorker.ExecuteArtemisSendEthTxWorkflow(ctx, t.SendEtherPayload)
	ou := c.Get("orgUser").(org_users.OrgUser)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("SendEtherGoerliTx, ExecuteArtemisSendEthTxWorkflow error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	resp := Response{Message: "send eth tx in progress"}
	return c.JSON(http.StatusAccepted, resp)
}

func (t *EthereumSendSignedTxRequest) SendGoerliSignedTx(c echo.Context) error {
	ctx := context.Background()
	err := artemis_ethereum_transcations.ArtemisEthereumGoerliTxBroadcastWorker.ExecuteArtemisSendSignedTxWorkflow(ctx, t.Transaction)
	ou := c.Get("orgUser").(org_users.OrgUser)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("SendGoerliSignedTx, ExecuteArtemisSendSignedTxWorkflow error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	resp := Response{Message: "tx in progress"}
	return c.JSON(http.StatusAccepted, resp)
}
