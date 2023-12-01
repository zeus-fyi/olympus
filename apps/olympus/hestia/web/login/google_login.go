package hestia_login

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/keys"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	create_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/keys"
	create_org_users "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/org_users"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
	hestia_billing "github.com/zeus-fyi/olympus/hestia/web/billing"
	aegis_sessions "github.com/zeus-fyi/olympus/pkg/aegis/sessions"
	quicknode_orchestrations "github.com/zeus-fyi/olympus/pkg/hestia/platform/quiknode/orchestrations"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
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
	pl, err := idtoken.Validate(ctx, g.Credential, GoogleOAuthConfig.ClientID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, nil)
	}
	if pl.Audience != GoogleOAuthConfig.ClientID {
		return c.JSON(http.StatusUnauthorized, nil)
	}
	claims := pl.Claims
	email, ok1 := claims["email"].(string)
	if !ok1 {
		return c.JSON(http.StatusUnauthorized, nil)
	}
	name, ok2 := claims["name"].(string)
	if !ok2 {
		return c.JSON(http.StatusUnauthorized, nil)
	}
	lastName, ok2 := claims["family_name"].(string)
	if !ok2 {
		return c.JSON(http.StatusUnauthorized, nil)
	}
	firstName := ""
	fullName := strings.Split(name, " ")
	if len(fullName) > 0 {
		firstName = fullName[0]
	}
	us := create_org_users.UserSignup{
		FirstName:        firstName,
		LastName:         lastName,
		EmailAddress:     email,
		Password:         rand.String(64),
		VerifyEmailToken: "",
	}
	err = key.GetUserFromEmail(ctx, email)
	if err == pgx.ErrNoRows {
		ou := create_org_users.OrgUser{}
		verifyToken, verr := ou.InsertSignUpOrgUserAndVerifyEmail(ctx, us)
		if verr != nil {
			log.Err(verr).Interface("user", us).Msg("SignupRequest, SignUp error")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		us.VerifyEmailToken = verifyToken
		err = create_keys.UpdateKeysFromVerifyEmail(ctx, us.VerifyEmailToken)
		if err != nil {
			log.Err(err).Str("token", us.VerifyEmailToken).Msg("VerifyEmailHandler, VerifyEmail error")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		err = key.GetUserFromEmail(ctx, email)

		ou2 := org_users.OrgUser{}
		ou2.UserID = key.UserID
		ou2.OrgID = key.OrgID
		_, err2 := create_keys.CreateUserAPIKey(ctx, ou2)
		if err2 != nil {
			log.Err(err2).Msg("CreateUserAPIKey error")
			err2 = nil
		}
	}
	if err != nil {
		log.Err(err).Interface("email", email).Msg("GetUserFromEmail error")
		return c.JSON(http.StatusUnauthorized, nil)
	}
	sessionID := rand.String(64)
	sessionKey := create_keys.NewCreateKey(key.UserID, sessionID)
	sessionKey.PublicKeyVerified = true
	sessionKey.PublicKeyTypeID = keys.SessionIDKeyTypeID
	sessionKey.PublicKeyName = "sessionID"
	oldKey, err := sessionKey.InsertUserSessionKey(ctx)
	if err != nil {
		log.Err(err).Interface("email", email).Msg("InsertUserSessionKey error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	if oldKey != "" {
		err = quicknode_orchestrations.HestiaQnWorker.ExecuteDeleteSessionCacheWorkflowWorkflow(ctx, oldKey)
		if err != nil {
			log.Err(err).Msg("ExecuteDeleteSessionCacheWorkflowWorkflow error")
			err = nil
		}
	}
	cookie := &http.Cookie{
		Name:     aegis_sessions.SessionIDNickname,
		Value:    sessionID,
		HttpOnly: true,
		Secure:   true,
		Domain:   Domain,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
	}
	c.SetCookie(cookie)
	if (key.UserID == 0) || (key.OrgID == 0) {

	}
	isInternal := false
	if key.OrgID == TemporalOrgID || key.OrgID == SamsOrgID {
		isInternal = true
	}
	li := LoginResponse{
		UserID:     key.UserID,
		SessionID:  sessionID,
		IsInternal: isInternal,
		TTL:        3600,
	}
	pd, err := hestia_billing.GetPlan(ctx, sessionID)
	if err != nil {
		log.Err(err).Msg("GetPlan error")
		err = nil
	} else {
		li.PlanDetailsUsage = &pd
	}
	return c.JSON(http.StatusOK, li)
}
