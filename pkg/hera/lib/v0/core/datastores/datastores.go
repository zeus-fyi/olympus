package datastores

import (
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/datastores/postgres"
)

type DatastoreAutogen struct {
	postgres.PgSchemaAutogen
}

func NewDatastoreAutogen() DatastoreAutogen {
	return DatastoreAutogen{}
}

func (d *DatastoreAutogen) NewInitPGDatastoreAutogen(dsnString string) DatastoreAutogen {
	pg := postgres.NewPgSchemaAutogen()
	pg.NewInitPgConnToSchemaAutogen(dsnString)
	return DatastoreAutogen{pg}
}
