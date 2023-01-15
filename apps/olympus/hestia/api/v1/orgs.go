package v1hestia

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	create_orgs "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/orgs"
)

type CreateOrgRequest struct {
	Name     string `db:"name" json:"name"`
	Metadata string `db:"metadata" json:"metadata,omitempty"`
}

func CreateOrgHandler(c echo.Context) error {
	request := new(CreateOrgRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateOrg(c)
}

func (o *CreateOrgRequest) CreateOrg(c echo.Context) error {
	ctx := context.Background()
	org := create_orgs.NewCreateNamedOrg(o.Name)
	err := org.InsertOrg(ctx)
	if err != nil {
		log.Err(err).Interface("org", o).Msg("CreateOrgRequest, CreateOrg error")
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, org)
}
