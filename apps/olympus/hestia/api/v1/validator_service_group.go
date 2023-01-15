package v1hestia

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/validator_service_group"
)

type CreateValidatorServiceRequest struct {
	hestia_autogen_bases.ValidatorServiceOrgGroupSlice
}

func CreateValidatorServiceRequestHandler(c echo.Context) error {
	request := new(CreateValidatorServiceRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateValidatorsServiceGroup(c)
}

func (v *CreateValidatorServiceRequest) CreateValidatorsServiceGroup(c echo.Context) error {
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	err := validator_service_group.InsertValidatorServiceOrgGroup(ctx, v.ValidatorServiceOrgGroupSlice, ou.OrgID)
	if err != nil {
		log.Ctx(ctx).Err(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, v.ValidatorServiceOrgGroupSlice)
}
