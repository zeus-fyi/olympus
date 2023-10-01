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
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2api "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"k8s.io/apimachinery/pkg/util/rand"
)

var GoogleOAuthConfig = &oauth2.Config{
	Endpoint: google.Endpoint,
}

func GoogleLoginHandler(c echo.Context) error {
	request := new(GoogleLoginRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.VerifyGoogleLogin(c)
}

type GoogleLoginRequest struct {
	Credential string `json:"credential"`
}

func (g *GoogleLoginRequest) VerifyGoogleLogin(c echo.Context) error {
	ctx := context.Background()
	key := read_keys.NewKeyReader()

	oauth2Service, err := oauth2api.NewService(ctx, option.WithTokenSource(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: GoogleOAuthConfig.ClientSecret})))
	if err != nil {
		log.Err(err).Msg("VerifyGoogleLogin: oauth2api.NewService error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.AccessToken(g.Credential)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		log.Err(err).Msg("VerifyGoogleLogin: tokenInfoCall.Do error")
		return c.JSON(http.StatusUnauthorized, nil)
	}
	// Ensure that the token is issued to your application
	if tokenInfo.Audience != GoogleOAuthConfig.ClientID {
		log.Err(err).Msg("VerifyGoogleLogin: tokenInfo.Audience != clientID")
		return c.JSON(http.StatusUnauthorized, nil)
	}
	err = key.GetUserFromEmail(ctx, tokenInfo.Email)
	if err != nil {
		log.Err(err).Interface("email", tokenInfo.Email).Msg("GetUserFromEmail error")
		return c.JSON(http.StatusUnauthorized, nil)
	}
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
