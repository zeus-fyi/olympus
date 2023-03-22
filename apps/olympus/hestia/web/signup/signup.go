package hestia_signup

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	create_org_users "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/org_users"
	hermes_email_notifications "github.com/zeus-fyi/olympus/pkg/hermes/email"
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
	ctx := context.Background()
	ou := create_org_users.OrgUser{}
	us := create_org_users.UserSignup{
		FirstName:    s.FirstName,
		LastName:     s.FirstName,
		EmailAddress: s.Email,
		Password:     s.Password,
	}
	verifyToken, err := ou.InsertSignUpOrgUserAndVerifyEmail(ctx, us)
	if err != nil {
		log.Err(err).Interface("user", us).Msg("SignupRequest, SignUp error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	us.VerifyEmailToken = verifyToken
	_, err = hermes_email_notifications.Hermes.SendSendGridEmailVerifyRequest(ctx, us)
	if err != nil {
		log.Err(err).Interface("user", us).Msg("SignupRequest, SignUp error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}
