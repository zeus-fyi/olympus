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
	hestia_billing "github.com/zeus-fyi/olympus/hestia/web/billing"
	aegis_sessions "github.com/zeus-fyi/olympus/pkg/aegis/sessions"
	quicknode_orchestrations "github.com/zeus-fyi/olympus/pkg/hestia/platform/quiknode/orchestrations"
	"k8s.io/apimachinery/pkg/util/rand"
)

func GoogleLoginHandler(c echo.Context) error {
	request := new(GoogleLoginRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.VerifyGoogleLogin(c)
}

type GoogleLoginRequest struct {
	Credential string `json:"credential"`
	ClientId   string `json:"clientId"`
	SelectBy   string `json:"select_by"`
}

func (g *GoogleLoginRequest) VerifyGoogleLogin(c echo.Context) error {
	ctx := context.Background()
	key := read_keys.NewKeyReader()

	// TODO verify here

	// TODO get user info, then look up user by email
	sessionID := rand.String(64)
	sessionKey := create_keys.NewCreateKey(key.UserID, sessionID)
	sessionKey.PublicKeyVerified = true
	sessionKey.PublicKeyTypeID = keys.SessionIDKeyTypeID
	sessionKey.PublicKeyName = "sessionID"
	oldKey, err := sessionKey.InsertUserSessionKey(ctx)
	if err != nil {
		//log.Err(err).Interface("email", l.Email).Msg("InsertUserSessionKey error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	if oldKey != "" {
		err = quicknode_orchestrations.HestiaQnWorker.ExecuteDeleteSessionCacheWorkflowWorkflow(ctx, oldKey)
		if err != nil {
			log.Err(err).Msg("ExecuteDeleteSessionCacheWorkflowWorkflow error")
			err = nil
		}
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
		Domain:   Domain,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(24 * time.Hour),
	}
	c.SetCookie(cookie)
	pd, err := hestia_billing.GetPlan(ctx, sessionID)
	if err != nil {
		log.Err(err).Msg("GetPlan error")
		err = nil
	} else {
		resp.PlanDetailsUsage = &pd
	}
	return c.JSON(http.StatusOK, resp)
}
