package create_org_users

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/keys"
	create_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/keys"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	"k8s.io/apimachinery/pkg/util/rand"
)

const Sn = "OrgUser"

func DoesUserExist(ctx context.Context, email string) bool {
	q := sql_query_templates.NewQueryParam("NewTestOrgUser", "org_users", "where", 1000, []string{})
	q.RawQuery = `SELECT EXISTS(SELECT user_id FROM users WHERE email = $1);`
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	var exists bool
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, email).Scan(&exists)
	if err != nil {
		return exists
	}
	return exists
}

func (o *OrgUser) InsertOrgUser(ctx context.Context, metadata []byte) error {
	q := sql_query_templates.NewQueryParam("NewTestOrgUser", "org_users", "where", 1000, []string{})
	q.RawQuery = `WITH new_user_id AS (
					INSERT INTO users(metadata)
					VALUES ($2)
					RETURNING user_id
				)
					INSERT INTO org_users(org_id, user_id)
					VALUES($1, (SELECT user_id FROM new_user_id))
					RETURNING (SELECT user_id FROM new_user_id)
				`
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	var userID int64
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, o.OrgID, metadata).Scan(&userID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	o.UserID = int(userID)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

const (
	UserDemoOrgID = 1677096191839528000

	EthereumEphemeryServiceID = 1677096782693758000
	IrisQuickNodeService      = "quickNodeMarketPlace"
	IrisService               = "iris"

	EthereumEphemeryService  = "ethereumEphemeryValidators"
	EthereumMainnetServiceID = 1677096791420465000
	EthereumMainnetService   = "ethereumMainnetValidators"
	EthereumGoerliService    = "ethereumGoerliValidators"
	EthereumGoerliServiceID  = 5

	ZeusServiceID         = 1677100016195486976
	ZeusService           = "zeus"
	ZeusWebhooksService   = "zeusWebhooks"
	ZeusWebhooksServiceID = 1677111502304971000
)

func (o *OrgUser) InsertOrgUserWithNewKeyForService(ctx context.Context, metadata []byte, keyname string, serviceID int) (string, error) {
	q := sql_query_templates.NewQueryParam("NewDemoOrgUser", "org_users", "where", 1000, []string{})
	q.RawQuery = `WITH new_user_id AS (
					INSERT INTO users(metadata)
					VALUES ($2)
					RETURNING user_id
				), cte_org_users AS (
					INSERT INTO org_users(org_id, user_id)
					VALUES($1, (SELECT user_id FROM new_user_id))
					RETURNING (SELECT user_id FROM new_user_id)
				), cte_users_key AS (
				  INSERT INTO users_keys(user_id, public_key_name, public_key_verified, public_key_type_id, public_key)
				  VALUES((SELECT user_id FROM new_user_id), $3, true, $4, $5)
			    ) INSERT INTO users_key_services(public_key, service_id)
			      VALUES($5, $6)
				  RETURNING (SELECT user_id FROM new_user_id)
				`
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	var userID int64
	userKey := rand.String(120)
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, o.OrgID, metadata, keyname, keys.BearerKeyTypeID, userKey, serviceID).Scan(&userID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return "", err
	}
	o.UserID = int(userID)
	return userKey, misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func (o *OrgUser) InsertOrgUserWithNewQuickNodeKeyForService(ctx context.Context, quickNodeCustomerID string) error {
	q := sql_query_templates.NewQueryParam("NewDemoOrgUser", "org_users", "where", 1000, []string{})
	q.RawQuery = `WITH new_user_id AS (
					INSERT INTO users(metadata)
					VALUES ('{}')
					RETURNING user_id
				), cte_org_users AS (
					INSERT INTO org_users(org_id, user_id)
					VALUES((SELECT org_id FROM orgs WHERE name = $3), (SELECT user_id FROM new_user_id))
				), cte_quicknode_service AS (
					INSERT INTO users_keys(user_id, public_key_name, public_key_verified, public_key_type_id, public_key)
					VALUES((SELECT user_id FROM new_user_id), $1, true, $2, $3)
					ON CONFLICT (public_key) DO UPDATE SET 
						public_key_verified = EXCLUDED.public_key_verified,
						user_id = EXCLUDED.user_id
				), cte_qn_service AS (
					INSERT INTO users_key_services(public_key, service_id)
					VALUES($3, $4)
					ON CONFLICT (public_key, service_id) DO NOTHING
				), final_insert AS (
					INSERT INTO users_key_services(public_key, service_id)
					VALUES($3, 1677096782693758000)
					ON CONFLICT (public_key, service_id) DO NOTHING
				)
				SELECT 1`
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	var userID int64
	// o.OrgID
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, "quickNodeMarketplaceCustomer", keys.QuickNodeCustomerID, quickNodeCustomerID, 11).Scan(&userID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	o.UserID = int(userID)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
func (o *OrgUser) InsertDemoOrgUserWithNewKey(ctx context.Context, metadata []byte, keyname string, serviceID int) (string, error) {
	q := sql_query_templates.NewQueryParam("NewDemoOrgUser", "org_users", "where", 1000, []string{})
	q.RawQuery = `WITH new_user_id AS (
					INSERT INTO users(metadata, email)
					VALUES ($2, $7)
					RETURNING user_id
				), cte_org_users AS (
					INSERT INTO org_users(org_id, user_id)
					VALUES($1, (SELECT user_id FROM new_user_id))
					RETURNING (SELECT user_id FROM new_user_id)
				), cte_users_key AS (
				  INSERT INTO users_keys(user_id, public_key_name, public_key_verified, public_key_type_id, public_key)
				  VALUES((SELECT user_id FROM new_user_id), $3, true, $4, $5)
			    ) INSERT INTO users_key_services(public_key, service_id)
			      VALUES($5, $6)
				  RETURNING (SELECT user_id FROM new_user_id)
				`
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))

	m := make(map[string]string)
	err := json.Unmarshal(metadata, &m)
	var email string
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error unmarshalling metadata")
	} else {
		if v, ok := m["email"]; ok {
			email = v
		}
	}
	var userID int64
	userKey := rand.String(120)
	err = apps.Pg.QueryRowWArgs(ctx, q.RawQuery, UserDemoOrgID, metadata, keyname, keys.BearerKeyTypeID, userKey, serviceID, &email).Scan(&userID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return "", err
	}
	o.UserID = int(userID)
	OtherServiceSetup(ctx, m, o.UserID)
	return userKey, misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func OtherServiceSetup(ctx context.Context, m map[string]string, userID int) {
	vc, ok := m["validatorCount"]
	if !ok {
		return
	}
	ethAddress, ok := m["ethereumAddress"]
	if !ok {
		log.Ctx(ctx).Info().Msg("no ethereum address found")
		return
	}
	nk := create_keys.NewCreateKey(userID, ethAddress)
	nk.PublicKeyVerified = true
	nk.PublicKeyName = "ethereumAddressEphemery"
	nk.PublicKeyTypeID = keys.EcdsaKeyTypeID
	nk.UserID = userID
	err := nk.InsertUserKey(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error inserting user key")
		return
	}
	validatorCount, err := strconv.Atoi(vc)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error inserting delivery schedule")
		return
	}
	sd := artemis_autogen_bases.EthScheduledDelivery{
		DeliveryScheduleType: "networkReset",
		ProtocolNetworkID:    hestia_req_types.EthereumEphemeryProtocolNetworkID,
		Amount:               validatorCount*artemis_validator_service_groups_models.GweiThirtyTwoEth + artemis_validator_service_groups_models.GweiGasFees,
		Units:                "gwei",
		PublicKey:            ethAddress,
	}
	err = artemis_validator_service_groups_models.InsertDeliverySchedule(ctx, sd)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error inserting delivery schedule")
		return
	}
}
