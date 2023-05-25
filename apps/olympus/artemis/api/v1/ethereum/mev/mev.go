package artemis_eth_mev

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

type MempoolTxRequest struct {
	ProtocolID  int  `json:"protocolID"`
	BlockNumber int  `json:"blockNumber,omitempty"`
	ProfitOnly  bool `json:"profitOnly,omitempty"`
}

func MempoolTxRequestHandler(c echo.Context) error {
	request := new(MempoolTxRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetMempoolTxs(c)
}

func (m *MempoolTxRequest) GetMempoolTxs(c echo.Context) error {
	ctx := context.Background()
	if m.ProtocolID == 0 {
		m.ProtocolID = hestia_req_types.EthereumMainnetProtocolNetworkID
	}
	switch m.BlockNumber {
	case 0:
		resp, err := artemis_validator_service_groups_models.SelectMempoolTxAtMaxBlockNumber(ctx, m.ProtocolID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("failed to get mempool txs")
			return c.JSON(http.StatusInternalServerError, err)
		}
		if m.ProfitOnly {
			return FilterZeroProfit(c, ctx, resp)
		}
		return c.JSON(http.StatusOK, resp)
	default:
		resp, err := artemis_validator_service_groups_models.SelectMempoolTxAtBlockNumber(ctx, m.ProtocolID, m.BlockNumber)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("failed to get mempool txs")
			return c.JSON(http.StatusInternalServerError, err)
		}
		if m.ProfitOnly {
			return FilterZeroProfit(c, ctx, resp)
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func FilterZeroProfit(c echo.Context, ctx context.Context, resp artemis_autogen_bases.EthMempoolMevTxSlice) error {
	tmp := artemis_autogen_bases.EthMempoolMevTxSlice{}
	for _, mempoolTx := range resp {
		tf := web3_client.TradeExecutionFlow{}
		b := []byte(mempoolTx.TxFlowPrediction)
		berr := json.Unmarshal(b, &tf)
		if berr != nil {
			log.Ctx(ctx).Error().Err(berr).Msg("failed to unmarshal tx flow prediction")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		if tf.SandwichPrediction.ExpectedProfit != "0" {
			tmp = append(tmp, mempoolTx)
		}
	}
	return c.JSON(http.StatusOK, tmp)
}
