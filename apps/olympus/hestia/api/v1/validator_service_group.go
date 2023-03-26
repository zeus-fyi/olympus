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
	create_org_users "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	eth_validators_service_requests "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validators_service_requests"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
)

type CreateValidatorServiceRequest struct {
	hestia_req_types.ServiceRequestWrapper         `json:"serviceRequestWrapper"`
	hestia_req_types.ValidatorServiceOrgGroupSlice `json:"validatorServiceOrgGroupSlice"`
}

func CreateValidatorServiceRequestHandler(c echo.Context) error {
	log.Info().Msg("Hestia: CreateValidatorServiceRequestHandler")
	request := new(CreateValidatorServiceRequest)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("CreateValidatorServiceRequestHandler")
		return err
	}
	return request.CreateValidatorsServiceGroup(c)
}

func (v *CreateValidatorServiceRequest) CreateValidatorsServiceGroup(c echo.Context) error {
	log.Info().Msg("Hestia: CreateValidatorServiceRequest: CreateValidatorsServiceGroup")
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)

	log.Ctx(ctx).Info().Interface("ou", ou).Interface("vsg", v.ValidatorServiceOrgGroupSlice).Msg("CreateValidatorsServiceGroup")
	vsr := eth_validators_service_requests.ValidatorServiceGroupWorkflowRequest{
		OrgID:                         ou.OrgID,
		ServiceRequestWrapper:         v.ServiceRequestWrapper,
		ValidatorServiceOrgGroupSlice: v.ValidatorServiceOrgGroupSlice,
	}
	for i, key := range v.ValidatorServiceOrgGroupSlice {
		v.ValidatorServiceOrgGroupSlice[i].Pubkey = strings_filter.AddHexPrefix(key.Pubkey)
		v.ValidatorServiceOrgGroupSlice[i].FeeRecipient = strings_filter.AddHexPrefix(key.FeeRecipient)
	}

	var network string
	bearer := c.Get("bearer").(string)
	switch v.ProtocolNetworkID {
	case hestia_req_types.EthereumMainnetProtocolNetworkID:
		key, err := auth.VerifyBearerTokenService(ctx, bearer, create_org_users.EthereumMainnetService)
		if err != nil || key.PublicKeyVerified == false {
			log.Err(err).Interface("orgUser", ou).Msg("CreateValidatorsServiceGroup: EthereumMainnetService unauthorized")
			return c.JSON(http.StatusUnauthorized, nil)
		}
		network = hestia_req_types.ProtocolNetworkIDToString(v.ProtocolNetworkID)
	case hestia_req_types.EthereumGoerliProtocolNetworkID:
		key, err := auth.VerifyBearerTokenService(ctx, bearer, create_org_users.EthereumGoerliService)
		if err != nil || key.PublicKeyVerified == false {
			log.Err(err).Interface("orgUser", ou).Msg("CreateValidatorsServiceGroup: EthereumGoerliService unauthorized")
			return c.JSON(http.StatusUnauthorized, nil)
		}
		network = hestia_req_types.ProtocolNetworkIDToString(v.ProtocolNetworkID)
	case hestia_req_types.EthereumEphemeryProtocolNetworkID:
		key, err := auth.VerifyBearerTokenService(ctx, bearer, create_org_users.EthereumEphemeryService)
		if err != nil || key.PublicKeyVerified == false {
			log.Err(err).Interface("orgUser", ou).Msg("CreateValidatorsServiceGroup: EthereumEphemeryService unauthorized")
			return c.JSON(http.StatusUnauthorized, nil)
		}
		network = hestia_req_types.ProtocolNetworkIDToString(v.ProtocolNetworkID)
	default:
		err := errors.New("unknown network")
		log.Ctx(ctx).Err(err).Msg("CreateValidatorServiceRequest")
		return c.JSON(http.StatusBadRequest, err)
	}

	log.Info().Msg("Hestia: CreateValidatorServiceRequest: Validating Service Auth")
	err := v.ServiceAuth.Validate()
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("service auth failed validation")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	b, err := json.Marshal(v.ServiceRequestWrapper)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("service auth failed json marshal")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	name := fmt.Sprintf("%s-%d-%s", v.GroupName, ou.OrgID, network)
	si := secretsmanager.CreateSecretInput{
		Name:         aws.String(name),
		Description:  aws.String(name),
		SecretBinary: b,
	}
	log.Info().Msg("Hestia: CreateValidatorServiceRequest: Service Auth Valid, Creating Secret")
	err = artemis_hydra_orchestrations_aws_auth.HydraSecretManagerAuthAWS.CreateNewSecret(ctx, si)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("%s", err.Error()))
		errCheckStr := fmt.Sprintf("the secret %s already exists", name)
		if strings.Contains(err.Error(), errCheckStr) {
			log.Err(err).Msg("secret already exists, updating to new values")
			su := &secretsmanager.UpdateSecretInput{
				SecretId:     aws.String(name),
				SecretBinary: b,
			}
			_, err = artemis_hydra_orchestrations_aws_auth.HydraSecretManagerAuthAWS.UpdateSecret(ctx, su)
			if err != nil {
				log.Ctx(ctx).Err(err).Msg("service auth failed to update secret")
				log.Info().Msgf("Hestia: CreateValidatorServiceRequest: Unexpected Error: %s", err.Error())
				return c.JSON(http.StatusInternalServerError, nil)
			}
		} else {
			log.Info().Msgf("Hestia: CreateValidatorServiceRequest: Unexpected Error: %s", err.Error())
			log.Ctx(ctx).Err(err).Msg("service auth failed to create secret")
			return c.JSON(http.StatusInternalServerError, nil)
		}
	}
	log.Info().Msg("Hestia: CreateValidatorsServiceGroup Secret Created: Init Validator Service Workflow")
	// clear auth, not needed anymore, and we don't want to log it in temporal
	vsr.ServiceAuth = hestia_req_types.ServiceAuthConfig{}
	resp := Response{}
	switch v.ProtocolNetworkID {
	case hestia_req_types.EthereumMainnetProtocolNetworkID:
		err = eth_validators_service_requests.ArtemisEthereumMainnetValidatorsRequestsWorker.ExecuteServiceNewValidatorsToCloudCtxNsWorkflow(ctx, vsr)
		if err != nil {
			log.Ctx(ctx).Err(err).Interface("network", network).Msg("ExecuteServiceNewValidatorsToCloudCtxNsWorkflow")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		resp.Message = "Ethereum Mainnet validators service request in progress"
		return c.JSON(http.StatusAccepted, resp)
	case hestia_req_types.EthereumGoerliProtocolNetworkID:
		err = eth_validators_service_requests.ArtemisEthereumGoerliValidatorsRequestsWorker.ExecuteServiceNewValidatorsToCloudCtxNsWorkflow(ctx, vsr)
		if err != nil {
			log.Ctx(ctx).Err(err).Interface("network", network).Msg("ExecuteServiceNewValidatorsToCloudCtxNsWorkflow")
			return c.JSON(http.StatusInternalServerError, nil)
		}
		resp.Message = "Ethereum Mainnet validators service request in progress"
		return c.JSON(http.StatusAccepted, resp)
	case hestia_req_types.EthereumEphemeryProtocolNetworkID:
		err = eth_validators_service_requests.ArtemisEthereumEphemeryValidatorsRequestsWorker.ExecuteServiceNewValidatorsToCloudCtxNsWorkflow(ctx, vsr)
		if err != nil {
			log.Ctx(ctx).Err(err).Interface("network", network).Msg("ExecuteServiceNewValidatorsToCloudCtxNsWorkflow")
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
