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
	aegis_sessions "github.com/zeus-fyi/olympus/pkg/aegis/sessions"
	quicknode_orchestrations "github.com/zeus-fyi/olympus/pkg/hestia/platform/quiknode/orchestrations"
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
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Warn().Msg("RefreshTokenRequestHandler: orgUser not found")
		return c.JSON(http.StatusUnauthorized, nil)
	}
	key := read_keys.NewKeyReader()

	sessionID := rand.String(64)
	sessionKey := create_keys.NewCreateKey(ou.UserID, sessionID)
	sessionKey.PublicKeyVerified = true
	sessionKey.PublicKeyTypeID = keys.SessionIDKeyTypeID
	sessionKey.PublicKeyName = "sessionID"
	oldKey, err := sessionKey.InsertUserSessionKey(ctx)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("InsertUserSessionKey error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	if oldKey != "" {
		err = quicknode_orchestrations.HestiaQnWorker.ExecuteDeleteSessionCacheWorkflowWorkflow(ctx, oldKey)
		if err != nil {
			log.Err(err).Msg("ExecuteDeleteSessionCacheWorkflowWorkflow error")
			err = nil
		}
	}
	isInternal := false
	if key.OrgID == TemporalOrgID || key.OrgID == SamsOrgID {
		isInternal = true
	}
	resp := LoginResponse{
		UserID:     key.UserID,
		SessionID:  sessionID,
		IsInternal: isInternal,
		TTL:        3600,
	}

	cookie := &http.Cookie{
		Name:     "cookieName",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, resp)
}

type UserAuthedServicesRequest struct {
}

func UsersServicesRequestHandler(c echo.Context) error {
	request := new(UserAuthedServicesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetAuthedServices(c)
}

func (l *UserAuthedServicesRequest) GetAuthedServices(c echo.Context) error {
	ctx := context.Background()
	sessionToken := ""
	cookie, err := c.Cookie(aegis_sessions.SessionIDNickname)
	if err == nil && cookie != nil {
		sessionToken = cookie.Value
	}
	k := read_keys.NewKeyReader()
	services, _, err := k.QueryUserAuthedServices(ctx, sessionToken)
	if err != nil {
		log.Err(err).Msg("InitV1Routes: QueryUserAuthedServices error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	var resp []string
	for _, service := range services {
		if service == "ethereumGoerliValidators" {
			resp = append(resp, "Goerli")
		}
		if service == "ethereumEphemeryValidators" {
			resp = append(resp, "Ephemery")
		}
	}
	return c.JSON(http.StatusOK, resp)
}
