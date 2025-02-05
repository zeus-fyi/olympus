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
	iris_service_plans "github.com/zeus-fyi/olympus/iris/api/v1/service_plans"
	aegis_sessions "github.com/zeus-fyi/olympus/pkg/aegis/sessions"
	quicknode_orchestrations "github.com/zeus-fyi/olympus/pkg/hestia/platform/quiknode/orchestrations"
	"k8s.io/apimachinery/pkg/util/rand"
)

var Domain string

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

	IsBillingSetup   bool                                         `json:"isBillingSetup,omitempty"`
	IsInternal       bool                                         `json:"isInternal,omitempty"`
	PlanDetailsUsage *iris_service_plans.PlanUsageDetailsResponse `json:"planUsageDetails,omitempty"`
}

const (
	TemporalOrgID = 7138983863666903883
	SamsOrgID     = 1701381301753642000
)

func (l *LoginRequest) VerifyPassword(c echo.Context) error {
	ctx := context.Background()
	key := read_keys.NewKeyReader()
	key.PublicKey = l.Password
	err := key.VerifyUserPassword(ctx, l.Email)
	if err != nil {
		log.Err(err).Interface("email", l.Email).Msg("VerifyPassword error")
		return c.JSON(http.StatusUnauthorized, nil)
	}
	if key.PublicKeyVerified == false {
		log.Err(err).Interface("email", l.Email).Msg("VerifyPassword StatusUnauthorized")
		return c.JSON(http.StatusUnauthorized, nil)
	}
	sessionID := rand.String(64)
	sessionKey := create_keys.NewCreateKey(key.UserID, sessionID)
	sessionKey.PublicKeyVerified = true
	sessionKey.PublicKeyTypeID = keys.SessionIDKeyTypeID
	sessionKey.PublicKeyName = "sessionID"
	oldKey, err := sessionKey.InsertUserSessionKey(ctx)
	if err != nil {
		log.Err(err).Interface("email", l.Email).Msg("InsertUserSessionKey error")
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
		UserID:         key.UserID,
		SessionID:      sessionID,
		IsInternal:     isInternal,
		IsBillingSetup: hestia_billing.CheckBillingCache(ctx, key.UserID),
		TTL:            3600,
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
