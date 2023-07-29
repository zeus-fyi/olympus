package hestia_quicknode_dashboard

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
		user := claims["name"].(string)
		organization := claims["organization_name"].(string)
		email := claims["email"].(string)
		ui := UserInfo{
			User:         user,
			Organization: organization,
			Email:        email,
		}
		resp.Status = "success"
		return c.JSON(http.StatusOK, ui)
	})
}

type UserInfo struct {
	User         string `json:"name"`
	Organization string `json:"organization"`
	Email        string `json:"email"`
}

type Response struct {
	Status string `json:"status"`
}

/*
	DashboardURL string `json:"dashboard-url"`
	AccessURL    string `json:"access-url"`
*/
