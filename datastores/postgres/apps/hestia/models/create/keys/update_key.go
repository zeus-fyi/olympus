package create_keys

import (
	"context"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/keys"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"k8s.io/apimachinery/pkg/util/rand"
)

const (
	UserDemoOrgID = 1677096191839528000

	EthereumEphemeryServiceID = 1677096782693758000
	EthereumEphemeryService   = "ethereumEphemeryValidators"
	EthereumMainnetServiceID  = 1677096791420465000
	EthereumMainnetService    = "ethereumMainnetValidators"
	EthereumGoerliService     = "ethereumGoerliValidators"
	EthereumGoerliServiceID   = 5

	ZeusServiceID         = 1677100016195486976
	ZeusService           = "zeus"
	ZeusWebhooksService   = "zeusWebhooks"
	ZeusWebhooksServiceID = 1677111502304971000
)

func UpdateKeysFromVerifyEmail(ctx context.Context, tokenID string) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
				  WITH cte_delete_verification_token AS (
				  	DELETE FROM users_keys WHERE public_key = $1 AND public_key_type_id = $2
					RETURNING user_id
				  ), cte_users_key AS (
					  INSERT INTO users_keys(user_id, public_key_name, public_key_verified, public_key_type_id, public_key)
					  VALUES((SELECT user_id FROM cte_delete_verification_token), $3, true, $4, $5)
				  ), cte_insert_basic_service_keys AS (
					 INSERT INTO users_key_services(service_id, public_key)
					 VALUES ($6, $5)
				  )
				  UPDATE users_keys
				  SET public_key_verified = true
				  WHERE public_key_type_id = $7 AND user_id = (SELECT user_id FROM cte_delete_verification_token)
				  RETURNING user_id
				  `
	var userID int
	bearer := rand.String(120)
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery,
		tokenID, keys.VerifyEmailTokenTypeID,
		"userDemo", keys.BearerKeyTypeID, bearer,
		EthereumEphemeryServiceID,
		keys.PassphraseKeyTypeID).Scan(&userID)
	if err != nil {
		return err
	}
	return nil
}
