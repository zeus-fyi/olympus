package postgres

import (
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs"
	"github.com/zeus-fyi/tables-to-go/pkg/database"
	"github.com/zeus-fyi/tables-to-go/pkg/table_formatting"
)

type PgSchemaAutogen struct {
	*database.Postgresql
	// map table to its columns
	table_formatting.TableContentMap
	StructMapToCodeGen map[string]structs.StructGen
}

func NewPgSchemaAutogen(dsnStringPgx string) PgSchemaAutogen {
	pg := PgSchemaAutogen{}
	pg.InitPG(dsnStringPgx)
	pg.TableContentMap = table_formatting.NewTableContentMap()
	pg.StructMapToCodeGen = make(map[string]structs.StructGen)
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
