package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/datastores/postgres/beacon_indexer/beacon_models"
)

type ValidatorBalancesRequest struct {
	ValidatorIndexes []int64
	LowerEpoch       int
	HigherEpoch      int
}

func HandleValidatorBalancesRequest(c echo.Context) error {
	log.Info().Msg("HandleValidatorBalancesRequest")
	request := new(ValidatorBalancesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	ctx := context.Background()
	vbe := beacon_models.ValidatorBalancesEpoch{}
	vb, err := vbe.SelectValidatorBalances(ctx, request.LowerEpoch, request.HigherEpoch, request.ValidatorIndexes)
	if err != nil || vb == nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	s, err := json.Marshal(vb.GetRawRowValues())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, string(s))
}

func HandleValidatorBalancesSumRequest(c echo.Context) error {
	log.Info().Msg("HandleValidatorBalancesSumRequest")
	request := new(ValidatorBalancesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	ctx := context.Background()
	vbe := beacon_models.ValidatorBalancesEpoch{}
	vb, err := vbe.SelectValidatorBalancesSum(ctx, request.LowerEpoch, request.HigherEpoch, request.ValidatorIndexes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	vb.LowerEpoch = request.LowerEpoch
	vb.HigherEpoch = request.HigherEpoch
	s, err := json.Marshal(vb)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, string(s))
}
