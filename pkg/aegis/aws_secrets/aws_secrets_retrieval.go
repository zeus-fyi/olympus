package aws_secrets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
)

func RetrieveSecretValue(ctx context.Context, ou org_users.OrgUser, ref string) (*SecretsRequest, error) {
	if len(ref) <= 0 {
		return nil, errors.New("ref is required")
	}
	sv, err := artemis_hydra_orchestrations_aws_auth.GetOrgSecret(ctx, FormatSecret(ou.OrgID))
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg(fmt.Sprintf("%s", err.Error()))
		return nil, err
	}
	m := make(map[string]SecretsKeyValue)
	err = json.Unmarshal(sv, &m)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg(fmt.Sprintf("%s", err.Error()))
		return nil, err
	}
	for k, v := range m {
		if k == ref {
			return &SecretsRequest{
				Name:  ref,
				Key:   v.Key,
				Value: v.Value,
			}, nil
		}
	}
	return nil, errors.New("ref not found")
}

func ReadSecretReferences(ctx context.Context, ou org_users.OrgUser) ([]SecretsRequest, error) {
	sv, err := artemis_hydra_orchestrations_aws_auth.GetOrgSecret(ctx, FormatSecret(ou.OrgID))
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg(fmt.Sprintf("%s", err.Error()))
		return nil, err
	}
	m := make(map[string]SecretsKeyValue)
	err = json.Unmarshal(sv, &m)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg(fmt.Sprintf("%s", err.Error()))
		return nil, err
	}
	var mr []SecretsRequest
	for name, v := range m {
		mr = append(mr, SecretsRequest{
			Name: name,
			Key:  v.Key,
		})
	}
	return mr, nil
}
