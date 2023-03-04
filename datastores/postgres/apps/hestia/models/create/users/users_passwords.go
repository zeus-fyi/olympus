package create_users

import (
	"context"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (u *User) InsertUserPassword(ctx context.Context, pw string) error {
	q := sql_query_templates.NewQueryParam("InsertUserPassword", "users_passwords", "where", 1000, []string{})
	ps, err := HashPassword(pw)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	q.RawQuery = `INSERT INTO users_passwords(user_id, password)
				  VALUES ($1, $2)
				  RETURNING password_id
				`
	var pwID int64
	err = apps.Pg.QueryRowWArgs(ctx, q.RawQuery, u.UserID, ps).Scan(&pwID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	return err
}
