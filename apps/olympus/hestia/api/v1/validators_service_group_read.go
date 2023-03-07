package v1hestia

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

type GetValidatorServiceInfo struct {
	hestia_req_types.ServiceRequestWrapper
	hestia_req_types.ValidatorServiceOrgGroupSlice
}

func GetValidatorServiceInfoHandler(c echo.Context) error {
	log.Info().Msg("Hestia: CreateValidatorServiceRequestHandler")
	request := new(GetValidatorServiceInfo)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("CreateValidatorServiceRequestHandler")
		return err
	}
	return request.GetValidatorServiceInfo(c)
}

func (v *GetValidatorServiceInfo) GetValidatorServiceInfo(c echo.Context) error {
	log.Info().Msg("Hestia: GetValidatorServiceInfo")
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	network := c.Param("network")
	pid := hestia_req_types.ProtocolNetworkStringToID(network)
	vs, err := artemis_validator_service_groups_models.SelectValidatorsServiceInfo(ctx, pid, ou.OrgID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	log.Ctx(ctx).Info().Interface("ou", ou).Msg("GetValidatorServiceInfo")
	return c.JSON(http.StatusOK, vs)
}
