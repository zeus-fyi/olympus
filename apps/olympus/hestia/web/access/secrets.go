package hestia_access_keygen

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
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	"golang.org/x/crypto/sha3"
)

func SecretsRequestHandler(c echo.Context) error {
	request := new(SecretsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateOrUpdateSecret(c, false)
}

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

func FormatSecret(orgID int) string {
	hash := sha3.New256()
	_, _ = hash.Write([]byte(fmt.Sprintf("org-%d-%s", orgID, "hestia")))
	// Get the resulting encoded byte slice
	sha3v := hash.Sum(nil)
	return fmt.Sprintf("%x", hash.Sum(sha3v))
}

type SecretsKeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
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
	m := make(map[string]SecretsKeyValue)
	err := a.Validate(isDelete)
	if err != nil {
		log.Info().Msgf("Hestia: CreateValidatorServiceRequest: Unexpected Error: %s", err.Error())
		return c.JSON(http.StatusBadRequest, nil)
	}
	ctx := context.Background()
	exists := artemis_hydra_orchestrations_aws_auth.HydraSecretManagerAuthAWS.DoesSecretExist(ctx, FormatSecret(ou.OrgID))
	if exists {
		sv, serr := artemis_hydra_orchestrations_aws_auth.GetOrgSecret(ctx, FormatSecret(ou.OrgID))
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
	m[a.Name] = SecretsKeyValue{
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
	orgSecretName := FormatSecret(ou.OrgID)
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

func SecretReadRequestHandler(c echo.Context) error {
	return RetrieveSecretValue(c)
}

func RetrieveSecretValue(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Warn().Interface("ou", ou).Msg("CreateOrUpdateSecret: orgUser not found")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	ctx := context.Background()
	ref := c.Param("ref")
	if len(ref) <= 0 {
		return c.JSON(http.StatusBadRequest, "ref is required")
	}
	sv, err := artemis_hydra_orchestrations_aws_auth.GetOrgSecret(ctx, FormatSecret(ou.OrgID))
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg(fmt.Sprintf("%s", err.Error()))
		return c.JSON(http.StatusInternalServerError, nil)
	}
	m := make(map[string]SecretsKeyValue)
	err = json.Unmarshal(sv, &m)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg(fmt.Sprintf("%s", err.Error()))
		return c.JSON(http.StatusInternalServerError, nil)
	}
	for k, v := range m {
		if k == ref {
			return c.JSON(http.StatusOK, SecretsRequest{
				Name:  ref,
				Key:   v.Key,
				Value: v.Value,
			})
		}
	}
	return c.JSON(http.StatusNotFound, nil)
}

func SecretsReadRequestHandler(c echo.Context) error {
	request := new(SecretsRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.ReadSecretReferences(c)
}

func (a *SecretsRequest) ReadSecretReferences(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Warn().Interface("ou", ou).Msg("ReadSecretReferences: orgUser not found")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	ctx := context.Background()
	sv, err := artemis_hydra_orchestrations_aws_auth.GetOrgSecret(ctx, FormatSecret(ou.OrgID))
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg(fmt.Sprintf("%s", err.Error()))
		return c.JSON(http.StatusInternalServerError, nil)
	}
	m := make(map[string]SecretsKeyValue)
	err = json.Unmarshal(sv, &m)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg(fmt.Sprintf("%s", err.Error()))
		return c.JSON(http.StatusInternalServerError, nil)
	}
	var mr []SecretsRequest
	for name, v := range m {
		mr = append(mr, SecretsRequest{
			Name: name,
			Key:  v.Key,
		})
	}
	return c.JSON(http.StatusOK, mr)
}

func SecretDeleteRequestHandler(c echo.Context) error {
	return DeleteSecretValue(c)
}

func DeleteSecretValue(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Warn().Interface("ou", ou).Msg("ReadSecretReferences: orgUser not found")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	ref := c.Param("ref")
	if len(ref) <= 0 {
		return c.JSON(http.StatusBadRequest, "ref is required")
	}
	sr := SecretsRequest{
		Name: ref,
	}
	return sr.CreateOrUpdateSecret(c, true)
}