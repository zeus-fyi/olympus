package hestia_signup

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type VerifyEmailRequest struct {
}

func VerifyEmailHandler(c echo.Context) error {
	request := new(VerifyEmailRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.VerifyEmail(c)
}

func (v *VerifyEmailRequest) VerifyEmail(c echo.Context) error {
	c.Param("token")

	return c.JSON(http.StatusOK, nil)
}
