package hestia_login

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(c echo.Context) error {
	request := new(LoginRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.VerifyPassword(c)
}

func (l *LoginRequest) VerifyPassword(c echo.Context) error {
	ctx := context.Background()
	key := read_keys.NewKeyReader()
	key.PublicKey = l.Password
	err := key.VerifyUserPassword(ctx, l.Email)
	if err != nil {
		log.Err(err).Interface("email", l.Email).Msg("VerifyPassword error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, nil)
}
