package hestia_signup

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	create_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/keys"
)

type VerifyEmailRequest struct {
}

func VerifyEmailHandler(c echo.Context) error {
	request := new(VerifyEmailRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.VerifyEmail(c)
}

func (v *VerifyEmailRequest) VerifyEmail(c echo.Context) error {
	token := c.Param("token")
	ctx := context.Background()
	err := create_keys.UpdateKeysFromVerifyEmail(ctx, token)
	if err != nil {
		log.Err(err).Str("token", token).Msg("VerifyEmailHandler, VerifyEmail error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}
