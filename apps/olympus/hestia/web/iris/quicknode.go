package hestia_quicknode_dashboard

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/keys"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	create_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/keys"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
	hestia_login "github.com/zeus-fyi/olympus/hestia/web/login"
	aegis_sessions "github.com/zeus-fyi/olympus/pkg/aegis/sessions"
	"k8s.io/apimachinery/pkg/util/rand"
)

var JWTAuthSecret = ""

func InitQuickNodeDashboardRoutes(e *echo.Echo) {
	e.GET("/quicknode/access", func(c echo.Context) error {
		resp := Response{
			Status: "",
		}
		token, err := jwt.Parse(c.QueryParam("jwt"), func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(JWTAuthSecret), nil
		})
		if err != nil {
			resp.Status = "error: failed to parse jwt token"
			return c.JSON(http.StatusBadRequest, resp)
		}
		claims, ok := token.Claims.(jwt.MapClaims) // by default claims is of type `jwt.MapClaims`
		if !ok {
			resp.Status = "error: failed to cast claims as jwt.MapClaims"
			return c.JSON(http.StatusBadRequest, resp)
		}
		user, ok1 := claims["name"].(string)
		if !ok1 {
			log.Warn().Interface("user", user).Msg("failed to cast claims[\"user\"] as string")
			log.Warn().Msg("failed to cast claims[\"name\"] as string")
		}
		organization, ok2 := claims["organization_name"].(string)
		if !ok2 {
			log.Warn().Interface("organization", organization).Msg("failed to cast claims[\"organization_name\"] as string")
			log.Warn().Msg("failed to cast claims[\"organization_name\"] as string")
		}
		email, ok3 := claims["email"].(string)
		if !ok3 {
			log.Warn().Interface("email", email).Msg("failed to cast claims[\"email\"] as string")
			log.Warn().Msg("failed to cast claims[\"email\"] as string")
		}
		quickNodeID, ok4 := claims["quicknode_id"].(string)
		if !ok4 {
			log.Warn().Msg("failed to cast claims[\"quicknode_id\"] as string")
		}
		//ui := User{
		//	Name:             user,
		//	OrganizationName: organization,
		//	Email:            email,
		//	QuickNodeID:      quickNodeID,
		//}
		resp.Status = "success"
		key := read_keys.NewKeyReader()
		key.PublicKey = quickNodeID
		ctx := context.Background()
		err = key.VerifyUserTokenServiceWithQuickNodePlan(ctx, quickNodeID)
		if err != nil {
			log.Err(err).Str("quickNodeID", quickNodeID).Msg("VerifyUserTokenServiceWithQuickNodePlan error")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		if key.PublicKeyVerified == false {
			return c.JSON(http.StatusInternalServerError, nil)
		}
		sessionID := rand.String(64)
		sessionKey := create_keys.NewCreateKey(key.UserID, sessionID)
		sessionKey.PublicKeyTypeID = keys.SessionIDKeyTypeID
		sessionKey.PublicKeyName = "sessionID"
		err = sessionKey.InsertUserSessionKey(ctx)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
		cookie := &http.Cookie{
			Name:     aegis_sessions.SessionIDNickname,
			Value:    sessionID,
			HttpOnly: true,
			Secure:   true,
			Domain:   hestia_login.Domain,
			SameSite: http.SameSiteNoneMode,
			Expires:  time.Now().Add(24 * time.Hour),
			Path:     "/",
		}
		c.SetCookie(cookie)
		li := hestia_login.LoginResponse{
			UserID:    key.UserID,
			SessionID: sessionID,
			TTL:       3600,
		}
		ou := org_users.NewOrgUser()
		ou.OrgID = key.OrgID

		pu, err := GetUserPlanInfo(ctx, ou, "test")
		if err != nil {
			log.Err(err).Msg("GetUserPlanInfo error")
		} else {
			li.PlanDetailsUsage = pu
		}
		return c.JSON(http.StatusOK, li)
	})

	e.GET("/quicknode/dashboard", func(c echo.Context) error {
		resp := Response{
			Status: "",
		}
		token, err := jwt.Parse(c.QueryParam("jwt"), func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(JWTAuthSecret), nil
		})
		if err != nil {
			resp.Status = "error: failed to parse jwt token"
			return c.JSON(http.StatusBadRequest, resp)
		}
		claims, ok := token.Claims.(jwt.MapClaims) // by default claims is of type `jwt.MapClaims`
		if !ok {
			resp.Status = "error: failed to cast claims as jwt.MapClaims"
			return c.JSON(http.StatusBadRequest, resp)
		}
		user, ok1 := claims["name"].(string)
		if !ok1 {
			log.Warn().Interface("user", user).Msg("failed to cast claims[\"user\"] as string")
			log.Warn().Msg("failed to cast claims[\"name\"] as string")
		}
		organization, ok2 := claims["organization_name"].(string)
		if !ok2 {
			log.Warn().Interface("organization", organization).Msg("failed to cast claims[\"organization_name\"] as string")
			log.Warn().Msg("failed to cast claims[\"organization_name\"] as string")
		}
		email, ok3 := claims["email"].(string)
		if !ok3 {
			log.Warn().Interface("email", email).Msg("failed to cast claims[\"email\"] as string")
			log.Warn().Msg("failed to cast claims[\"email\"] as string")
		}
		quickNodeID, ok4 := claims["quicknode_id"].(string)
		if !ok4 {
			log.Warn().Msg("failed to cast claims[\"quicknode_id\"] as string")
		}
		//ui := User{
		//	Name:             user,
		//	OrganizationName: organization,
		//	Email:            email,
		//	QuickNodeID:      quickNodeID,
		//}
		resp.Status = "success"
		key := read_keys.NewKeyReader()
		key.PublicKey = quickNodeID
		ctx := context.Background()
		err = key.VerifyUserTokenServiceWithQuickNodePlan(ctx, quickNodeID)
		if err != nil {
			log.Err(err).Str("quickNodeID", quickNodeID).Msg("VerifyUserTokenServiceWithQuickNodePlan error")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		if key.PublicKeyVerified == false {
			return c.JSON(http.StatusInternalServerError, nil)
		}
		sessionID := rand.String(64)
		sessionKey := create_keys.NewCreateKey(key.UserID, sessionID)
		sessionKey.PublicKeyTypeID = keys.SessionIDKeyTypeID
		sessionKey.PublicKeyName = "sessionID"
		err = sessionKey.InsertUserSessionKey(ctx)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
		cookie := &http.Cookie{
			Name:     aegis_sessions.SessionIDNickname,
			Value:    sessionID,
			HttpOnly: true,
			Secure:   true,
			Domain:   hestia_login.Domain,
			SameSite: http.SameSiteNoneMode,
			Expires:  time.Now().Add(24 * time.Hour),
			Path:     "/",
		}
		c.SetCookie(cookie)
		li := hestia_login.LoginResponse{
			UserID:    key.UserID,
			SessionID: sessionID,
			TTL:       3600,
		}
		ou := org_users.NewOrgUser()
		ou.OrgID = key.OrgID
		pu, err := GetUserPlanInfo(ctx, ou, "test")
		if err != nil {
			log.Err(err).Msg("GetUserPlanInfo error")
		} else {
			li.PlanDetailsUsage = pu
		}
		return c.JSON(http.StatusOK, li)
	})
}

type User struct {
	Name             string `json:"name"`
	Email            string `json:"email"`
	OrganizationName string `json:"organization_name"`
	QuickNodeID      string `json:"quicknode-id"`
}

type Response struct {
	Status string `json:"status"`
}
