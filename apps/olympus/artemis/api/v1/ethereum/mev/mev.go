package artemis_eth_mev

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

type MempoolTxRequest struct {
	ProtocolID  int `json:"protocolID"`
	BlockNumber int `json:"blockNumber,omitempty"`
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
		return c.JSON(http.StatusOK, resp)
	default:
		resp, err := artemis_validator_service_groups_models.SelectMempoolTxAtBlockNumber(ctx, m.ProtocolID, m.BlockNumber)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("failed to get mempool txs")
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, resp)
	}
}
