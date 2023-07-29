package hestia_quicknode_dashboard

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

var JWTAuthSecret = ""

// dashboardUrl := fmt.Sprintf("%s?jwt=%s", dashboardURL, jwtToken)

func InitQuickNodeDashboardRoutes(e *echo.Echo) {
	e.GET("/v1/quicknode/access", func(c echo.Context) error {
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
		ui := User{
			Name:             user,
			OrganizationName: organization,
			Email:            email,
			QuickNodeID:      quickNodeID,
		}
		resp.Status = "success"
		return c.JSON(http.StatusOK, ui)
	})

	e.GET("/v1/quicknode/dashboard", func(c echo.Context) error {
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
		ui := User{
			Name:             user,
			OrganizationName: organization,
			Email:            email,
			QuickNodeID:      quickNodeID,
		}
		resp.Status = "success"
		return c.JSON(http.StatusOK, ui)
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
