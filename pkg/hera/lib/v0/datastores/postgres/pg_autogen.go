package postgres

import (
	"github.com/zeus-fyi/tables-to-go/pkg/database"
	"github.com/zeus-fyi/tables-to-go/pkg/table_formatting"
)

type PgSchemaAutogen struct {
	*database.Postgresql
	// map table to its columns
	table_formatting.TableContentMap
}

func NewPgSchemaAutogen(dsnStringPgx string) PgSchemaAutogen {
	pg := PgSchemaAutogen{}
	pg.InitPG(dsnStringPgx)
	pg.TableContentMap = table_formatting.NewTableContentMap()
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
