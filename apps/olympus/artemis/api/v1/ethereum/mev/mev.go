package artemis_eth_mev

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
)

type MempoolTxRequest struct {
	ProtocolID  int `json:"protocolID"`
	BlockNumber int `json:"blockNumber"`
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
	resp, err := artemis_validator_service_groups_models.SelectMempoolTxAtBlockNumber(ctx, m.ProtocolID, m.BlockNumber)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, resp)
}
