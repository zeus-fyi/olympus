package hestia_signup

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type SignupRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func SignUpHandler(c echo.Context) error {
	request := new(SignupRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.SignUp(c)
}

func (s *SignupRequest) SignUp(c echo.Context) error {
	//ctx := context.Background()

	// TODO, insert user into db, setup service, send email to verify
	// create verify link -> user clicks link -> activate user
	return c.JSON(http.StatusOK, nil)
}
