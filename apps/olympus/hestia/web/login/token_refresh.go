package hestia_login

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/keys"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	create_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/keys"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
	"k8s.io/apimachinery/pkg/util/rand"
)

type TokenRefreshRequest struct {
}

func TokenRefreshRequestHandler(c echo.Context) error {
	request := new(TokenRefreshRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.RefreshToken(c)
}

type TokenRefreshResponse struct {
	UserID    int    `json:"userID"`
	SessionID string `json:"sessionID"`
	TTL       int    `json:"ttl"`
}

func (l *TokenRefreshRequest) RefreshToken(c echo.Context) error {
	ctx := context.Background()

	ou := c.Get("orgUser").(org_users.OrgUser)
	key := read_keys.NewKeyReader()

	sessionID := rand.String(64)
	sessionKey := create_keys.NewCreateKey(ou.UserID, sessionID)
	sessionKey.PublicKeyVerified = true
	sessionKey.PublicKeyTypeID = keys.SessionIDKeyTypeID
	sessionKey.PublicKeyName = "sessionID"
	err := sessionKey.InsertUserSessionKey(ctx)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("InsertUserSessionKey error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	resp := LoginResponse{
		UserID:    key.UserID,
		SessionID: sessionID,
		TTL:       3600,
	}
	//cookie := &http.Cookie{
	//	Name:     "cookieName",
	//	Value:    sessionID,
	//	Path:     "/",
	//	HttpOnly: true,
	//	Secure:   true,
	//	SameSite: http.SameSiteLaxMode,
	//}
	//c.SetCookie(cookie)
	return c.JSON(http.StatusOK, resp)
}
