package read_keys

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/keys"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_stripe "github.com/zeus-fyi/olympus/pkg/hestia/stripe"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const Sn = "KeyReader"

func (k *OrgUserKey) QueryVerifyUserBearerToken() sql_query_templates.QueryParams {
	var q sql_query_templates.QueryParams
	query := fmt.Sprintf(`
	SELECT usk.public_key_verified, ou.org_id, ou.user_id
	FROM users_keys usk
	INNER JOIN key_types kt ON kt.key_type_id = usk.public_key_type_id
	INNER JOIN org_users ou ON ou.user_id = usk.user_id
	WHERE public_key = $1 AND usk.public_key_verified = true
	`)
	q.RawQuery = query
	return q
}

func (k *OrgUserKey) QueryVerifyUserPassword() sql_query_templates.QueryParams {
	var q sql_query_templates.QueryParams
	query := fmt.Sprintf(`
	SELECT usk.public_key_verified, ou.org_id, ou.user_id
	FROM users u
	INNER JOIN org_users ou ON ou.user_id = u.user_id
	INNER JOIN users_keys usk ON usk.user_id = ou.user_id
	WHERE public_key = crypt($1, public_key) AND u.email = $2
	`)
	q.RawQuery = query
	return q
}

func (k *OrgUserKey) VerifyUserPassword(ctx context.Context, email string) error {
	q := k.QueryVerifyUserPassword()
	log.Debug().Interface("VerifyUserBearerToken:", q.LogHeader(Sn))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, k.PublicKey, email).Scan(&k.PublicKeyVerified, &k.OrgID, &k.UserID)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("VerifyUserPassword error")
		k.PublicKeyVerified = false
	}
	if k.PublicKeyVerified == false {
		return errors.New("unauthorized key")
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func (k *OrgUserKey) QueryUserByEmail() sql_query_templates.QueryParams {
	var q sql_query_templates.QueryParams
	query := fmt.Sprintf(`
	SELECT ou.org_id, ou.user_id
	FROM users u
	INNER JOIN org_users ou ON ou.user_id = u.user_id
	INNER JOIN users_keys usk ON usk.user_id = ou.user_id
	WHERE u.email = $1
	`)
	q.RawQuery = query
	return q
}

func (k *OrgUserKey) GetUserFromEmail(ctx context.Context, email string) error {
	q := k.QueryUserByEmail()
	log.Debug().Interface("GetUserFromEmail:", q.LogHeader(Sn))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, email).Scan(&k.OrgID, &k.UserID)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("GetUserFromEmail error")
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func (k *OrgUserKey) VerifyQuickNodeToken(ctx context.Context) error {
	var q sql_query_templates.QueryParams
	query := fmt.Sprintf(`
	SELECT usk.public_key_verified, ou.org_id, ou.user_id
	FROM users_keys usk
	INNER JOIN key_types kt ON kt.key_type_id = usk.public_key_type_id
	INNER JOIN org_users ou ON ou.user_id = usk.user_id
	WHERE public_key = $1 AND (usk.public_key_type_id = $2 OR usk.public_key_type_id = $3 OR usk.public_key_type_id = $4) AND usk.public_key_verified = true
	LIMIT 1
	`)
	q.RawQuery = query
	log.Debug().Interface("VerifyQuickNodeToken:", q.LogHeader(Sn))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, k.PublicKey, keys.QuickNodeCustomerID, keys.BearerKeyTypeID, keys.SessionIDKeyTypeID).Scan(&k.PublicKeyVerified, &k.OrgID, &k.UserID)
	if err != nil {
		k.PublicKeyVerified = false
	}
	if k.PublicKeyVerified == false {
		return errors.New("unauthorized key")
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func (k *OrgUserKey) VerifyUserBearerToken(ctx context.Context) error {
	q := k.QueryVerifyUserBearerToken()
	log.Debug().Interface("VerifyUserBearerToken:", q.LogHeader(Sn))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, k.PublicKey).Scan(&k.PublicKeyVerified, &k.OrgID, &k.UserID)
	if err != nil {
		k.PublicKeyVerified = false
	}
	if k.PublicKeyVerified == false {
		return errors.New("unauthorized key")
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func (k *OrgUserKey) QueryVerifyUserTokenService() sql_query_templates.QueryParams {
	var q sql_query_templates.QueryParams
	query := fmt.Sprintf(`
	WITH cte_get_user_key AS (
		SELECT usk.user_id AS user_id
		FROM users_keys usk
		WHERE usk.public_key = $1
	) 
	SELECT usk.public_key_verified, ou.org_id, ou.user_id
	FROM users_keys usk
	INNER JOIN key_types kt ON kt.key_type_id = usk.public_key_type_id
	INNER JOIN users_key_services uksvc ON uksvc.public_key = usk.public_key
	INNER JOIN services svcs ON svcs.service_id = uksvc.service_id
	INNER JOIN org_users ou ON ou.user_id = usk.user_id
	WHERE usk.user_id = (SELECT user_id FROM cte_get_user_key) AND svcs.service_name = $2
	`)
	q.RawQuery = query
	return q
}

func (k *OrgUserKey) VerifyUserTokenService(ctx context.Context, serviceName string) error {
	q := k.QueryVerifyUserTokenService()
	log.Debug().Interface("QueryVerifyUserTokenService:", q.LogHeader(Sn))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, k.PublicKey, serviceName).Scan(&k.PublicKeyVerified, &k.OrgID, &k.UserID)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("VerifyUserTokenService error")
		k.PublicKeyVerified = false
		err = k.VerifyQuickNodeToken(ctx)
		if err != nil {
			return err
		}
	}
	if k.PublicKeyVerified == false {
		return errors.New("unauthorized key")
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func (k *OrgUserKey) QueryVerifyUserTokenServiceWithQuickNodePlan() sql_query_templates.QueryParams {
	var q sql_query_templates.QueryParams
	query := fmt.Sprintf(`
		SELECT usk.public_key_verified, ou.org_id, ou.user_id, uks.service_id
		FROM users_keys usk
		INNER JOIN key_types kt ON kt.key_type_id = usk.public_key_type_id
		INNER JOIN org_users ou ON ou.user_id = usk.user_id
		INNER JOIN users_key_services uks ON uks.public_key = usk.public_key
		WHERE usk.public_key = $1 AND uks.service_id = $2 AND usk.public_key_verified = true
		LIMIT 1
	`)
	q.RawQuery = query
	return q
}

func (k *OrgUserKey) VerifyUserTokenServiceWithQuickNodePlan(ctx context.Context, serviceName string) error {
	q := k.QueryVerifyUserTokenServiceWithQuickNodePlan()
	log.Debug().Interface("QueryVerifyUserTokenServiceWithQuickNodePlan:", q.LogHeader(Sn))
	serviceID := 11
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, k.PublicKey, serviceID).Scan(&k.PublicKeyVerified, &k.OrgID, &k.UserID, &serviceID)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("VerifyUserTokenServiceWithQuickNodePlan error")
		k.PublicKeyVerified = false
	}
	if k.PublicKeyVerified == false {
		return errors.New("unauthorized key")
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func (k *OrgUserKey) GetUserID() int {
	return k.UserID
}

func (k *OrgUserKey) IsVerified() bool {
	return k.PublicKeyVerified
}

func (k *OrgUserKey) QueryUserToken() sql_query_templates.QueryParams {
	var q sql_query_templates.QueryParams
	query := fmt.Sprintf(`
	SELECT usk.public_key
	FROM users_keys usk
	INNER JOIN key_types kt ON kt.key_type_id = usk.public_key_type_id
	INNER JOIN org_users ou ON ou.user_id = usk.user_id
	WHERE ou.org_id = $1 AND usk.user_id = $2 AND usk.public_key_type_id = $3
	`)
	q.RawQuery = query
	return q
}

func (k *OrgUserKey) QueryUserBearerToken(ctx context.Context, ou org_users.OrgUser) error {
	q := k.QueryUserToken()
	log.Debug().Interface("QueryUserBearerToken:", q.LogHeader(Sn))
	k.OrgID = ou.OrgID
	k.UserID = ou.UserID
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, ou.OrgID, ou.UserID, keys.BearerKeyTypeID).Scan(&k.PublicKey)
	if err != nil {
		k.PublicKeyVerified = false
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func (k *OrgUserKey) GetUserAuthedServices() sql_query_templates.QueryParams {
	var q sql_query_templates.QueryParams
	query := fmt.Sprintf(`
	WITH cte_get_user_key AS (
		SELECT usk.user_id AS user_id
		FROM users_keys usk
		WHERE usk.public_key = $1
	) 
	SELECT usk.public_key, svcs.service_name, ou.org_id, ou.user_id, COALESCE(qms.plan, 'generic') as plan
	FROM users_keys usk
	INNER JOIN key_types kt ON kt.key_type_id = usk.public_key_type_id
	INNER JOIN users_key_services uksvc ON uksvc.public_key = usk.public_key
	INNER JOIN services svcs ON svcs.service_id = uksvc.service_id
	INNER JOIN org_users ou ON ou.user_id = usk.user_id
	LEFT JOIN quicknode_marketplace_customer qms ON qms.quicknode_id = usk.public_key
	WHERE usk.user_id = (SELECT user_id FROM cte_get_user_key) AND usk.public_key_verified = true
	GROUP BY usk.public_key, svcs.service_name, qms.plan, ou.org_id, ou.user_id
	`)
	q.RawQuery = query
	return q
}

func (k *OrgUserKey) GetUserKeysToServicesQuery() sql_query_templates.QueryParams {
	var q sql_query_templates.QueryParams
	query := fmt.Sprintf(`
	WITH cte_get_user_key AS (
		SELECT usk.user_id AS user_id
		FROM users_keys usk
		WHERE usk.public_key = $1
	)
	SELECT usk.public_key
	FROM users_keys usk
	WHERE usk.user_id = (SELECT user_id FROM cte_get_user_key)
	`)
	q.RawQuery = query
	return q
}

func (k *OrgUserKey) GetUserKeysToServices(ctx context.Context, token string) (map[string]string, error) {
	q := k.GetUserKeysToServicesQuery()
	log.Debug().Interface("QueryUserAuthedServices:", q.LogHeader(Sn))

	keysFound := make(map[string]string)
	rows, err := apps.Pg.Query(ctx, q.RawQuery, token)
	log.Err(err).Interface("QueryUserAuthedServices: Query: ", q.RawQuery)
	if err != nil {
		return keysFound, err
	}
	defer rows.Close()
	for rows.Next() {
		rowErr := rows.Scan(&k.PublicKey)
		if rowErr != nil {
			log.Err(rowErr).Interface("QueryUserAuthedServices: Query: ", q.RawQuery)
			return keysFound, rowErr
		}
		keysFound[k.PublicKey] = k.PublicKey
	}
	return keysFound, misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func (k *OrgUserKey) QueryUserAuthedServices(ctx context.Context, token string) ([]string, map[string]string, error) {
	q := k.GetUserAuthedServices()
	log.Debug().Interface("QueryUserAuthedServices:", q.LogHeader(Sn))

	m := make(map[string]string)
	keysFound := make(map[string]string)
	var services []string
	rows, err := apps.Pg.Query(ctx, q.RawQuery, token)
	log.Err(err).Interface("QueryUserAuthedServices: Query: ", q.RawQuery)
	if err != nil {
		return services, keysFound, err
	}
	defer rows.Close()
	for rows.Next() {
		var serviceName string
		var plan string
		rowErr := rows.Scan(&k.PublicKey, &serviceName, &k.OrgID, &k.UserID, &plan)
		if rowErr != nil {
			log.Err(rowErr).Interface("QueryUserAuthedServices: Query: ", q.RawQuery)
			return services, keysFound, rowErr
		}
		// todo fix, this is not getting unique keys
		keysFound[k.PublicKey] = k.PublicKey
		services = append(services, serviceName)

		if m[serviceName] == "" || m[serviceName] == "generic" {
			m[serviceName] = plan
		}
	}
	k.Services = m
	return services, keysFound, misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func (k *OrgUserKey) QueryGetCustomerStripeID() sql_query_templates.QueryParams {
	var q sql_query_templates.QueryParams
	query := fmt.Sprintf(`
	SELECT usk.public_key
	FROM users_keys usk
	INNER JOIN key_types kt ON kt.key_type_id = usk.public_key_type_id
	INNER JOIN org_users ou ON ou.user_id = usk.user_id
	INNER JOIN users u ON u.user_id = ou.user_id
	WHERE u.user_id = $1 AND usk.public_key_type_id = $2
	`)
	q.RawQuery = query
	return q
}

func (k *OrgUserKey) QueryGetUserInfo() sql_query_templates.QueryParams {
	var q sql_query_templates.QueryParams
	query := fmt.Sprintf(`
	SELECT u.first_name, u.last_name, u.email
	FROM users u
	WHERE u.user_id = $1
	`)
	q.RawQuery = query
	return q
}

func GetOrCreateCustomerStripeIDWithUserID(ctx context.Context, userID int) (string, error) {
	k := OrgUserKey{
		Key: keys.Key{
			UsersKeys: autogen_bases.UsersKeys{
				UserID: userID,
			},
		},
	}
	return k.GetOrCreateCustomerStripeID(ctx)
}

func (k *OrgUserKey) GetOrCreateCustomerStripeID(ctx context.Context) (string, error) {
	q := k.QueryGetCustomerStripeID()
	log.Debug().Interface("GetOrCreateCustomerStripeID:", q.LogHeader(Sn))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, k.UserID, keys.StripeCustomerID).Scan(&k.PublicKey)
	if err == pgx.ErrNoRows {
		q = k.QueryGetUserInfo()
		firstName, lastName, email := "", "", ""
		err = apps.Pg.QueryRowWArgs(ctx, q.RawQuery, k.UserID).Scan(&firstName, &lastName, &email)
		if err != nil {
			log.Ctx(ctx).Err(err).Msg("CreateOrGetCustomerStripeID error")
			return "", err
		}
		c, cerr := hestia_stripe.CreateCustomer(ctx, k.UserID, firstName, lastName, email)
		if cerr != nil {
			log.Ctx(ctx).Err(cerr).Msg("CreateOrGetCustomerStripeID error")
			return "", cerr
		}
		return c.ID, nil
	}
	return k.PublicKey, misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func DeactivateQuickNodeApiKey(ctx context.Context, qid string) (int, error) {
	var q sql_query_templates.QueryParams
	query := fmt.Sprintf(`
	WITH cte_update_key AS (
		UPDATE users_keys
		SET public_key_verified = false
		WHERE public_key = $1
		RETURNING user_id
	) SELECT org_id FROM org_users WHERE user_id = (SELECT user_id FROM cte_update_key)
	`)
	q.RawQuery = query
	orgID := 0
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, qid).Scan(&orgID)
	if err == pgx.ErrNoRows {
		log.Warn().Msg("No key to deactivate")
		return orgID, nil
	}
	return orgID, err
}
