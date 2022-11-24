package athena_jwt

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	v1_common_routes "github.com/zeus-fyi/olympus/athena/api/v1/common"
	athena_jwt "github.com/zeus-fyi/olympus/athena/pkg/jwt"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

type TokenRequestJWT struct {
	JWT string
}

func (t *TokenRequestJWT) Create(c echo.Context) error {
	err := athena_jwt.SetToken(v1_common_routes.CommonManager.DataDir, t.JWT)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

func ReplaceToken(p filepaths.Path, token string) error {
	err := p.RemoveFileInPath()
	if err != nil {
		log.Err(err).Msg("error removing jwt token")
		return err
	}
	return SetToken(p, token)
}

func SetToken(p filepaths.Path, token string) error {
	p.DirOut = p.DirIn
	err := p.WriteToFileOutPath([]byte(token))
	if err != nil {
		log.Err(err).Msg("error setting jwt token")
	}
	return err
}

func SetTokenToDefault(p filepaths.Path, tokenFileName, defaultToken string) error {
	p.FnIn = tokenFileName
	p.FnOut = tokenFileName
	if !p.FileInPathExists() {
		err := SetToken(p, defaultToken)
		return err
	} else {
		return ReplaceToken(p, defaultToken)
	}
}
