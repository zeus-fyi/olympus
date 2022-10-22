package postgres

import (
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/datastores/common"
	"github.com/zeus-fyi/tables-to-go/pkg/database"
)

type PgSchemaAutogen struct {
	*database.Postgresql
	// map table to its columns
	Tables map[string]common.Columns
}

func NewPgSchemaAutogen(dsnStringPgx string) PgSchemaAutogen {
	pg := PgSchemaAutogen{}
	pg.InitPG(dsnStringPgx)
	return pg
}

func (d *PgSchemaAutogen) InitPG(dsnStringPgx string) {
	pgConf, err := d.PgxConfigToSqlX(dsnStringPgx)
	if err != nil {
		panic(err)
	}
	d.Postgresql = database.NewPostgresql(pgConf)
	err = d.Connect()
	if err != nil {
		panic(err)
	}
	return
}
