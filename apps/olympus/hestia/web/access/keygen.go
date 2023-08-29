package hestia_access_keygen

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	create_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/keys"
)

type AccessKeyGenRequest struct {
}

func AccessKeyGenRequestHandler(c echo.Context) error {
	request := new(AccessKeyGenRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.KeyGen(c)
}

type AccessKeyGenResp struct {
	ApiKeyName   string `json:"apiKeyName"`
	ApiKeySecret string `json:"apiKeySecret"`
}

func (a *AccessKeyGenRequest) KeyGen(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	ctx := context.Background()
	apiKey, err := create_keys.CreateUserAPIKey(ctx, ou)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("CreateUserAPIKey error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	resp := AccessKeyGenResp{
		ApiKeyName:   apiKey.PublicKeyName,
		ApiKeySecret: apiKey.PublicKey,
	}
	return c.JSON(http.StatusOK, resp)
}
