package postgres

import (
	"github.com/fraenky8/tables-to-go/pkg/database"
)

type PgSchemaAutogen struct {
	*database.Postgresql
}

func NewPgSchemaAutogen(dsnStringPgx string) PgSchemaAutogen {
	pg := PgSchemaAutogen{}
	pgConf, err := pg.PgxConfigToSqlX(dsnStringPgx)
	if err != nil {
		panic(err)
	}
	pg.Postgresql = database.NewPostgresql(pgConf)
	err = pg.Connect()
	if err != nil {
		panic(err)
	}
	return pg
}
