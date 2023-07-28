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
		resp := QuickNodeResponse{
			Status:       "",
			DashboardURL: "",
			AccessURL:    "",
		}
		token, ok := c.Get("user").(*jwt.Token) // by default token is stored under `user` key
		if !ok {
			resp.Status = "error: JWT token missing or invalid"
			return c.JSON(http.StatusBadRequest, resp)
		}
		// claims
		_, ok = token.Claims.(jwt.MapClaims) // by default claims is of type `jwt.MapClaims`
		if !ok {
			resp.Status = "error: failed to cast claims as jwt.MapClaims"
			return c.JSON(http.StatusBadRequest, resp)
		}
		/*
			c.JSON(200, gin.H{
				"status":        "success",
				"dashboard-url": scheme + c.Request.Host + "/dashboard",
				"access-url":    scheme + c.Request.Host + "/api", // Note: should be protected by API key
			})
		*/
		resp.Status = "success"
		return c.JSON(http.StatusOK, resp)
	})

}

type QuickNodeResponse struct {
	Status       string `json:"status"`
	DashboardURL string `json:"dashboard-url"`
	AccessURL    string `json:"access-url"`
}
