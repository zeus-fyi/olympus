package v1

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/beacon_indexer/beacon_models"
)

type ValidatorsRequest struct {
	ValidatorIndexes []int64
}

func HandleValidatorsRequest(c echo.Context) error {
	log.Info().Msg("HandleValidatorsRequest")
	request := new(ValidatorsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	ctx := context.Background()
	vs := beacon_models.Validators{}
	v, err := vs.SelectValidators(ctx, request.ValidatorIndexes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, v.Validators)
}
