package hestia_login

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/keys"
	create_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/keys"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
	aegis_sessions "github.com/zeus-fyi/olympus/pkg/aegis/sessions"
	"k8s.io/apimachinery/pkg/util/rand"
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
	UserID    int    `json:"userID"`
	SessionID string `json:"sessionID"`
	TTL       int    `json:"ttl"`
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

	sessionID := rand.String(64)
	sessionKey := create_keys.NewCreateKey(key.UserID, sessionID)
	sessionKey.PublicKeyVerified = true
	sessionKey.PublicKeyTypeID = keys.SessionIDKeyTypeID
	sessionKey.PublicKeyName = "sessionID"
	err = sessionKey.InsertUserSessionKey(ctx)
	if err != nil {
		log.Err(err).Interface("email", l.Email).Msg("InsertUserSessionKey error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	resp := LoginResponse{
		UserID:    key.UserID,
		SessionID: sessionID,
		TTL:       3600,
	}
	cookie := &http.Cookie{
		Name:     aegis_sessions.SessionIDNickname,
		Value:    sessionID,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
	}
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, resp)
}
