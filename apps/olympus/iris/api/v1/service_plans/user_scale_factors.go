package iris_service_plans

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const (
	UpdateTableLatencyScaleFactor = "/table/:groupName/scale/latency"
	UpdateTableErrorScaleFactor   = "/table/:groupName/scale/error"
	UpdateTableDecayScaleFactor   = "/table/:groupName/scale/decay"
)

type ScaleFactorsUpdateRequest struct {
}

func UpdateLatencyScaleFactorRequestHandler(c echo.Context) error {
	log.Info().Msg("Iris: UpdateLatencyScaleFactorRequestHandler")
	request := new(ScaleFactorsUpdateRequest)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("Iris: UpdateLatencyScaleFactorRequestHandler")
		return err
	}
	return request.UpdateLatencyScaleFactor(c)
}

func (p *ScaleFactorsUpdateRequest) UpdateLatencyScaleFactor(c echo.Context) error {
	return nil
}

func UpdateErrorScaleFactorRequestHandler(c echo.Context) error {
	log.Info().Msg("Iris: UpdateErrorScaleFactorRequestHandler")
	request := new(ScaleFactorsUpdateRequest)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("Iris: UpdateErrorScaleFactorRequestHandler")
		return err
	}
	return request.UpdateErrorScaleFactor(c)
}

func (p *ScaleFactorsUpdateRequest) UpdateErrorScaleFactor(c echo.Context) error {
	return nil
}

func UpdateDecayScaleFactorRequestHandler(c echo.Context) error {
	log.Info().Msg("Iris: UpdateDecayScaleFactorRequestHandler")
	request := new(ScaleFactorsUpdateRequest)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("Iris: UpdateDecayScaleFactorRequestHandler")
		return err
	}
	return request.UpdateDecayScaleFactor(c)
}

func (p *ScaleFactorsUpdateRequest) UpdateDecayScaleFactor(c echo.Context) error {
	return nil
}
