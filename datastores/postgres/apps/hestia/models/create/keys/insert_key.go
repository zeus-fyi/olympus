package create_keys

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/keys"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"golang.org/x/crypto/bcrypt"
)

const Sn = "UserKey"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (k *Key) InsertUserKey(ctx context.Context) error {
	q := sql_query_templates.QueryParams{}
	if k.PublicKeyTypeID == keys.PassphraseKeyTypeID {
		hashedPassword, err := HashPassword(k.PublicKey)
		if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
			return err
		}
		k.PublicKey = hashedPassword
	}
	q.RawQuery = `
				  INSERT INTO users_keys(public_key, user_id, public_key_name, public_key_verified, public_key_type_id)
				  VALUES ($1, $2, $3, $4, $5)
				  `
	r, err := apps.Pg.Exec(ctx, q.RawQuery, k.PublicKey, k.UserID, k.PublicKeyName, k.PublicKeyVerified, k.PublicKeyTypeID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("InsertUserKey: %s, Rows Affected: %d", q.LogHeader(Sn), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func (k *Key) InsertUserSessionKey(ctx context.Context) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `WITH cte_delete_prev_session_keys AS (
				  	DELETE FROM users_keys WHERE user_id = $2 AND public_key_type_id = $5
				  )
				  INSERT INTO users_keys(public_key, user_id, public_key_name, public_key_verified, public_key_type_id)
				  VALUES ($1, $2, $3, $4, $5)
				  `
	r, err := apps.Pg.Exec(ctx, q.RawQuery, k.PublicKey, k.UserID, k.PublicKeyName, k.PublicKeyVerified, keys.SessionIDKeyTypeID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("InsertUserKey: %s, Rows Affected: %d", q.LogHeader(Sn), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
