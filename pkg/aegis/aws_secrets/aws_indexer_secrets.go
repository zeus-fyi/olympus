package aws_secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
)

func GetTelegramToken(ctx context.Context, orgID int) (string, error) {
	sv, err := artemis_hydra_orchestrations_aws_auth.GetOrgSecret(ctx, FormatSecret(orgID))
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("%s", err.Error()))
		return "", err
	}
	m := make(map[string]SecretsKeyValue)
	err = json.Unmarshal(sv, &m)
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("%s", err.Error()))
		return "", err
	}

	token := ""
	tv, ok := SecretsCache.Get("telegram_token")
	if ok {
		token = tv.(string)
	}
	if len(token) == 0 {
		for k, v := range m {
			if k == "telegram" {
				if v.Key == "token" {
					SecretsCache.Set("telegram_token", v.Value, cache.DefaultExpiration)
					token = v.Value
				}
			}
		}
	}
	return token, err
}

var (
	SecretCache = cache.New(time.Hour*24*3, cache.DefaultExpiration)
)

func ClearOrgSecretCache(ou org_users.OrgUser) {
	SecretCache.Delete(FormatSecret(ou.OrgID))
}

func GetMockingbirdPlatformSecrets(ctx context.Context, ou org_users.OrgUser, platform string) (*OAuth2PlatformSecret, error) {
	m := make(map[string]SecretsKeyValue)
	svCached, ok := SecretCache.Get(FormatSecret(ou.OrgID))
	if ok {
		skv, cacheOk := svCached.(map[string]SecretsKeyValue)
		if cacheOk {
			m = skv
		}
	}
	if len(m) == 0 {
		sv, err := artemis_hydra_orchestrations_aws_auth.GetOrgSecret(ctx, FormatSecret(ou.OrgID))
		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("%s", err.Error()))
			return nil, err
		}
		err = json.Unmarshal(sv, &m)
		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("%s", err.Error()))
			return nil, err
		}
		SecretCache.Set(FormatSecret(ou.OrgID), m, cache.DefaultExpiration)
	}

	mp := MockingBirdPlatformNames(platform)
	op := &OAuth2PlatformSecret{
		Platform: platform,
	}

	for mkeyName, mockingbird := range mp {
		svItem, sok := m[mkeyName]
		if sok && svItem.Key == mockingbird {
			if mkeyName == fmt.Sprintf("%s-oauth2-public", platform) {
				op.OAuth2Public = svItem.Value
			}
			if mkeyName == fmt.Sprintf("%s-oauth2-secret", platform) {
				op.OAuth2Secret = svItem.Value
			}
			if mkeyName == fmt.Sprintf("%s-username", platform) {
				op.Username = svItem.Value
			}
			if mkeyName == fmt.Sprintf("%s-password", platform) {
				op.Password = svItem.Value
			}
			if mkeyName == fmt.Sprintf("%s-api-key", platform) {
				op.ApiKey = svItem.Value
			}
		}
	}
	return op, nil
}

func MockingBirdPlatformNames(platform string) map[string]string {
	return map[string]string{
		fmt.Sprintf("%s-oauth2-public", platform): "mockingbird",
		fmt.Sprintf("%s-oauth2-secret", platform): "mockingbird",
		fmt.Sprintf("%s-username", platform):      "mockingbird",
		fmt.Sprintf("%s-password", platform):      "mockingbird",
		fmt.Sprintf("%s-api-key", platform):       "mockingbird",
	}
}

type OAuth2PlatformSecret struct {
	Platform     string `json:"platform"`
	OAuth2Public string `json:"oauth2Public"`
	OAuth2Secret string `json:"oauth2Secret"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	ApiKey       string `json:"apiKey,omitempty"`
}