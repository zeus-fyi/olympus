package create_org_users

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/keys"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"k8s.io/apimachinery/pkg/util/rand"
)

const Sn = "OrgUser"

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
	UserDemoOrgID             = 1677096191839528000
	EthereumEphemeryServiceID = 1677096782693758000
	EthereumEphemeryService   = "ethereumEphemeryValidators"
	EthereumMainnetServiceID  = 1677096791420465000
	EthereumMainnetService    = "ethereumMainnetValidators"
)

func (o *OrgUser) InsertDemoOrgUserWithNewKey(ctx context.Context, metadata []byte, keyname string, serviceID int) (string, error) {
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
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, UserDemoOrgID, metadata, keyname, keys.BearerKeyTypeID, userKey, serviceID).Scan(&userID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return "", err
	}
	o.UserID = int(userID)
	return userKey, misc.ReturnIfErr(err, q.LogHeader(Sn))
}
