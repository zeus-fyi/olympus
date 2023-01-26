package v1hestia

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	eth_validators_service_requests "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validators_service_requests"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

type CreateValidatorServiceRequest struct {
	hestia_req_types.ServiceRequestWrapper
	hestia_req_types.ValidatorServiceOrgGroupSlice
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
	log.Ctx(ctx).Info().Interface("ou", ou).Interface("vsg", v.ValidatorServiceOrgGroupSlice).Msg("CreateValidatorsServiceGroup")
	vsr := eth_validators_service_requests.ValidatorServiceGroupWorkflowRequest{
		ServiceRequestWrapper:         v.ServiceRequestWrapper,
		ValidatorServiceOrgGroupSlice: v.ValidatorServiceOrgGroupSlice,
	}
	var err error
	switch v.ProtocolNetworkID {
	case hestia_req_types.EthereumMainnetProtocolNetworkID:
		err = eth_validators_service_requests.ArtemisEthereumMainnetValidatorsRequestsWorker.ExecuteServiceNewValidatorsToCloudCtxNsWorkflow(ctx, vsr)
	case hestia_req_types.EthereumEphemeryProtocolNetworkID:
		err = eth_validators_service_requests.ArtemisEthereumEphemeryValidatorsRequestsWorker.ExecuteServiceNewValidatorsToCloudCtxNsWorkflow(ctx, vsr)
	default:
		return c.JSON(http.StatusBadRequest, nil)
	}
	if err != nil {
		log.Ctx(ctx).Err(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusAccepted, nil)
}
