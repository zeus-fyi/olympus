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

func NewPgSchemaAutogen() PgSchemaAutogen {
	pg := PgSchemaAutogen{
		StructMapToCodeGen: make(map[string]structs.StructGen),
	}
	return pg
}

func (d *PgSchemaAutogen) NewInitPgConnToSchemaAutogen(dsnStringPgx string) {
	d.InitPG(dsnStringPgx)
	d.TableContentMap = table_formatting.NewTableContentMap()
	d.StructMapToCodeGen = make(map[string]structs.StructGen)
	return
}

func (d *PgSchemaAutogen) InitPG(dsnStringPgx string) {
	pgConf, err := d.PgxConfigToSqlX(dsnStringPgx)
	if err != nil {
		panic(err)
	}
	d.Postgresql = database.NewPostgresql(pgConf)
	err = d.Postgresql.Connect()
	if err != nil {
		panic(err)
	}
	d.Settings = pgConf
}
