package eth_validator_signature_requests

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	artemis_validator_signature_service_routing "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/signature_routing"
	bls_serverless_signing "github.com/zeus-fyi/zeus/pkg/aegis/aws/serverless_signing"
	"time"

	aegis_inmemdbs "github.com/zeus-fyi/zeus/pkg/aegis/inmemdbs"
)

const (
	waitForTxRxTimeout    = 15 * time.Minute
	submitSignedTxTimeout = 5 * time.Minute
)

type ArtemisEthereumValidatorSignatureRequestActivities struct {
}

func NewArtemisEthereumValidatorSignatureRequestActivities() ArtemisEthereumValidatorSignatureRequestActivities {
	return ArtemisEthereumValidatorSignatureRequestActivities{}
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

type Eth2SignResponse struct {
	Signature string `json:"signature"`
}

func (d *ArtemisEthereumValidatorSignatureRequestActivities) GetActivities() ActivitiesSlice {
	return []interface{}{d.RequestValidatorSignatures}
}

type Resty struct {
	*resty.Client
}

// TODO add signature auth, todo make each group an async activity

func (d *ArtemisEthereumValidatorSignatureRequestActivities) RequestValidatorSignatures(ctx context.Context, sigRequests aegis_inmemdbs.EthereumBLSKeySignatureRequests) (aegis_inmemdbs.EthereumBLSKeySignatureResponses, error) {
	sigResponses := aegis_inmemdbs.EthereumBLSKeySignatureResponses{}
	gm := artemis_validator_signature_service_routing.GroupSigRequestsByGroupName(ctx, sigRequests)
	r := Resty{}
	r.Client = resty.New()
	for groupName, signReqs := range gm {
		auth, err := artemis_validator_signature_service_routing.GetGroupAuthFromInMemFS(ctx, groupName)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to get group auth")
			return sigResponses, err
		}
		sr := bls_serverless_signing.SignatureRequests{
			SecretName:        auth.SecretName,
			SignatureRequests: aegis_inmemdbs.EthereumBLSKeySignatureRequests{Map: signReqs.Map},
		}
		respJson := aegis_inmemdbs.EthereumBLSKeySignatureResponses{Map: make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureResponse)}
		_, err = r.R().
			SetResult(&respJson).
			SetBody(sr).
			Post(auth.AuthLamdbaAWS.ServiceURL)
		if err != nil {
			log.Ctx(ctx).Error().Err(err)
			return sigResponses, err
		}
		for k, v := range respJson.Map {
			sigResponses.Map[k] = v
		}
	}
	return sigResponses, nil
}
