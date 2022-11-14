package athena_jwt_route

import (
	"net/http"

	"github.com/labstack/echo/v4"
	v1_common_routes "github.com/zeus-fyi/olympus/athena/api/v1/common"
	athena_jwt "github.com/zeus-fyi/olympus/athena/pkg/jwt"
)

type TokenRequestJWT struct {
	JWT string
}

func (t *TokenRequestJWT) Create(c echo.Context) error {
	err := athena_jwt.SetToken(v1_common_routes.DataDir, t.JWT)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}
