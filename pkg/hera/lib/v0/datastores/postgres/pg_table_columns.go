package postgres

import (
	"strconv"

	"github.com/fraenky8/tables-to-go/pkg/settings"
	"github.com/jackc/pgx/v4/pgxpool"
)

func (d *PgSchemaAutogen) PgxConfigToSqlX(dsnStringPgx string) (*settings.Settings, error) {
	c, err := pgxpool.ParseConfig(dsnStringPgx)
	conf := c.ConnConfig
	pgSettings := settings.New()
	pgSettings.User = conf.User
	pgSettings.Pswd = conf.Password
	pgSettings.Host = conf.Host
	pgSettings.Port = strconv.Itoa(int(conf.Port))
	pgSettings.DbName = conf.Database
	return pgSettings, err
}
