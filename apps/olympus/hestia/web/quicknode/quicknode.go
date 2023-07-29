package hestia_quicknode_dashboard

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

var JWTAuthSecret = ""

func InitQuickNodeDashboardRoutes(e *echo.Echo) {
	eg := e.Group("/v1/quicknode")
	eg.Use(echojwt.JWT(middleware.JWTConfig{
		SigningKey: []byte(JWTAuthSecret),
	}))
	e.GET("/dashboard", func(c echo.Context) error {
		resp := Response{
			Status: "",
		}
		token, ok := c.Get("jwt").(*jwt.Token)
		if !ok {
			resp.Status = "error: JWT token missing or invalid"
			return c.JSON(http.StatusBadRequest, resp)
		}

		claims, ok := token.Claims.(jwt.MapClaims) // by default claims is of type `jwt.MapClaims`
		if !ok {
			resp.Status = "error: failed to cast claims as jwt.MapClaims"
			return c.JSON(http.StatusBadRequest, resp)
		}
		user, ok1 := claims["name"].(string)
		if !ok1 {
			log.Warn().Msg("failed to cast claims[\"name\"] as string")
		}
		organization, ok2 := claims["organization_name"].(string)
		if !ok2 {
			log.Warn().Msg("failed to cast claims[\"organization_name\"] as string")
		}
		email, ok3 := claims["email"].(string)
		if !ok3 {
			log.Warn().Msg("failed to cast claims[\"email\"] as string")
		}
		quickNodeID, ok4 := claims["quicknode_id"].(string)
		if !ok4 {
			log.Warn().Msg("failed to cast claims[\"quicknode_id\"] as string")
		}
		ui := UserInfo{
			User:         user,
			Organization: organization,
			Email:        email,
			QuickNodeID:  quickNodeID,
		}
		resp.Status = "success"
		return c.JSON(http.StatusOK, ui)
	})
}

type UserInfo struct {
	User         string `json:"name"`
	Organization string `json:"organization"`
	Email        string `json:"email"`
	QuickNodeID  string `json:"quickNodeID"`
}

type Response struct {
	Status string `json:"status"`
}
