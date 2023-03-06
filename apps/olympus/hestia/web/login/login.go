package hestia_login

import (
	"context"
	"net/http"

	"github.com/google/uuid"
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

type LoginResponse struct {
	UserID    int       `json:"userID"`
	SessionID uuid.UUID `json:"sessionID"`
	TTL       int       `json:"ttl"`
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
	resp := LoginResponse{
		UserID:    key.UserID,
		SessionID: uuid.New(),
		TTL:       3600,
	}
	return c.JSON(http.StatusOK, resp)
}
