package athena_ethereum_import_validators

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type ImportLighthouseValidatorsRequest struct {
}

func ImportLighthouseValidatorsHandler(c echo.Context) error {
	request := new(ImportLighthouseValidatorsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.ImportValidators(c)
}

func (lh *ImportLighthouseValidatorsRequest) ImportValidators(c echo.Context) error {
	log.Info().Msg("ImportLighthouseValidatorsRequest: ImportValidators")

	log.Info().Msg("ImportLighthouseValidatorsRequest: Import Sync Finished")
	return c.JSON(http.StatusOK, nil)
}
