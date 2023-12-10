package aws_secrets

import (
	"fmt"

	"golang.org/x/crypto/sha3"
)

type SecretsKeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
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
