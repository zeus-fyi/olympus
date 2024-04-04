package aws_secrets

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
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
	SecretCache = cache.New(time.Hour*1, cache.DefaultExpiration)
)

func ClearOrgSecretCache(ou org_users.OrgUser) {
	SecretCache.Delete(FormatSecret(ou.OrgID))
}

var (
	CredBasePath = "/.aws/credentials"
	ConfigPath   = "/.aws/config"
)

func GetServiceAccountSecrets(ctx context.Context, ou org_users.OrgUser) (ServiceAccountPlatformSecrets, error) {
	sps := ServiceAccountPlatformSecrets{
		AwsEksServiceMap: make(map[string]aegis_aws_auth.AuthAWS),
	}
	m := make(map[string]SecretsKeyValue)
	sv, err := artemis_hydra_orchestrations_aws_auth.GetOrgSecret(ctx, FormatSecret(ou.OrgID))
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("%s", err.Error()))
		return sps, err
	}
	err = json.Unmarshal(sv, &m)
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("%s", err.Error()))
		return sps, err
	}
	for secretName, v := range m {
		if strings.HasPrefix(secretName, "zeus-aws-eks-") {
			secretNameWithoutPrefixSuffix := strings.TrimPrefix(secretName, "zeus-aws-eks-")
			if _, ok := sps.AwsEksServiceMap[v.Key]; !ok {
				sps.AwsEksServiceMap[v.Key] = aegis_aws_auth.AuthAWS{}
			}
			tmp := sps.AwsEksServiceMap[v.Key]
			if strings.HasSuffix(secretName, "-service-account-access-key") {
				tmp.AccessKey = v.Value
				secretNameWithoutPrefixSuffix = strings.TrimSuffix(secretNameWithoutPrefixSuffix, "-service-account-access-key")
			}
			if strings.HasSuffix(secretName, "-service-account-secret-key") {
				tmp.SecretKey = v.Value
				secretNameWithoutPrefixSuffix = strings.TrimSuffix(secretNameWithoutPrefixSuffix, "-service-account-secret-key")
			}
			tmp.Region = secretNameWithoutPrefixSuffix
			sps.AwsEksServiceMap[v.Key] = tmp
		}
	}
	for _, v := range sps.AwsEksServiceMap {
		err = AddOrUpdateProfile(CredBasePath, fmt.Sprintf("%d", ou.OrgID), v.AccessKey, v.SecretKey)
		if err != nil {
			return sps, err
		}
		err = AddOrUpdateConfig(ConfigPath, fmt.Sprintf("%d", ou.OrgID), v.Region)
		if err != nil {
			return sps, err
		}
	}

	return sps, err
}
func AddOrUpdateConfig(filePath, profileName, region string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("failed to open or create file: %v", err)
	}
	defer file.Close()

	// Variables to hold the new contents of the file and a flag to mark if the profile was found.
	var newContents []string
	profileFound := false
	profileHeader := fmt.Sprintf("[profile %s]", profileName)
	scanner := bufio.NewScanner(file)

	// Scan through the file line by line.
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the current line marks the beginning of the profile.
		if strings.TrimSpace(line) == profileHeader {
			profileFound = true
			newContents = append(newContents, line)
			newContents = append(newContents, fmt.Sprintf("region = %s", region))
			// Skip the next two lines assuming they are the keys we want to update.
			scanner.Scan() // Skip aws_access_key_id
			scanner.Scan() // Skip aws_secret_access_key
			continue
		}
		newContents = append(newContents, line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// If the profile was not found, add it to the end of the file.
	if !profileFound {
		newContents = append(newContents, profileHeader)
		newContents = append(newContents, fmt.Sprintf("region = %s", region))
	}

	// Write the new contents back to the file.
	return writeToFile(filePath, newContents)
}

func AddOrUpdateProfile(filePath, profileName, awsAccessKeyID, awsSecretAccessKey string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("failed to open or create file: %v", err)
	}
	defer file.Close()

	// Variables to hold the new contents of the file and a flag to mark if the profile was found.
	var newContents []string
	profileFound := false
	profileHeader := fmt.Sprintf("[%s]", profileName)
	scanner := bufio.NewScanner(file)

	// Scan through the file line by line.
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the current line marks the beginning of the profile.
		if strings.TrimSpace(line) == profileHeader {
			profileFound = true
			newContents = append(newContents, line)
			newContents = append(newContents, fmt.Sprintf("aws_access_key_id = %s", awsAccessKeyID))
			newContents = append(newContents, fmt.Sprintf("aws_secret_access_key = %s", awsSecretAccessKey))
			// Skip the next two lines assuming they are the keys we want to update.
			scanner.Scan() // Skip aws_access_key_id
			scanner.Scan() // Skip aws_secret_access_key
			continue
		}
		newContents = append(newContents, line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// If the profile was not found, add it to the end of the file.
	if !profileFound {
		newContents = append(newContents, profileHeader)
		newContents = append(newContents, fmt.Sprintf("aws_access_key_id = %s", awsAccessKeyID))
		newContents = append(newContents, fmt.Sprintf("aws_secret_access_key = %s", awsSecretAccessKey))
	}

	// Write the new contents back to the file.
	return writeToFile(filePath, newContents)
}

// Helper function to write the new contents to the file.
func writeToFile(filePath string, contents []string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range contents {
		_, err = writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}
	}

	return writer.Flush()
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
	mpAdd := MockingBirdPlatformNamesCloudServices("mb")
	for k, v := range mpAdd {
		mp[k] = v
	}
	for mkeyName, mockingbird := range mp {
		//fmt.Println(mkeyName, "mkey", mockingbird, "mockingbird")
		svItem, sok := m[mkeyName]
		if sok && svItem.Key == mockingbird {
			if strings.Contains(mkeyName, "twillio") {
				ss := strings.Split(svItem.Value, ":")
				if len(ss) == 2 {
					op.TwillioAccount = ss[0]
					op.TwillioAuth = ss[1]
				}
			}
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
			if strings.HasSuffix(mkeyName, "client-id") {
				op.ClientID = svItem.Value
			}
			if strings.HasSuffix(mkeyName, "consumer-secret") {
				op.ConsumerSecret = svItem.Value
			}
			if strings.HasSuffix(mkeyName, "consumer-key") {
				op.ConsumerPublic = svItem.Value
			}
			if strings.HasSuffix(mkeyName, "access-token-secret") {
				op.AccessTokenSecret = svItem.Value
			}
			if strings.HasSuffix(mkeyName, "access-token-public") {
				op.AccessTokenPublic = svItem.Value
			}
			if strings.HasSuffix(mkeyName, platform) {
				op.BearerToken = svItem.Value
			}
			if strings.HasSuffix(mkeyName, "-s3-access") {
				op.S3AccessKey = svItem.Value
			}
			if strings.HasSuffix(mkeyName, "-s3-secret") {
				op.S3SecretKey = svItem.Value
			}
		}
	}
	return op, nil
}

func MockingBirdPlatformNames(platform string) map[string]string {
	return map[string]string{
		fmt.Sprintf("%s-oauth2-public", platform):       "mockingbird",
		fmt.Sprintf("%s-oauth2-secret", platform):       "mockingbird",
		fmt.Sprintf("%s-username", platform):            "mockingbird",
		fmt.Sprintf("%s-password", platform):            "mockingbird",
		fmt.Sprintf("%s-api-key", platform):             "mockingbird",
		fmt.Sprintf("%s-consumer-key", platform):        "mockingbird",
		fmt.Sprintf("%s-consumer-secret", platform):     "mockingbird",
		fmt.Sprintf("%s-client-id", platform):           "mockingbird",
		fmt.Sprintf("%s-access-token-secret", platform): "mockingbird",
		fmt.Sprintf("%s-access-token-public", platform): "mockingbird",
		fmt.Sprintf("%s", platform):                     "mockingbird",
	}
}

func MockingBirdPlatformNamesCloudServices(platform string) map[string]string {
	return map[string]string{
		fmt.Sprintf("%s-s3-access", platform): "mockingbird-s3-ovh-us-west-or",
		fmt.Sprintf("%s-s3-secret", platform): "mockingbird-s3-ovh-us-west-or",
	}
}

type OAuth2PlatformSecret struct {
	S3AccessKey       string `json:"s3AccessKey"`
	S3SecretKey       string `json:"s3SecretKey"`
	TwillioAccount    string `json:"twillioAccount"`
	TwillioAuth       string `json:"twillioAuth"`
	ConsumerPublic    string `json:"consumerPublic"`
	ConsumerSecret    string `json:"consumerSecret"`
	Platform          string `json:"platform"`
	ClientID          string `json:"clientID"`
	OAuth2Public      string `json:"oauth2Public"`
	OAuth2Secret      string `json:"oauth2Secret"`
	Username          string `json:"username,omitempty"`
	Password          string `json:"password,omitempty"`
	ApiKey            string `json:"apiKey,omitempty"`
	AccessTokenPublic string `json:"accessTokenPublic,omitempty"`
	AccessTokenSecret string `json:"accessTokenSecret,omitempty"`
	BearerToken       string `json:"bearerToken,omitempty"`
}

type ServiceAccountPlatformSecrets struct {
	AwsEksServiceMap map[string]aegis_aws_auth.AuthAWS `json:"awsEksServiceMap"`
}

func GetDockerSecret(ctx context.Context, ou org_users.OrgUser, sn string) (*DockerPlatformSecret, error) {
	m := make(map[string]SecretsKeyValue)
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
	for secretName, v := range m {
		if v.Key == "dockerconfigjson" && secretName == sn {
			dps := &DockerPlatformSecret{DockerAuthJson: v.Value}
			return dps, err
		}
	}
	return nil, err
}

type DockerPlatformSecret struct {
	DockerAuthJson string `json:"dockerAuthJSON"`
}
