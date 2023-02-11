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
	hestia_aws_secrets_auth "github.com/zeus-fyi/olympus/hestia/auth"
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
		log.Ctx(ctx).Error().Err(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	la := v.ServiceAuth
	b, err := json.Marshal(la)
	if err != nil {
		log.Ctx(ctx).Error().Err(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	name := fmt.Sprintf("%s-%d-%s", v.GroupName, ou.OrgID, network)
	si := secretsmanager.CreateSecretInput{
		Name:         aws.String(name),
		Description:  aws.String(name),
		SecretBinary: b,
		SecretString: nil,
	}
	err = hestia_aws_secrets_auth.HestiaSecretManagerAuthAWS.CreateNewSecret(ctx, si)
	if err != nil {
		errCheckStr := fmt.Sprintf("the secret %s already exists", name)
		if strings.Contains(err.Error(), errCheckStr) {
			fmt.Println("Secret already exists, skipping")
		} else {
			log.Ctx(ctx).Error().Err(err)
			return c.JSON(http.StatusInternalServerError, nil)
		}
	}

	switch v.ProtocolNetworkID {
	case hestia_req_types.EthereumMainnetProtocolNetworkID:
		err = eth_validators_service_requests.ArtemisEthereumMainnetValidatorsRequestsWorker.ExecuteServiceNewValidatorsToCloudCtxNsWorkflow(ctx, vsr)
		if err != nil {
			log.Ctx(ctx).Error().Err(err)
			return c.JSON(http.StatusAccepted, nil)
		}
	case hestia_req_types.EthereumEphemeryProtocolNetworkID:
		err = eth_validators_service_requests.ArtemisEthereumEphemeryValidatorsRequestsWorker.ExecuteServiceNewValidatorsToCloudCtxNsWorkflow(ctx, vsr)
		if err != nil {
			log.Ctx(ctx).Error().Err(err)
			return c.JSON(http.StatusAccepted, nil)
		}
	default:
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusBadRequest, nil)
}
