package hestia_login

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
)

type SecretsRequest struct {
	Name  string `json:"name"`
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
}

func (a *SecretsRequest) Validate(isDelete bool) error {
	if len(a.Name) <= 0 {
		return fmt.Errorf("name is required")
	}
	if isDelete {
		return nil
	}
	if len(a.Key) <= 0 {
		return fmt.Errorf("key is required")
	}
	if len(a.Value) <= 0 {
		return fmt.Errorf("value is required")
	}
	return nil
}
func (a *SecretsRequest) CreateOrUpdateSecret(c echo.Context, isDelete bool) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Warn().Interface("ou", ou).Msg("CreateOrUpdateSecret: orgUser not found")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if ou.OrgID <= 0 {
		log.Warn().Interface("ou", ou).Msg("CreateOrUpdateSecret: orgID not found")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	m := make(map[string]aws_secrets.SecretsKeyValue)
	err := a.Validate(isDelete)
	if err != nil {
		log.Info().Msgf("Hestia: CreateValidatorServiceRequest: Unexpected Error: %s", err.Error())
		return c.JSON(http.StatusBadRequest, nil)
	}
	a.Name = strings.TrimSpace(a.Name)
	a.Key = strings.TrimSpace(a.Key)
	a.Value = strings.TrimSpace(a.Value)

	ctx := context.Background()
	exists := artemis_hydra_orchestrations_aws_auth.HydraSecretManagerAuthAWS.DoesSecretExist(ctx, aws_secrets.FormatSecret(ou.OrgID))
	if exists {
		sv, serr := artemis_hydra_orchestrations_aws_auth.GetOrgSecret(ctx, aws_secrets.FormatSecret(ou.OrgID))
		if serr != nil {
			log.Err(serr).Interface("ou", ou).Msg(fmt.Sprintf("%s", serr.Error()))
			return c.JSON(http.StatusInternalServerError, nil)
		}

		err = json.Unmarshal(sv, &m)
		if err != nil {
			log.Err(err).Interface("ou", ou).Msg(fmt.Sprintf("%s", err.Error()))
			return c.JSON(http.StatusInternalServerError, nil)
		}
		for k, v := range m {
			if k == a.Name {
				if v.Key == a.Key && v.Value == a.Value {
					log.Info().Msg("Hestia Auth Valid, No Update Needed")
					return c.JSON(http.StatusOK, nil)
				}
			}
		}
	}
	log.Info().Msg("Hestia: Secret")
	m[a.Name] = aws_secrets.SecretsKeyValue{
		Key:   a.Key,
		Value: a.Value,
	}
	if isDelete {
		delete(m, a.Name)
	}
	b, err := json.Marshal(m)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg(fmt.Sprintf("%s", err.Error()))
		return c.JSON(http.StatusInternalServerError, nil)
	}
	orgSecretName := aws_secrets.FormatSecret(ou.OrgID)
	si := secretsmanager.CreateSecretInput{
		Name:         aws.String(orgSecretName),
		Description:  aws.String(orgSecretName),
		SecretBinary: b,
	}

	err = artemis_hydra_orchestrations_aws_auth.HydraSecretManagerAuthAWS.CreateNewSecret(ctx, si)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("%s", err.Error()))
		errCheckStr := fmt.Sprintf("the secret %s already exists", orgSecretName)
		if strings.Contains(err.Error(), errCheckStr) {
			log.Err(err).Msg("secret already exists, updating to new values")
			su := &secretsmanager.UpdateSecretInput{
				SecretId:     aws.String(orgSecretName),
				SecretBinary: b,
			}
			_, err = artemis_hydra_orchestrations_aws_auth.HydraSecretManagerAuthAWS.UpdateSecret(ctx, su)
			if err != nil {
				log.Err(err).Msg("service auth failed to update secret")
				log.Info().Msgf("Hestia: CreateOrUpdateSecret: Unexpected Error: %s", err.Error())
				return c.JSON(http.StatusInternalServerError, nil)
			}
		} else {
			log.Info().Msgf("Hestia: CreateOrUpdateSecret: Unexpected Error: %s", err.Error())
			log.Err(err).Msg("service auth failed to create secret")
			return c.JSON(http.StatusInternalServerError, nil)
		}
	}
	return c.JSON(http.StatusOK, nil)
}
