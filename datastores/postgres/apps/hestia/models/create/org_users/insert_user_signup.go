package create_org_users

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/keys"
	create_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/keys"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"k8s.io/apimachinery/pkg/util/rand"
)

type UserSignup struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	EmailAddress     string `json:"email"`
	Password         string `json:"password"`
	VerifyEmailToken string `json:"verify_email_token"`
}

func (o *OrgUser) InsertSignUpOrgUserAndVerifyEmail(ctx context.Context, us UserSignup) (string, error) {
	q := sql_query_templates.NewQueryParam("InsertSignUpOrgUserWithNewKey", "org_users", "where", 1000, []string{})
	q.RawQuery = `WITH check_cte_user_id AS (
						SELECT user_id
						FROM users
						WHERE email = $3
						LIMIT 1 
					), cte_new_user_id AS (
						INSERT INTO users(first_name, last_name, email, metadata)
						SELECT $1, $2, $3, '{}'
					 	WHERE (SELECT user_id FROM check_cte_user_id) IS NULL
					  	RETURNING user_id
				    ), cte_user_id AS (
						SELECT COALESCE((SELECT user_id FROM cte_new_user_id), (SELECT user_id FROM check_cte_user_id)) AS user_id
					), cte_password_insert AS (
			  			INSERT INTO users_keys(user_id, public_key_name, public_key_verified, public_key_type_id, public_key)
				  		SELECT user_id, $4, $5, $6, $7
						FROM cte_user_id
						WHERE NOT EXISTS (
							SELECT 1
							FROM users_keys u
							WHERE user_id = (SELECT user_id FROM cte_user_id) AND public_key_type_id = $6
					  )
				    ), cte_create_org AS (
						INSERT INTO orgs (name, metadata)	
						SELECT $8, '{}'	
						RETURNING org_id
					), cte_org_users AS (
						INSERT INTO org_users(org_id, user_id)
						SELECT (SELECT org_id FROM cte_create_org), (SELECT user_id FROM cte_user_id)
						WHERE NOT EXISTS (
							SELECT 1
							FROM org_users ou
							WHERE ou.user_id = (SELECT user_id FROM cte_user_id) AND ou.org_id = (SELECT org_id FROM cte_create_org)
						  )
					), cte_delete_prev_verify_token AS (
						DELETE FROM users_keys WHERE user_id = (SELECT user_id FROM cte_user_id) AND public_key_type_id = $9
					), cte_verify_token AS (
	 				  INSERT INTO users_keys(user_id, public_key_name, public_key_verified, public_key_type_id, public_key)
					  VALUES ((SELECT user_id FROM cte_user_id), $10, $11, $9, $12)
					) SELECT user_id FROM cte_user_id
				`
	log.Debug().Interface("InsertSignUpOrgUserWithNewKey:", q.LogHeader(Sn))
	hashedPassword, err := create_keys.HashPassword(us.Password)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return "", err
	}
	signupPwKey := create_keys.NewCreateKey(0, hashedPassword)
	signupPwKey.PublicKeyVerified = false
	signupPwKey.PublicKeyName = "userLoginPassword"
	signupPwKey.PublicKeyTypeID = keys.PassphraseKeyTypeID

	signupKey := create_keys.NewCreateKey(0, rand.String(64))
	signupKey.PublicKeyVerified = false
	signupKey.PublicKeyTypeID = keys.VerifyEmailTokenTypeID
	signupKey.PublicKeyName = "verifyEmailToken"
	us.VerifyEmailToken = signupKey.PublicKey

	demoOrgName := fmt.Sprintf("demoOrg-%s", rand.String(12))
	var userID int64
	err = apps.Pg.QueryRowWArgs(ctx, q.RawQuery,
		us.FirstName, us.LastName, us.EmailAddress,
		signupPwKey.PublicKeyName, signupPwKey.PublicKeyVerified, signupPwKey.PublicKeyTypeID, signupPwKey.PublicKey,
		demoOrgName,
		signupKey.PublicKeyTypeID, signupKey.PublicKeyName, signupKey.PublicKeyVerified, signupKey.PublicKey,
	).Scan(&userID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return "", err
	}
	o.UserID = int(userID)
	return signupKey.PublicKey, misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func UpdateUserEmail(ctx context.Context, userID int, email string) error {
	if userID <= 0 {
		return fmt.Errorf("UpdateUserEmail: userID <= 0")
	}
	q := sql_query_templates.NewQueryParam("UpdateUserEmail", "users", "where", 1000, []string{})
	q.RawQuery = `UPDATE users
				  SET email = $2
				  WHERE user_id = $1 AND email IS NULL;`
	_, err := apps.Pg.Exec(ctx, q.RawQuery, userID, email)
	if err == pgx.ErrNoRows {
		log.Warn().Int("userID", userID).Str("email", email).Msg("UpdateUserEmail: pgx.ErrNoRows")
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}
