package hestia_access_keygen

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func SecretsRequestHandler(c echo.Context) error {
	request := new(SecretsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateOrUpdateSecret(c)
}

type SecretsRequest struct {
	Name  string `json:"name"`
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
}

func (a *SecretsRequest) CreateOrUpdateSecret(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Warn().Interface("ou", ou).Msg("CreateOrUpdateSecret: orgUser not found")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

func SecretReadRequestHandler(c echo.Context) error {
	return RetrieveSecretValue(c)
}

func RetrieveSecretValue(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Warn().Interface("ou", ou).Msg("CreateOrUpdateSecret: orgUser not found")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	//ctx := context.Background()
	ref := c.Param("ref")
	if len(ref) <= 0 {
		return c.JSON(http.StatusBadRequest, "ref is required")
	}
	return c.JSON(http.StatusOK, SecretsRequest{
		Name:  "te",
		Key:   "aa",
		Value: "wfdsfsd",
	})
}

func SecretsReadRequestHandler(c echo.Context) error {
	request := new(SecretsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.ReadSecretReferences(c)
}

func (a *SecretsRequest) ReadSecretReferences(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Warn().Interface("ou", ou).Msg("ReadSecretReferences: orgUser not found")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, []SecretsRequest{
		{
			Name: "test",
			Key:  "test",
		},
		{
			Name: "test2",
			Key:  "test2",
		},
	})
}

func SecretDeleteRequestHandler(c echo.Context) error {
	return DeleteSecretValue(c)
}

func DeleteSecretValue(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Warn().Interface("ou", ou).Msg("ReadSecretReferences: orgUser not found")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	//ctx := context.Background()
	ref := c.Param("ref")
	if len(ref) <= 0 {
		return c.JSON(http.StatusBadRequest, "ref is required")
	}
	return c.JSON(http.StatusOK, nil)
}
