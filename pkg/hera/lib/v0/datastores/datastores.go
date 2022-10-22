package datastores

import "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/datastores/postgres"

type DatastoreAutogen struct {
	postgres.PgSchemaAutogen
}

func NewDatastoreAutogen() DatastoreAutogen {
	return DatastoreAutogen{}
}
