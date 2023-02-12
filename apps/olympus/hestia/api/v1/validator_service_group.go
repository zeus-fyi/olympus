package v1hestia

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
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
		log.Error().Err(err).Msg("CreateValidatorServiceRequestHandler")
		return err
	}
	return request.CreateValidatorsServiceGroup(c)
}

func (v *CreateValidatorServiceRequest) CreateValidatorsServiceGroup(c echo.Context) error {
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	log.Ctx(ctx).Info().Interface("ou", ou).Interface("vsg", v.ValidatorServiceOrgGroupSlice).Msg("CreateValidatorsServiceGroup")
	vsr := eth_validators_service_requests.ValidatorServiceGroupWorkflowRequest{
		OrgID:                         ou.OrgID,
		ServiceRequestWrapper:         v.ServiceRequestWrapper,
		ValidatorServiceOrgGroupSlice: v.ValidatorServiceOrgGroupSlice,
	}

	var network string
	switch v.ProtocolNetworkID {
	case hestia_req_types.EthereumMainnetProtocolNetworkID:
		network = hestia_req_types.ProtocolNetworkIDToString(v.ProtocolNetworkID)
	case hestia_req_types.EthereumEphemeryProtocolNetworkID:
		network = hestia_req_types.ProtocolNetworkIDToString(v.ProtocolNetworkID)
	default:
		return c.JSON(http.StatusBadRequest, errors.New("unknown network"))

	}
	err := v.ServiceAuth.Validate()
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("service auth failed validation")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	la := v.ServiceAuth
	b, err := json.Marshal(la)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("service auth failed json marshal")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	name := fmt.Sprintf("%s-%d-%s", v.GroupName, ou.OrgID, network)
	si := secretsmanager.CreateSecretInput{
		Name:         aws.String(name),
		Description:  aws.String(name),
		SecretBinary: b,
		SecretString: nil,
	}
	err = artemis_hydra_orchestrations_aws_auth.HydraSecretManagerAuthAWS.CreateNewSecret(ctx, si)
	if err != nil {
		errCheckStr := fmt.Sprintf("the secret %s already exists", name)
		if strings.Contains(err.Error(), errCheckStr) {
			fmt.Println("Secret already exists, updating to new values")
		} else {
			log.Ctx(ctx).Error().Err(err)
			return c.JSON(http.StatusInternalServerError, nil)
		}
	}
	// clear auth, not needed anymore, and we don't want to log it in temporal
	vsr.ServiceAuth = hestia_req_types.ServiceAuthConfig{}
	resp := Response{}
	switch v.ProtocolNetworkID {
	case hestia_req_types.EthereumMainnetProtocolNetworkID:
		err = eth_validators_service_requests.ArtemisEthereumMainnetValidatorsRequestsWorker.ExecuteServiceNewValidatorsToCloudCtxNsWorkflow(ctx, vsr)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Interface("network", network).Msg("ExecuteServiceNewValidatorsToCloudCtxNsWorkflow")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		resp.Message = "Ethereum Mainnet validators service request in progress"
		return c.JSON(http.StatusAccepted, resp)
	case hestia_req_types.EthereumEphemeryProtocolNetworkID:
		err = eth_validators_service_requests.ArtemisEthereumEphemeryValidatorsRequestsWorker.ExecuteServiceNewValidatorsToCloudCtxNsWorkflow(ctx, vsr)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Interface("network", network).Msg("ExecuteServiceNewValidatorsToCloudCtxNsWorkflow")
			return c.JSON(http.StatusInternalServerError, resp)
		}
		resp.Message = "Ethereum Ephemery validators service request in progress"
		return c.JSON(http.StatusAccepted, resp)
	default:
		return c.JSON(http.StatusBadRequest, nil)
	}
}

type Response struct {
	Message string `json:"message"`
}
